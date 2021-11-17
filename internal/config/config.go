package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/go-redis/redis/v8"
	"github.com/ianmuhia/bookings/internal/models"
	"github.com/sirupsen/logrus"
)

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Logrus        *logrus.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
	Cache         *redis.Client
}
