package dbrepo

import (
	"github.com/ianmuhia/bookings/internals/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts the reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction inserts a room restrictions into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for room id and false if no availability
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

//GetRoomByID gets a room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	return room, nil
}
