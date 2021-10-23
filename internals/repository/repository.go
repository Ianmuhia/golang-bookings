package repository

import "github.com/ianmuhia/bookings/internals/models"

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) error
}
