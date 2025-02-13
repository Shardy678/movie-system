package routes

import (
	"movie-system/internal/auth"
	"movie-system/internal/handlers"
	"movie-system/metrics"
	"net/http"
)

func SetupRoutes(movieHandler *handlers.MovieHandler, showtimeHandler *handlers.ShowtimeHandler, authHandler *handlers.AuthHandler, reservationHandler *handlers.ReservationHandler) {
	// Movie routes
	http.HandleFunc("/movies", metrics.RequestCounter(movieHandler.HandleGetMovies))
	http.Handle("/movies/add", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(movieHandler.HandleAddMovie))))
	http.Handle("/movies/update/", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(movieHandler.HandleUpdateMovie))))
	http.Handle("/movies/delete/", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(movieHandler.HandleDeleteMovie))))

	// User routes
	http.HandleFunc("/auth/signup", metrics.RequestCounter(authHandler.SignUp))
	http.HandleFunc("/auth/login", metrics.RequestCounter(authHandler.LogIn))

	// Showtime routes
	http.HandleFunc("/showtimes", metrics.RequestCounter(showtimeHandler.HandleGetShowtimes))
	http.Handle("/showtimes/add", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(showtimeHandler.HandleAddShowtime))))
	http.Handle("/showtimes/update/", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(showtimeHandler.HandleUpdateShowtime))))
	http.Handle("/showtimes/delete/", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(showtimeHandler.HandleDeleteShowtime))))
	http.HandleFunc("/showtimes/seats/", metrics.RequestCounter(showtimeHandler.HandleGetSeats))

	// Reservation routes
	http.Handle("/reserve/add", auth.RoleMiddleware("user", http.HandlerFunc(metrics.RequestCounter(reservationHandler.HandleReservation))))
	http.HandleFunc("/reserve/delete/", metrics.RequestCounter(reservationHandler.HandleCancelReservation))
	http.HandleFunc("/reserve", metrics.RequestCounter(reservationHandler.HandleGetReservations))
	http.Handle("/reserve/all", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(reservationHandler.HandleGetAllReservations))))
	http.Handle("/reserve/movie/", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(reservationHandler.HandleGetReservationsPerMovie))))

	// Revenue routes
	http.Handle("/revenue", auth.RoleMiddleware("admin", http.HandlerFunc(metrics.RequestCounter(reservationHandler.HandleGetTotalRevenue))))

	// Metrics routes
	http.Handle("/metrics", metrics.MetricsHandler())
}
