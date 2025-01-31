package repositories

import (
	"context"
	"errors"
	"fmt"
	"movie-system/internal/models"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	DB *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{DB: db}
}

func (r *ReservationRepository) ReserveSeat(ctx context.Context, reservation *models.Reservation) error {
	seatsArray := "{" + strings.Join(reservation.Seats, ",") + "}"
	var count int
	checkSeatsQuery := `
		SELECT COUNT(*)
		FROM reservations
		WHERE showtime_id = $1
		AND seats && $2;
	`
	err := r.DB.QueryRow(ctx, checkSeatsQuery, reservation.ShowtimeID, seatsArray).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking seat availability: %w", err)
	}

	if count > 0 {
		return errors.New("one or more seats are already reserved")
	}

	insertReservationQuery := `
		INSERT INTO reservations (user_id, movie_id, showtime_id, seats)
		VALUES ($1, $2, $3, $4) RETURNING id;
	`

	var reservationID int
	err = r.DB.QueryRow(ctx, insertReservationQuery, reservation.UserID, reservation.MovieID, reservation.ShowtimeID, seatsArray).Scan(&reservationID)
	if err != nil {
		return fmt.Errorf("error creating reservation: %w", err)
	}

	incrementReservedQuery := `
		UPDATE showtimes
		SET reserved = reserved + $1
		WHERE id = $2;
	`

	numSeats := len(reservation.Seats)
	_, err = r.DB.Exec(ctx, incrementReservedQuery, numSeats, reservation.ShowtimeID)
	if err != nil {
		return fmt.Errorf("error updating reserved seats: %w", err)
	}

	fmt.Printf("Reservation successful! Reservation ID: %d\n", reservationID)
	return nil
}
