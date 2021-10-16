package main

import (
	"github.com/go-chi/chi"
	"github.com/ianmuhia/bookings/internals/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
	//do nothing testpassed
	default:
		t.Errorf("type is not *chi.Mux, but is %T", v)
	}
}
