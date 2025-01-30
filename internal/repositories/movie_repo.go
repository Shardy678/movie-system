package repositories

import (
	"context"
	"log"
	"movie-system/config"
	"movie-system/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepository struct {
	DB *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{DB: db}
}

func (repo *MovieRepository) InsertMovie(ctx context.Context, movie *models.Movie) error {
	_, err := repo.DB.Exec(ctx, `
		INSERT INTO movies (title, description, genre, poster_image)
		VALUES ($1, $2, $3, $4)`, movie.Title, movie.Description, movie.Genre, movie.PosterImage)
	return err
}

func (repo *MovieRepository) GetMovies(ctx context.Context, genre string) ([]models.Movie, error) {
	var rows pgx.Rows
	var err error

	if genre != "" {
		rows, err = repo.DB.Query(ctx, `
		SELECT id, title, description, genre, poster_image
		FROM movies
		WHERE genre ILIKE $1`, "%"+genre+"%")
	} else {
		rows, err = config.DB.Query(ctx, `
		SELECT id, title, description, genre, poster_image
		FROM movies`)
	}

	if err != nil {
		log.Printf("error fetching movies: %v", err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Genre, &movie.PosterImage); err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	return movies, nil
}

func (repo *MovieRepository) UpdateMovie(ctx context.Context, id int, movie *models.Movie) error {
	_, err := repo.DB.Exec(ctx, `
		UPDATE movies
		SET title = $1, description = $2, genre = $3, poster_image = $4
		WHERE id = $5`,
		movie.Title, movie.Description, movie.Genre, movie.PosterImage, id)
	return err
}

func (repo *MovieRepository) DeleteMovie(ctx context.Context, id int) error {
	_, err := repo.DB.Exec(ctx, "DELETE FROM movies WHERE id = $1", id)
	return err
}
