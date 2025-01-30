package services

import (
	"context"
	"movie-system/internal/models"
	"movie-system/internal/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	repo      *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repositories.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret}
}

func (s *AuthService) SignUp(ctx context.Context, user *models.User) error {
	return s.repo.SignUp(ctx, user)
}

func (s *AuthService) LogIn(ctx context.Context, username, password string) (string, error) {
	role, err := s.repo.AuthenticateUser(ctx, username, password)
	if err != nil {
		return "", err
	}

	token, err := s.GenerateJWT(username, role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GenerateJWT(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
