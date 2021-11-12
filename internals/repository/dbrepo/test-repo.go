package dbrepo

import (
	"errors"
	"github.com/ianmuhia/bookings/internals/models"
	"time"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts the reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2 the fail otherwise fail
	if res.RoomID == 2 {
		return 0, errors.New("some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restrictions into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}

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
	if id > 2 {
		return room, errors.New("some error")
	}

	return room, nil
}
