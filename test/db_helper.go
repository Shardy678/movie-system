package test

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupTestDB() (*pgxpool.Pool, error) {
	connStr := "postgres://nosweat:password@localhost:5432/movie_system_test"

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	log.Println("Successfully connected to the test database")
	return db, nil
}

func ClearTestDB(db *pgxpool.Pool) error {
	tables := []string{
		"reservations",
		"showtimes",
		"movies",
		"users",
	}

	ctx := context.Background()

	_, err := db.Exec(ctx, "SET CONSTRAINTS ALL DEFERRED")
	if err != nil {
		return fmt.Errorf("failed to defer constraints: %w", err)
	}

	for _, table := range tables {
		_, err := db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	_, err = db.Exec(ctx, "SET CONSTRAINTS ALL IMMEDIATE")
	if err != nil {
		return fmt.Errorf("failed to restore constraints: %w", err)
	}

	return nil
}
