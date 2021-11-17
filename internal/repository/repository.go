package repository

import (
	"time"

	"github.com/ianmuhia/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	UpdateUser(u models.User) error
	GetUserById(id int) (models.User, error)
	Authenticate(email, testPassord string) (int, string, error)
}
