package repositories

import (
	"context"
	"fmt"
	"log"
	"movie-system/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) SignUp(ctx context.Context, user *models.User) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	// Insert user into DB
	_, err = repo.db.Exec(ctx, `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, $3)`, user.Username, string(hashedPassword), user.Role)

	// Log SQL error if it occurs
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}

	// Verify user was inserted
	var count int
	err = repo.db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE username = $1`, user.Username).Scan(&count)
	if err != nil {
		log.Printf("Error verifying user insertion: %v", err)
		return err
	}

	if count == 0 {
		log.Println("User was not inserted into the database!")
		return fmt.Errorf("user was not inserted into the database")
	}

	log.Printf("User %s successfully inserted into the database", user.Username)
	return nil
}

func (repo *UserRepository) AuthenticateUser(ctx context.Context, username, password string) (string, error) {
	var storedHash, role string

	err := repo.db.QueryRow(ctx, `
		SELECT password_hash, role FROM users WHERE username = $1`, username).
		Scan(&storedHash, &role)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		return "", err
	}

	return role, nil
}

func (repo *UserRepository) GetUserID(ctx context.Context, username string) (int, error) {
	var userID int
	err := repo.db.QueryRow(ctx, "SELECT id from users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}
