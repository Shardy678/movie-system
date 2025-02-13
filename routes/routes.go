package routes

import (
	"movie-system/internal/auth"
	"movie-system/internal/handlers"
	"movie-system/metrics"
	"net/http"
)

func SetupRoutes(mh *handlers.MovieHandler, sh *handlers.ShowtimeHandler, ah *handlers.AuthHandler, rh *handlers.ReservationHandler) {
	// Middleware chain function
	middleware := func(role string, handlerFunc http.HandlerFunc) http.Handler {
		return metrics.RequestCounter(auth.RoleMiddleware(role, handlerFunc))
	}

	// Movie routes
	http.Handle("/movies", middleware("user", http.HandlerFunc(mh.HandleGetMovies)))
	http.Handle("/movies/add", middleware("admin", mh.HandleAddMovie))
	http.Handle("/movies/update/", middleware("admin", mh.HandleUpdateMovie))
	http.Handle("/movies/delete/", middleware("admin", mh.HandleDeleteMovie))

	// User routes
	http.Handle("/auth/signup", middleware("user", http.HandlerFunc(ah.SignUp)))
	http.Handle("/auth/login", middleware("user", http.HandlerFunc(ah.LogIn)))

	// Showtime routes
	http.Handle("/showtimes", middleware("user", sh.HandleGetShowtimes))
	http.Handle("/showtimes/add", middleware("admin", sh.HandleAddShowtime))
	http.Handle("/showtimes/update/", middleware("admin", sh.HandleUpdateShowtime))
	http.Handle("/showtimes/delete/", middleware("admin", sh.HandleDeleteShowtime))
	http.Handle("/showtimes/seats/", middleware("user", sh.HandleGetSeats))

	// Reservation routes
	http.Handle("/reserve/add", middleware("user", rh.HandleReservation))
	http.Handle("/reserve/delete/", middleware("user", rh.HandleCancelReservation))
	http.Handle("/reserve", middleware("user", rh.HandleGetReservations))
	http.Handle("/reserve/all", middleware("admin", rh.HandleGetAllReservations))
	http.Handle("/reserve/movie/", middleware("admin", rh.HandleGetReservationsPerMovie))

	// Revenue routes
	http.Handle("/revenue", middleware("admin", rh.HandleGetTotalRevenue))

	// Metrics routes
	http.Handle("/metrics", metrics.MetricsHandler())
}
