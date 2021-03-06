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

// GetUserByID returns a user by it's id
func (m *testingDBRepo) GetUserByID(id int) (models.User, error) {

	var user models.User

	return user, nil

}

// UpdateUser updates a user in the database
func (m *testingDBRepo) UpdateUser(u models.User) error {

	return nil
}

// Authenticate Authenticates a user
func (m *testingDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	var id int
	var hashedPassword string

	return id, hashedPassword, nil
}

// AllReservations Returns all the reservations
func (m *testingDBRepo) AllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// AllNewReservations Returns All the reservations not yet processed
func (m *testingDBRepo) AllNewReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation

	return reservations, nil
}

// GetReservationByID Returns the reservation with the given ID
func (m *testingDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	var reservation models.Reservation

	return reservation, nil
}

// UpdateReservation updates a reservation in the database
func (m *testingDBRepo) UpdateReservation(r models.Reservation) error {

	return nil
}

// DeleteReservation Deletes a reservation from the database
func (m *testingDBRepo) DeleteReservation(id int) error {

	return nil
}

// UpdateProcessedForReservation Updates processed status for a reservation by id
func (m *testingDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

// AllRooms Returns All the rooms
func (m *testingDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

// GetRstrictionsForRoomByDate Return all the restrictions for the given room
func (m *testingDBRepo) GetRstrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {

	var restrictions []models.RoomRestriction

	return restrictions, nil
}

// InsertRoomBlock Inserts Owner block for the room for the given date
func (m *testingDBRepo) InsertRoomBlock(roomID int, startDate time.Time) error {

	return nil
}

// DeleteRoomBlock Deletes Owner block for the room for the given date
func (m *testingDBRepo) DeleteRoomBlock(id int) error {

	return nil
}
