package config

import (
	"github.com/alexedwards/scs/v2"
	"github.com/ianmuhia/bookings/internals/models"
	"html/template"
	"log"
)

// AppConfig holds the application wide config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
