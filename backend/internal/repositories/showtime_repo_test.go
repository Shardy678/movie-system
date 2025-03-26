package repositories

import (
	"context"
	"movie-system/internal/models"
	"movie-system/test"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShowtimeRepository(t *testing.T) {
	db, err := test.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer db.Close()

	err = test.ClearTestDB(db)
	if err != nil {
		t.Fatalf("Failed to clear test database: %v", err)
	}

	repo := NewShowtimeRepository(db)
	ctx := context.Background()

	assert.NoError(t, err)

	t.Run("InsertShowtime", func(t *testing.T) {
		err := test.ClearTestDB(db)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO movies (id, title, description, genre, poster_image) VALUES
			(1, 'Test Movie 1', 'Test Description 1', 'Action', 'poster1.jpg')
		`)
		assert.NoError(t, err)

		showtime := &models.Showtime{
			MovieID:   1,
			StartTime: time.Now(),
			Capacity:  100,
			Reserved:  0,
		}

		err = repo.InsertShowtime(ctx, showtime)
		assert.NoError(t, err)
	})

	t.Run("GetShowtimes", func(t *testing.T) {
		err := test.ClearTestDB(db)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO movies (id, title, description, genre, poster_image) VALUES
			(1, 'Test Movie 1', 'Test Description 1', 'Action', 'poster1.jpg')
		`)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO showtimes (id, movie_id, start_time, capacity, reserved) VALUES
			(1, 1, '2024-01-01 10:00:00', 100, 0)
		`)
		assert.NoError(t, err)

		showtimes, err := repo.GetShowtimes(ctx)
		assert.NoError(t, err)
		assert.Len(t, showtimes, 1)
		assert.Equal(t, uint(100), showtimes[0].Capacity)
	})

	t.Run("UpdateShowtime", func(t *testing.T) {
		err := test.ClearTestDB(db)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO movies (id, title, description, genre, poster_image) VALUES
			(1, 'Test Movie 1', 'Test Description 1', 'Action', 'poster1.jpg'),
			(2, 'Test Movie 2', 'Test Description 2', 'Comedy', 'poster2.jpg')
		`)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO showtimes (id, movie_id, start_time, capacity, reserved) VALUES
			(1, 1, '2024-01-01 10:00:00', 100, 0)
		`)
		assert.NoError(t, err)

		updatedShowtime := &models.Showtime{
			MovieID:   2,
			StartTime: time.Now(),
			Capacity:  150,
			Reserved:  10,
		}

		err = repo.UpdateShowtime(ctx, 1, updatedShowtime)
		assert.NoError(t, err)

		showtimes, err := repo.GetShowtimes(ctx)
		assert.NoError(t, err)
		assert.Equal(t, uint(150), showtimes[0].Capacity)
		assert.Equal(t, uint(10), showtimes[0].Reserved)
	})

	t.Run("DeleteShowtime", func(t *testing.T) {
		err := test.ClearTestDB(db)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO movies (id, title, description, genre, poster_image) VALUES
			(1, 'Test Movie 1', 'Test Description 1', 'Action', 'poster1.jpg')
		`)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO showtimes (id, movie_id, start_time, capacity, reserved) VALUES
			(1, 1, '2024-01-01 10:00:00', 100, 0)
		`)
		assert.NoError(t, err)

		showtimes, err := repo.GetShowtimes(ctx)
		assert.NoError(t, err)
		assert.Len(t, showtimes, 1)

		err = repo.DeleteShowtime(ctx, 1)
		assert.NoError(t, err)

		showtimes, err = repo.GetShowtimes(ctx)
		assert.NoError(t, err)
		assert.Len(t, showtimes, 0)
	})

	t.Run("GetAvailableSeats", func(t *testing.T) {
		err := test.ClearTestDB(db)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO movies (id, title, description, genre, poster_image) VALUES
			(1, 'Test Movie 1', 'Test Description 1', 'Action', 'poster1.jpg')
		`)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO showtimes (id, movie_id, start_time, capacity, reserved) VALUES
			(1, 1, '2024-01-01 10:00:00', 100, 0)
		`)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO users (id, username, password_hash, role) VALUES
			(1, 'testuser1', 'password1', 'user')
		`)
		assert.NoError(t, err)

		_, err = db.Exec(ctx, `
			INSERT INTO reservations (user_id, movie_id, showtime_id, seats) VALUES
			(1, 1, 1, ARRAY['A1', 'A2']::text[])
		`)
		assert.NoError(t, err)

		availableSeats, err := repo.GetAvailableSeats(ctx, 1)
		assert.NoError(t, err)

		assert.NotContains(t, availableSeats, "A1")
		assert.NotContains(t, availableSeats, "A2")
		assert.Contains(t, availableSeats, "A3")
		assert.Contains(t, availableSeats, "B1")
	})
}
