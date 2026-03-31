package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	authRepo "github.com/adesepriansyah/task-list-timesheet-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Service interface for auth logic.
type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
}

type service struct {
	repo authRepo.AuthRepository
}

// NewService creates a new auth service.
func NewService(repo authRepo.AuthRepository) Service {
	return &service{repo}
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("unauthorized")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("unauthorized")
	}

	// Generate token
	token, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}

	// Set expiration (e.g., 24 hours)
	expiredAt := time.Now().Add(24 * time.Hour)

	if err := s.repo.CreateToken(ctx, user.ID, token, expiredAt); err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteToken(ctx, token)
}

func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
