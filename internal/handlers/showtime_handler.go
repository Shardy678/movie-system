package handlers

import (
	"context"
	"encoding/json"
	"movie-system/internal/models"
	"movie-system/internal/repositories"
	"net/http"
	"strconv"
	"strings"
)

type ShowtimeHandler struct {
	Repo *repositories.ShowtimeRepository
}

func NewShowtimeHandler(repo *repositories.ShowtimeRepository) *ShowtimeHandler {
	return &ShowtimeHandler{Repo: repo}
}

func (h *ShowtimeHandler) HandleAddShowtime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var showtime models.Showtime
	if err := json.NewDecoder(r.Body).Decode(&showtime); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	err := h.Repo.InsertShowtime(context.Background(), &showtime)
	if err != nil {
		http.Error(w, "Failed to add showtime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Showtime added successfully"})
}

func (h *ShowtimeHandler) HandleGetShowtimes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	showtimes, err := h.Repo.GetShowtimes(context.Background())
	if err != nil {
		http.Error(w, "Failed to fetch showtimes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(showtimes)
}

func (h *ShowtimeHandler) HandleUpdateShowtime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/showtimes/update/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid showtime ID", http.StatusBadRequest)
		return
	}

	var showtime models.Showtime
	if err := json.NewDecoder(r.Body).Decode(&showtime); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = h.Repo.UpdateShowtime(context.Background(), id, &showtime)
	if err != nil {
		http.Error(w, "Failed to update showtime", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Showtime updated successfully"})
}

func (h *ShowtimeHandler) HandleDeleteShowtime(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/showtimes/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid showtime ID", http.StatusBadRequest)
		return
	}

	err = h.Repo.DeleteShowtime(context.Background(), id)
	if err != nil {
		http.Error(w, "Failed to delete showtime", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Showtime deleted successfully"})
}
