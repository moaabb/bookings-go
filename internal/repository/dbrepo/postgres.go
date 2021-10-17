package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/moaabb/bookings-go/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {

	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `INSERT INTO reservations
			(first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var newID int

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction Inserts the restriction to the room
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions
			(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
			VALUES
			($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomID Takes a roomID and the dates and returns, if any, if it's available
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start_date, end_date time.Time, room_id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var numRows int

	stmt := `
		SELECT
			count(id)
		FROM
			room_restrictions
		WHERE
			room_id = $1
			AND $2 <= end_date AND $3 >= start_date
		`

	err := m.DB.QueryRowContext(ctx, stmt, room_id, start_date, end_date).Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms Takes the dates and return a slice of the available rooms for that time window
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		SELECT
			r.id, r.room_name
		FROM
			rooms r
		WHERE
			r.id NOT IN
		(SELECT rr.room_id FROM room_restrictions rr WHERE $1 < end_date AND $2 > start_date)
	`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(
			&room.ID,
			&room.RoomName,
		)

		if err != nil {
			return rooms, nil
		}

		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil

}

// GetRoomByID Receives and ID and returns the correspondent Room
func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var room models.Room

	query := `
		SELECT
			id, room_name, created_at, updated_at
		FROM
			rooms
		WHERE
			id = $1
	`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserByID returns a user by it's id
func (m *postgresDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	query := `
		SELECT
			id, first_name, last_name, email, password, access_level, created_at, updated_at
		FROM
			user
		WHERE
			id = $1
	`
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil

}

// UpdateUser updates a user in the database
func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate Authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	query := `
		SELECT id, password FROM users WHERE email = $1
	`
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&id,
		&hashedPassword,
	)
	if err != nil {
		return 0, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("Incorrect Password!")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

// AllReservations Returns all the reservations
func (m *postgresDBRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		SELECT
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.processed, r.created_at, r.updated_at, rm.room_name, rm.id
		FROM
			reservations r
		LEFT JOIN
			rooms rm ON (r.room_id = rm.id)
		ORDER BY
			start_date ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.Processed,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.RoomName,
			&i.Room.ID,
		)
		if err != nil {
			return reservations, nil
		}

		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// AllNewReservations Returns All the reservations not yet processed
func (m *postgresDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	query := `
		SELECT
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.processed, r.created_at, r.updated_at, rm.room_name, rm.id
		FROM
			reservations r
		LEFT JOIN
			rooms rm ON (r.room_id = rm.id)
		WHERE
			r.processed = 0
		ORDER BY
			start_date ASC
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservations, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Reservation
		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.Processed,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.RoomName,
			&i.Room.ID,
		)
		if err != nil {
			return reservations, nil
		}

		reservations = append(reservations, i)
	}
	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

// GetReservationByID Returns the reservation with the given ID
func (m *postgresDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservation models.Reservation

	query := `
		SELECT
			r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.processed, r.created_at, r.updated_at, rm.room_name, rm.id
		FROM
			reservations r
		LEFT JOIN
			rooms rm ON (r.room_id = rm.id)
		WHERE
			r.id = $1
	`

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.FirstName,
		&reservation.LastName,
		&reservation.Email,
		&reservation.Phone,
		&reservation.StartDate,
		&reservation.EndDate,
		&reservation.RoomID,
		&reservation.Processed,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
		&reservation.Room.RoomName,
		&reservation.Room.ID,
	)
	if err != nil {
		return reservation, err
	}

	return reservation, nil
}

// UpdateReservation updates a reservation in the database
func (m *postgresDBRepo) UpdateReservation(r models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE reservations SET first_name = $1, last_name = $2, email = $3, phone = $4, updated_at = $5
		WHERE id = $6
	`

	_, err := m.DB.ExecContext(ctx, query,
		r.FirstName,
		r.LastName,
		r.Email,
		r.Phone,
		time.Now(),
		r.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteReservation Deletes a reservation from the database
func (m *postgresDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM reservations WHERE id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProcessedForReservation Updates processed status for a reservation by id
func (m *postgresDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE reservations SET processed = $1 WHERE id = $2
	`

	_, err := m.DB.ExecContext(ctx, query, processed, id)
	if err != nil {
		return err
	}

	return nil
}

// AllRooms Returns All the rooms
func (m *postgresDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rooms []models.Room

	query := `
		SELECT id, room_name, created_at, updated_at FROM rooms
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Room
		err = rows.Scan(
			&r.ID,
			&r.RoomName,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, r)
	}
	if err = rows.Err(); err != nil {
		return rooms, err
	}

	return rooms, nil
}

// GetRstrictionsForRoomByDate Return all the restrictions for the given room
func (m *postgresDBRepo) GetRstrictionsForRoomByDate(roomID int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var restrictions []models.RoomRestriction

	query := `
		SELECT id, start_date, end_date, room_id, COALESCE(reservation_id, 0), restriction_id
		FROM room_restrictions
		WHERE start_date >= $1 AND $2 >= end_date AND room_id = $3
	`
	rows, err := m.DB.QueryContext(ctx, query, start, end, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.RoomRestriction
		err = rows.Scan(
			&r.ID,
			&r.StartDate,
			&r.EndDate,
			&r.RoomID,
			&r.ReservationID,
			&r.RestrictionID,
		)
		if err != nil {
			return nil, err
		}

		restrictions = append(restrictions, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restrictions, nil
}

// InsertRoomBlock Inserts Owner block for the room for the given date
func (m *postgresDBRepo) InsertRoomBlock(roomID int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO room_restrictions
		(start_date, end_date, room_id, restriction_id, created_at, updated_at)
		VALUES
		($1, $2, $3, $4, $5, $6)
	`

	_, err := m.DB.ExecContext(ctx, query,
		startDate,
		startDate.AddDate(0, 0, 1),
		roomID,
		2,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRoomBlock Deletes Owner block for the room for the given date
func (m *postgresDBRepo) DeleteRoomBlock(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		DELETE FROM room_restrictions WHERE id = $1
	`

	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
