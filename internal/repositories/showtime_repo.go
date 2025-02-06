package repositories

import (
	"context"
	"fmt"
	"log"
	"movie-system/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShowtimeRepository struct {
	DB *pgxpool.Pool
}

func NewShowtimeRepository(db *pgxpool.Pool) *ShowtimeRepository {
	return &ShowtimeRepository{DB: db}
}

func (repo *ShowtimeRepository) InsertShowtime(ctx context.Context, showtime *models.Showtime) error {
	_, err := repo.DB.Exec(ctx, `
		INSERT INTO showtimes (movie_id, start_time, capacity, reserved)
		VALUES ($1,$2,$3,$4)`,
		showtime.MovieID, showtime.StartTime, showtime.Capacity, showtime.Reserved)
	return err
}

func (repo *ShowtimeRepository) GetShowtimes(ctx context.Context) ([]models.Showtime, error) {
	var rows pgx.Rows
	var err error

	rows, err = repo.DB.Query(ctx, `
			SELECT id, movie_id, start_time, capacity, reserved
			FROM showtimes`)
	if err != nil {
		log.Printf("error fetching showtimes: %v", err)
		return nil, err
	}
	defer rows.Close()

	var showtimes []models.Showtime
	for rows.Next() {
		var showtime models.Showtime
		if err := rows.Scan(&showtime.ID, &showtime.MovieID, &showtime.StartTime, &showtime.Capacity, &showtime.Reserved); err != nil {
			log.Printf("error scanning showtime: %v", err)
			return nil, err
		}
		showtimes = append(showtimes, showtime)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error during rows iteration: %v", err)
		return nil, err
	}

	return showtimes, err
}

func (repo *ShowtimeRepository) UpdateShowtime(ctx context.Context, id int, showtime *models.Showtime) error {
	_, err := repo.DB.Exec(ctx, `
		UPDATE showtimes
		SET movie_id = $1, start_time = $2, capacity = $3, reserved = $4
		WHERE id = $5`,
		showtime.MovieID, showtime.StartTime, showtime.Capacity, showtime.Reserved, id)
	return err
}

func (repo *ShowtimeRepository) DeleteShowtime(ctx context.Context, id int) error {
	_, err := repo.DB.Exec(ctx, "DELETE FROM showtimes WHERE id = $1", id)
	return err
}

func (repo *ShowtimeRepository) GetAvailableSeats(ctx context.Context, id int) ([]string, error) {
	query := `
		SELECT seats
		FROM reservations
		WHERE showtime_id = $1
	`
	rows, err := repo.DB.Query(ctx, query, id)
	if err != nil {
		log.Printf("error fetching reserved seats: %v", err)
		return nil, err
	}
	defer rows.Close()

	reservedSeats := make(map[string]bool)
	for rows.Next() {
		var seatsArray []string
		if err := rows.Scan(&seatsArray); err != nil {
			log.Printf("error scanning seats: %v", err)
			return nil, err
		}
		for _, seat := range seatsArray {
			reservedSeats[seat] = true
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("error during rows iteration: %v", err)
		return nil, err
	}

	var capacity int
	err = repo.DB.QueryRow(ctx, `
		SELECT capacity
		FROM showtimes
		WHERE id = $1`, id).Scan(&capacity)
	if err != nil {
		log.Printf("error fetching showtime capacity: %v", err)
		return nil, err
	}

	allSeats := generateAllSeats(capacity)
	var availableSeats []string
	for _, seat := range allSeats {
		if !reservedSeats[seat] {
			availableSeats = append(availableSeats, seat)
		}
	}
	return availableSeats, nil
}

func generateAllSeats(capacity int) []string {
	rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	columns := capacity / len(rows)
	var seats []string
	for _, row := range rows {
		for col := 1; col <= columns; col++ {
			seats = append(seats, fmt.Sprintf("%s%d", row, col))
		}
	}
	return seats
}
