package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"movie-system/internal/models"
	"movie-system/internal/repositories"
	"net/http"
)

type ReservationHandler struct {
	Repo *repositories.ReservationRepository
}

func NewReservationHandler(repo *repositories.ReservationRepository) *ReservationHandler {
	return &ReservationHandler{Repo: repo}
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
