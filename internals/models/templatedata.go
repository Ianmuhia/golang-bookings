package models

import "github.com/ianmuhia/bookings/internals/forms"

// TemplateData holds data sent from handlers package
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	floatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
