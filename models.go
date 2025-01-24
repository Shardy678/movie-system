package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           uint          `json:"id"`
	Username     string        `gorm:"unique;notnull" json:"username"`
	PasswordHash string        `json:"password_hash"`
	Role         string        `gorm:"default:'user'" json:"role"`
	Reservations []Reservation `gorm:"foreignKey:UserID" json:"reservations"`
}

type Movie struct {
	gorm.Model
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Showtime    time.Time `json:"showtime"`
}

type Reservation struct {
	gorm.Model
	ID      uint `json:"id"`
	UserID  uint `json:"user_id"`
	MovieID uint `json:"movie_id"`
	Seats   uint `json:"seats"`
}
