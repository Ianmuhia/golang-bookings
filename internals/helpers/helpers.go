package helpers

import "github.com/ianmuhia/bookings/internals/config"

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}
