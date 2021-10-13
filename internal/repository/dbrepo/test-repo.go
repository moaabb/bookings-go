package dbrepo

import (
	"time"

	"github.com/moaabb/bookings-go/internal/models"
)

func (m *testingDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database
func (m *testingDBRepo) InsertReservation(res models.Reservation) (int, error) {
	var newID int

	return newID, nil
}

// InsertRoomRestriction Inserts the restriction to the room
func (m *testingDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	return nil
}

// SearchAvailabilityByDatesByRoomID Takes a roomID and the dates and returns, if any, if it's available
func (m *testingDBRepo) SearchAvailabilityByDatesByRoomID(start_date, end_date time.Time, room_id int) (bool, error) {

	return false, nil
}

// SearchAvailabilityForAllRooms Takes the dates and return a slice of the available rooms for that time window
func (m *testingDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

// GetRoomByID Receives and ID and returns the correspondent Room
func (m *testingDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	return room, nil
}
