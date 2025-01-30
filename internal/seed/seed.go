package seed

import (
	"context"
	"fmt"
	"movie-system/config"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() error {
	ctx := context.Background()

	adminPassword := "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	_, err = config.DB.Exec(ctx, `
			INSERT INTO users (username, password_hash, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (username) DO NOTHING`,
		"admin", string(hashedPassword), "admin", time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}
	movies := []struct {
		Title       string
		Description string
		Genre       string
		PosterImage string
	}{
		{
			Title:       "Inception",
			Description: "A thief who steals corporate secrets through the use of dream-sharing technology.",
			Genre:       "Sci-Fi",
			PosterImage: "https://m.media-amazon.com/images/M/MV5BMjAxMzY3NjcxNF5BMl5BanBnXkFtZTcwNTI5OTM0Mw@@._V1_FMjpg_UX1000_.jpg",
		},
		{
			Title:       "The Dark Knight",
			Description: "When the menace known as the Joker emerges, Batman must confront chaos.",
			Genre:       "Action",
			PosterImage: "https://m.media-amazon.com/images/M/MV5BMTMxNTMwODM0NF5BMl5BanBnXkFtZTcwODAyMTk2Mw@@._V1_FMjpg_UX1000_.jpg",
		},
		{
			Title:       "Interstellar",
			Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
			Genre:       "Sci-Fi",
			PosterImage: "https://m.media-amazon.com/images/M/MV5BYzdjMDAxZGItMjI2My00ODA1LTlkNzItOWFjMDU5ZDJlYWY3XkEyXkFqcGc@._V1_FMjpg_UX1000_.jpg",
		},
	}
	for _, movie := range movies {
		_, err := config.DB.Exec(ctx, `
				INSERT INTO movies (title, description, genre, poster_image)
				VALUES ($1,$2,$3,$4)
				ON CONFLICT (title) DO NOTHING`,
			movie.Title, movie.Description, movie.Genre, movie.PosterImage)
		if err != nil {
			return fmt.Errorf("failed to insert movie: %w", err)
		}
	}
	// showtimes := []struct {
	// 	MovieID   uint
	// 	StartTime time.Time
	// 	Capacity  uint
	// 	Reserved  uint
	// }{
	// 	{
	// 		MovieID:   1, // Inception
	// 		StartTime: time.Now().Add(24 * time.Hour),
	// 		Capacity:  100,
	// 		Reserved:  0,
	// 	},
	// 	{
	// 		MovieID:   2, // The Dark Knight
	// 		StartTime: time.Now().Add(48 * time.Hour),
	// 		Capacity:  150,
	// 		Reserved:  0,
	// 	},
	// 	{
	// 		MovieID:   3, // Interstellar
	// 		StartTime: time.Now().Add(72 * time.Hour),
	// 		Capacity:  200,
	// 		Reserved:  0,
	// 	},
	// }

	// for _, showtime := range showtimes {
	// 	_, err := db.Exec(ctx, `
	// 		INSERT INTO showtimes (movie_id, start_time, capacity, reserved)
	// 		VALUES ($1, $2, $3, $4)`,
	// 		showtime.MovieID, showtime.StartTime, showtime.Capacity, showtime.Reserved,
	// 	)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to insert showtime: %w", err)
	// 	}
	// }

	return nil
}
