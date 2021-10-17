package repository

import (
	"time"

	"github.com/moaabb/bookings-go/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	// InsertReservation inserts a reservation into the database
	InsertReservation(res models.Reservation) (int, error)
	// InsertRoomRestriction Inserts the restriction to the room
	InsertRoomRestriction(r models.RoomRestriction) error
	// SearchAvailabilityByDatesByRoomID Takes a roomID and the dates and returns, if any, if it's available
	SearchAvailabilityByDatesByRoomID(start_date, end_date time.Time, room_id int) (bool, error)
	// SearchAvailabilityForAllRooms Takes the dates and return a slice of the available rooms for that time window
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	// GetRoomByID Receives and ID and returns the correspondent Room
	GetRoomByID(id int) (models.Room, error)

	// GetUserByID returns a user by it's id
	GetUserByID(id int) (models.User, error)
	// UpdateUser updates a user in the database
	UpdateUser(u models.User) error
	// Authenticate Authenticates a user
	Authenticate(email, testPassword string) (int, string, error)

	// AllReservations Returns all the reservations
	AllReservations() ([]models.Reservation, error)
	// AllNewReservations Returns All the reservations not yet processed
	AllNewReservations() ([]models.Reservation, error)
	// GetReservationByID Returns the reservation with the given ID
	GetReservationByID(id int) (models.Reservation, error)
	// UpdateReservation updates a reservation in the database
	UpdateReservation(r models.Reservation) error
	// DeleteReservation Deletes a reservation from the database
	DeleteReservation(id int) error
	// UpdateProcessedForReservation Updates processed status for a reservation by id
	UpdateProcessedForReservation(id, processed int) error

	// AllRooms Returns All the rooms
	AllRooms() ([]models.Room, error)

	// GetRstrictionsForRoomByDate Return all the restrictions for the given room
	GetRstrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error)
	// InsertRoomBlock Inserts Owner block for the room for the given date
	InsertRoomBlock(roomID int, startDate time.Time) error
	// DeleteRoomBlock Deletes Owner block for the room for the given date
	DeleteRoomBlock(id int) error
}
