package models

import "time"

type User struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Movie struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	PosterImage string `json:"poster_image"`
}

type Reservation struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	MovieID    uint      `json:"movie_id"`
	ShowtimeID uint      `json:"showtime_id"`
	CreatedAt  time.Time `json:"created_at"`
	Seats      []string  `json:"seats"`
}

type Showtime struct {
	ID        uint      `json:"id"`
	MovieID   uint      `json:"movie_id"`
	StartTime time.Time `json:"start_time"`
	Capacity  uint      `json:"capacity"`
	Reserved  uint      `json:"reserved"`
}

type MovieReservationCount struct {
	MovieID          int    `json:"movie_id"`
	MovieTitle       string `json:"movie_title"`
	ReservationCount int    `json:"reservation_count"`
	SeatCount        int    `json:"seat_count"`
}
