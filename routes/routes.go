package routes

import (
	"movie-system/internal/auth"
	"movie-system/internal/handlers"
	"net/http"
)

func SetupRoutes(movieHandler *handlers.MovieHandler, showtimeHandler *handlers.ShowtimeHandler, authHandler *handlers.AuthHandler, reservationHandler *handlers.ReservationHandler) {
	// Movie routes
	http.HandleFunc("/movies", movieHandler.HandleGetMovies)
	http.Handle("/movies/add", auth.RoleMiddleware("admin", http.HandlerFunc(movieHandler.HandleAddMovie)))
	http.Handle("/movies/update/", auth.RoleMiddleware("admin", http.HandlerFunc(movieHandler.HandleUpdateMovie)))
	http.Handle("/movies/delete/", auth.RoleMiddleware("admin", http.HandlerFunc(movieHandler.HandleDeleteMovie)))

	// User routes
	http.HandleFunc("/auth/signup", authHandler.SignUp)
	http.HandleFunc("/auth/login", authHandler.LogIn)

	// Showtime routes
	http.HandleFunc("/showtimes", showtimeHandler.HandleGetShowtimes)
	http.Handle("/showtimes/add", auth.RoleMiddleware("admin", http.HandlerFunc(showtimeHandler.HandleAddShowtime)))
	http.Handle("/showtimes/update/", auth.RoleMiddleware("admin", http.HandlerFunc(showtimeHandler.HandleUpdateShowtime)))
	http.Handle("/showtimes/delete/", auth.RoleMiddleware("admin", http.HandlerFunc(showtimeHandler.HandleDeleteShowtime)))
	http.HandleFunc("/showtimes/seats/", showtimeHandler.HandleGetSeats)

	// Reservation routes
	http.Handle("/reserve/add", auth.RoleMiddleware("user", http.HandlerFunc(reservationHandler.HandleReservation)))
	http.HandleFunc("/reserve/delete/", reservationHandler.HandleCancelReservation)
	http.HandleFunc("/reserve", reservationHandler.HandleGetReservations)
	http.Handle("/reserve/all", auth.RoleMiddleware("admin", http.HandlerFunc(reservationHandler.HandleGetAllReservations)))
	http.Handle("/reserve/movie/", auth.RoleMiddleware("admin", http.HandlerFunc(reservationHandler.HandleGetReservationsPerMovie)))

	// Revenue routes
	http.Handle("/revenue", auth.RoleMiddleware("admin", http.HandlerFunc(reservationHandler.HandleGetTotalRevenue)))
}
