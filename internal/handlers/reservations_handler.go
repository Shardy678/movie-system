package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"movie-system/internal/models"
	"movie-system/internal/repositories"
	"movie-system/internal/services"
	"net/http"
	"strconv"
	"strings"
)

type ReservationHandler struct {
	Repo        *repositories.ReservationRepository
	AuthService *services.AuthService
}

func NewReservationHandler(repo *repositories.ReservationRepository, authservice *services.AuthService) *ReservationHandler {
	return &ReservationHandler{
		Repo:        repo,
		AuthService: authservice}
}

func (h *ReservationHandler) HandleReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if reservation.UserID == 0 || reservation.ShowtimeID == 0 || len(reservation.Seats) == 0 {
		http.Error(w, "Missing required fields (user_id, showtime_id, or seats)", http.StatusBadRequest)
		return
	}

	if err := h.Repo.ReserveSeat(context.Background(), &reservation); err != nil {
		http.Error(w, fmt.Sprintf("Error creating resevation: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":        "Reservation succesfull",
		"reservation_id": reservation.ID,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func (h *ReservationHandler) HandleCancelReservation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	reservationIDStr := strings.TrimPrefix(r.URL.Path, "/reserve/delete/")
	reservationID, err := strconv.Atoi(reservationIDStr)
	if err != nil {
		http.Error(w, "invalid reservation ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.CancelReservation(context.Background(), reservationID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error cancelling reservation: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Reservation cancelled successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func (h *ReservationHandler) HandleGetReservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := h.AuthService.ExtractUserIDFromJWT(tokenString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
		return
	}

	reservations, err := h.Repo.GetReservations(context.Background(), userID)
	if err != nil {
		http.Error(w, "Error fetching reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}

func (h *ReservationHandler) HandleGetAllReservations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	reservations, err := h.Repo.GetAllReservations(context.Background())
	if err != nil {
		http.Error(w, "Error fetching reservations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reservations)
}
