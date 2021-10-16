package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/ianmuhia/bookings/internals/config"
	"github.com/ianmuhia/bookings/internals/handlers"
	"github.com/ianmuhia/bookings/internals/models"
	"github.com/ianmuhia/bookings/internals/render"
	"log"
	"net/http"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager

const (
	port = ":8080"
)

func main() {
	gob.Register(models.Reservation{})

	app.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	log.Printf("starting app in %s", port)

	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
