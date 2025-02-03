package repositories

import (
	"context"
	"errors"
	"fmt"
	"movie-system/internal/models"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
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

	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transcation: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				fmt.Printf("error rolling back transcation: %v\n", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				fmt.Printf("error committing transcation: %v\n", commitErr)
			}
		}
	}()

	checkSeatsQuery := `
		SELECT COUNT(*)
		FROM reservations
		WHERE showtime_id = $1
		AND seats && $2;
	`
	err = tx.QueryRow(ctx, checkSeatsQuery, reservation.ShowtimeID, seatsArray).Scan(&count)
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
	err = tx.QueryRow(ctx, insertReservationQuery, reservation.UserID, reservation.MovieID, reservation.ShowtimeID, seatsArray).Scan(&reservationID)
	if err != nil {
		return fmt.Errorf("error creating reservation: %w", err)
	}

	incrementReservedQuery := `
		UPDATE showtimes
		SET reserved = reserved + $1
		WHERE id = $2;
	`

	numSeats := len(reservation.Seats)
	_, err = tx.Exec(ctx, incrementReservedQuery, numSeats, reservation.ShowtimeID)
	if err != nil {
		return fmt.Errorf("error updating reserved seats: %w", err)
	}

	return nil
}

func (r *ReservationRepository) CancelReservation(ctx context.Context, reservationID int) error {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				fmt.Printf("error rolling back transcation: %v\n", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				fmt.Printf("error committing transcation: %v\n", commitErr)
			}
		}
	}()

	var showtimeID int
	var seatsArray []string
	checkReservationQuery := `
		SELECT showtime_id, seats
		FROM reservations
		WHERE id = $1;
	`
	err = tx.QueryRow(ctx, checkReservationQuery, reservationID).Scan(&showtimeID, &seatsArray)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("reservation not found")
		}
		return fmt.Errorf("error checking reservation: %w", err)
	}

	var showtimeTime time.Time
	checkShowtimeQuery := `
		SELECT start_time
		FROM showtimes 
		WHERE id = $1;
	`
	err = tx.QueryRow(ctx, checkShowtimeQuery, showtimeID).Scan(&showtimeTime)
	if err != nil {
		return fmt.Errorf("error checking showtime: %w", err)
	}

	if showtimeTime.Before(time.Now()) {
		return errors.New("cannot cancel a reservation for a past showtime")
	}

	numSeats := len(seatsArray)
	decrementReservedQuery := `
		UPDATE showtimes
		SET reserved = reserved - $1
		WHERE id = $2;
	`
	_, err = tx.Exec(ctx, decrementReservedQuery, numSeats, showtimeID)
	if err != nil {
		return fmt.Errorf("error updating reserved seats: %w", err)
	}

	deleteReservationQuery := `
	DELETE FROM reservations
	WHERE id = $1;
	`
	_, err = tx.Exec(ctx, deleteReservationQuery, reservationID)
	if err != nil {
		return fmt.Errorf("error deleting reservation: %w", err)
	}

	fmt.Printf("Reservation %d canceled successfully.\n", reservationID)
	return nil
}
