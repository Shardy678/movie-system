package main

import (
	"context"
	"fmt"
	"log"
	"movie-system/config"
	"movie-system/internal/handlers"
	"movie-system/internal/middleware"
	"movie-system/internal/repositories"
	"movie-system/internal/seed"
	"movie-system/internal/services"
	"movie-system/routes"
	"net/http"
	"os"
	"time"
)

func main() {

	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("SECRET_KEY is not set in the environment")
	}

	var err error
	config.DB, err = config.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.DB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	config.Ping(ctx)

	if err := seed.SeedDatabase(); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	movieRepo := repositories.NewMovieRepository(config.DB)
	movieHandler := handlers.NewMovieHandler(movieRepo)

	showtimeRepo := repositories.NewShowtimeRepository(config.DB)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeRepo)

	userRepo := repositories.NewUserRepository(config.DB)
	authService := services.NewAuthService(userRepo, jwtSecret)

	authHandler := handlers.NewAuthHandler(authService)

	reservationRepo := repositories.NewReservationRepository(config.DB)
	reservationService := services.NewReservationService(reservationRepo)
	reservationHandler := handlers.NewReservationHandler(reservationRepo, authService, reservationService)

	routes.SetupRoutes(movieHandler, showtimeHandler, authHandler, reservationHandler)

	corsHandler := middleware.CORS(http.DefaultServeMux.ServeHTTP)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(corsHandler)))
}
