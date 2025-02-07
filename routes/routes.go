package routes

import (
	"movie-system/internal/auth"
	"movie-system/internal/handlers"
	"net/http"
)

func SetupRoutes(movieHandler *handlers.MovieHandler, showtimeHandler *handlers.ShowtimeHandler, authHandler *handlers.AuthHandler, reservationHandler *handlers.ReservationHandler) {
	// Movie routes
	http.HandleFunc("/movies", movieHandler.HandleGetMovies)                                                       // GET /movies
	http.Handle("/movies/add", auth.RoleMiddleware("admin", http.HandlerFunc(movieHandler.HandleAddMovie)))        // POST /movies/add
	http.Handle("/movies/update/", auth.RoleMiddleware("admin", http.HandlerFunc(movieHandler.HandleUpdateMovie))) // PUT /movies/update/{id}
	http.Handle("/movies/delete/", auth.RoleMiddleware("admin", http.HandlerFunc(movieHandler.HandleDeleteMovie))) // DELETE /movies/delete/{id}

	// User routes
	http.HandleFunc("/auth/signup", authHandler.SignUp) // POST /auth/signup
	http.HandleFunc("/auth/login", authHandler.LogIn)   // POST /auth/login

	// Showtime routes
	http.HandleFunc("/showtimes", showtimeHandler.HandleGetShowtimes)                                                       // GET /showtimes
	http.Handle("/showtimes/add", auth.RoleMiddleware("admin", http.HandlerFunc(showtimeHandler.HandleAddShowtime)))        // POST /showtimes/add
	http.Handle("/showtimes/update/", auth.RoleMiddleware("admin", http.HandlerFunc(showtimeHandler.HandleUpdateShowtime))) // PUT /showtimes/update/{id}
	http.Handle("/showtimes/delete/", auth.RoleMiddleware("admin", http.HandlerFunc(showtimeHandler.HandleDeleteShowtime))) // DELETE /showtimes/delete/{id}
	http.HandleFunc("/showtimes/seats/", showtimeHandler.HandleGetSeats)                                                    // GET /showtimes/seats/{id}

	// Reservation routes
	http.Handle("/reserve/add", auth.RoleMiddleware("user", http.HandlerFunc(reservationHandler.HandleReservation)))                 // POST /reserve/add
	http.HandleFunc("/reserve/delete/", reservationHandler.HandleCancelReservation)                                                  // DELETE /reserve/delete
	http.HandleFunc("/reserve", reservationHandler.HandleGetReservations)                                                            // GET /reserve
	http.Handle("/reserve/all", auth.RoleMiddleware("admin", http.HandlerFunc(reservationHandler.HandleGetAllReservations)))         // GET /reserve/all
	http.Handle("/reserve/movie/", auth.RoleMiddleware("admin", http.HandlerFunc(reservationHandler.HandleGetReservationsPerMovie))) // GET /reserve/movie

	// Revenue routes
	http.Handle("/revenue", auth.RoleMiddleware("admin", http.HandlerFunc(reservationHandler.HandleGetTotalRevenue))) // GET /revenue
}
