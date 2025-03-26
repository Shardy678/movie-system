package services

import (
	"context"
	"fmt"
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

func (s *AuthService) ExtractUserIDFromJWT(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %v", err)
	}

	var username string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var ok bool
		username, ok = claims["username"].(string)
		if !ok {
			return 0, fmt.Errorf("username not found in token")
		}

	} else {
		return 0, fmt.Errorf("invalid token or token claims are malformed")
	}

	userID, err := s.repo.GetUserID(context.Background(), username)
	if err != nil {
		return 0, fmt.Errorf("failed to get user ID: %w", err)
	}
	return userID, nil
}
