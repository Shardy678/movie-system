package repositories

import (
	"context"
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
