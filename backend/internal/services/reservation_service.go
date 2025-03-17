package services

import (
	"context"
	"fmt"
	"movie-system/internal/models"
	"movie-system/internal/repositories"
)

type ReservationService struct {
	repo *repositories.ReservationRepository
}

func NewReservationService(repo *repositories.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) GetReservationsPerMovie(movieID int) ([]models.MovieReservationCount, error) {
	counts, err := s.repo.GetReservationsPerMovie(context.Background(), movieID)
	if err != nil {
		return nil, fmt.Errorf("error getting reservations per movie: %v", err)
	}
	return counts, nil
}
