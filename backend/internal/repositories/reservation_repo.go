package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (repo *ReservationRepository) GetReservations(ctx context.Context, id int) ([]models.Reservation, error) {
	var rows pgx.Rows
	var err error

	query := `
		SELECT id, user_id, movie_id, showtime_id, seats, created_at
		FROM reservations
		WHERE user_id = $1
	`

	rows, err = repo.DB.Query(ctx, query, id)
	if err != nil {
		log.Printf("error fetching reservations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var reservation models.Reservation
		if err := rows.Scan(&reservation.ID, &reservation.UserID, &reservation.MovieID, &reservation.ShowtimeID, &reservation.Seats, &reservation.CreatedAt); err != nil {
			log.Printf("error scanning reservations: %v", err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error during rows iteration: %v", err)
		return nil, err
	}

	return reservations, err
}

func (repo *ReservationRepository) GetAllReservations(ctx context.Context) ([]models.Reservation, error) {
	var rows pgx.Rows
	var err error

	query := `
		SELECT id, user_id, movie_id, showtime_id, seats, created_at
		FROM reservations
	`

	rows, err = repo.DB.Query(ctx, query)
	if err != nil {
		log.Printf("error fetching reservations: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var reservation models.Reservation
		err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.MovieID,
			&reservation.ShowtimeID,
			&reservation.Seats,
			&reservation.CreatedAt,
		)
		if err != nil {
			log.Printf("error scanning reservations: %v", err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error during rows iteration: %v", err)
		return nil, err
	}

	return reservations, nil
}

// GetReservationsPerMovie returns the count of reservations for each movie
func (r *ReservationRepository) GetReservationsPerMovie(ctx context.Context, movieID int) ([]models.MovieReservationCount, error) {
	query := `
		SELECT 
			m.id,
			m.title,
			COUNT(r.id) as reservation_count,
			COALESCE(SUM(array_length(r.seats, 1)), 0) as total_seats
		FROM movies m
		LEFT JOIN reservations r ON m.id = r.movie_id
		WHERE m.id = $1
		GROUP BY m.id, m.title
		ORDER BY reservation_count DESC`

	rows, err := r.DB.Query(ctx, query, movieID)
	if err != nil {
		return nil, fmt.Errorf("error querying reservations per movie: %v", err)
	}
	defer rows.Close()

	var results []models.MovieReservationCount
	for rows.Next() {
		var count models.MovieReservationCount
		if err := rows.Scan(
			&count.MovieID,
			&count.MovieTitle,
			&count.ReservationCount,
			&count.SeatCount,
		); err != nil {
			return nil, fmt.Errorf("error scanning reservation count: %v", err)
		}
		results = append(results, count)
	}

	return results, nil
}

// GetTotalRevenue returns the total revenue and seat count across all movies
func (r *ReservationRepository) GetTotalRevenue(ctx context.Context) (int, map[string]int, int, error) {
	query := `
		SELECT 
			m.title,
			COALESCE(SUM(array_length(r.seats, 1)), 0) as seats_reserved,
			COALESCE(SUM(array_length(r.seats, 1)), 0) * 5 as revenue
		FROM movies m
		LEFT JOIN reservations r ON m.id = r.movie_id
		GROUP BY m.id, m.title
		ORDER BY revenue DESC`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return 0, nil, 0, fmt.Errorf("error querying total revenue: %v", err)
	}
	defer rows.Close()

	revenues := make(map[string]int)
	totalRevenue := 0
	totalSeatsReserved := 0

	for rows.Next() {
		var (
			movieTitle    string
			seatsReserved int
			revenue       int
		)
		if err := rows.Scan(&movieTitle, &seatsReserved, &revenue); err != nil {
			return 0, nil, 0, fmt.Errorf("error scanning revenue data: %v", err)
		}

		revenues[movieTitle] = revenue
		totalRevenue += revenue
		totalSeatsReserved += seatsReserved
	}

	if err := rows.Err(); err != nil {
		return 0, nil, 0, fmt.Errorf("error during rows iteration: %v", err)
	}

	return totalSeatsReserved, revenues, totalRevenue, nil
}
