package repositories

import (
	"context"
	"fmt"
	"log"
	"movie-system/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ShowtimeRepository struct {
	DB *pgxpool.Pool
}

func NewShowtimeRepository(db *pgxpool.Pool) *ShowtimeRepository {
	return &ShowtimeRepository{DB: db}
}

func (repo *ShowtimeRepository) InsertShowtime(ctx context.Context, showtime *models.Showtime) error {
	query := `
		INSERT INTO showtimes (movie_id, start_time, capacity, reserved)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := repo.DB.QueryRow(ctx, query,
		showtime.MovieID,
		showtime.StartTime,
		showtime.Capacity,
		showtime.Reserved,
	).Scan(&showtime.ID)
	if err != nil {
		log.Printf("error inserting showtime: %v", err)
		return err
	}
	log.Printf("inserted showtime with id: %d", showtime.ID)
	return nil
}



func (repo *ShowtimeRepository) GetShowtimes(ctx context.Context) ([]models.Showtime, error) {
	rows, err := repo.DB.Query(ctx, `
		SELECT id, movie_id, start_time, capacity, reserved 
		FROM showtimes`)
	if err != nil {
		log.Printf("error fetching showtimes: %v", err)
		return nil, fmt.Errorf("error fetching showtimes: %w", err)
	}
	defer rows.Close()

	var showtimes []models.Showtime
	for rows.Next() {
		var showtime models.Showtime
		if err := rows.Scan(
			&showtime.ID,
			&showtime.MovieID,
			&showtime.StartTime,
			&showtime.Capacity,
			&showtime.Reserved,
		); err != nil {
			log.Printf("error scanning showtime: %v", err)
			return nil, fmt.Errorf("error scanning showtime: %w", err)
		}
		showtimes = append(showtimes, showtime)
	}

	if err := rows.Err(); err != nil {
		log.Printf("error during rows iteration: %v", err)
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return showtimes, nil
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
    const maxRows = 10
    const maxSeatsPerRow = 10

    rows := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
    
    if capacity > maxRows*maxSeatsPerRow {
        capacity = maxRows * maxSeatsPerRow
    }
    
    var seats []string
    row := 0
    col := 1
    
    for i := 0; i < capacity; i++ {
        seats = append(seats, fmt.Sprintf("%s%d", rows[row], col))
        col++
        if col > maxSeatsPerRow {
            col = 1
            row++
        }
    }
    return seats
}

