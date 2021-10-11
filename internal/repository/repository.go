package repository

import (
	"time"

	"github.com/moaabb/bookings-go/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start_date, end_date time.Time, room_id int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
}
