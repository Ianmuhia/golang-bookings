package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-redis/redis/v8"
	"github.com/ianmuhia/bookings/internal/config"
	"github.com/ianmuhia/bookings/internal/driver"
	"github.com/ianmuhia/bookings/internal/handlers"
	"github.com/ianmuhia/bookings/internal/helpers"
	"github.com/ianmuhia/bookings/internal/models"
	"github.com/ianmuhia/bookings/internal/render"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	// "go.elastic.co/ecslogrus"
)

const portNumber = ":8080"

const LOG_FILE = "./logrus.log"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var logger *logrus.Logger

// main is the main application
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	log.Println("Starting mail listener")

	listenForMail()

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app.InfoLog = infoLog

	///set up logrus for logging
	logger = logrus.New()
	logger.SetFormatter(&ecslogrus.Formatter{
		PrettyPrint: true,
	})
	logger.ReportCaller = true

	logger.SetOutput(os.Stdout)

	file, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	// defer file.Close()
	logger.SetOutput(file)

	app.Logrus = logger

	///set up logrus for logging

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// setup cache
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	app.Cache = rdb
	log.Println("setup redis complete")

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=*#*Johnte2536")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Println(err)
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
