package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/ianmuhia/bookings/internals/config"
	"github.com/ianmuhia/bookings/internals/driver"
	"github.com/ianmuhia/bookings/internals/handlers"
	"github.com/ianmuhia/bookings/internals/helpers"
	"github.com/ianmuhia/bookings/internals/models"
	"github.com/ianmuhia/bookings/internals/render"
	"log"
	"net/http"
	"os"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

const (
	port = ":8080"
)

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	log.Printf("starting app in %s", port)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	gob.Register(models.Reservation{})

	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=booking user=postgres password=*#*Johnte2536")
	if err != nil {
		log.Fatal("can't connect to database")
	}
	log.Println("connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)
	return db, nil
}
