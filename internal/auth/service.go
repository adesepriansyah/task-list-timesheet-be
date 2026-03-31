package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
	authRepo "github.com/adesepriansyah/task-list-timesheet-be/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// Service interface for auth logic.
type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, token string) error
	Register(ctx context.Context, name, email, password string) error
	GetUserInfo(ctx context.Context, token string) (*entity.User, time.Time, error)
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

func (s *service) Register(ctx context.Context, name, email, password string) error {
	// Check if user already exists
	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user != nil {
		return errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, name, email, string(hashedPassword))
}

func (s *service) GetUserInfo(ctx context.Context, token string) (*entity.User, time.Time, error) {
	// Find and validate token
	userToken, err := s.repo.FindToken(ctx, token)
	if err != nil {
		return nil, time.Time{}, err
	}
	if userToken == nil {
		return nil, time.Time{}, errors.New("unauthorized")
	}

	// Check if token is expired
	if time.Now().After(userToken.ExpiredAt) {
		return nil, time.Time{}, errors.New("unauthorized")
	}

	// Get user info
	user, err := s.repo.FindUserByID(ctx, userToken.UserID)
	if err != nil {
		return nil, time.Time{}, err
	}
	if user == nil {
		return nil, time.Time{}, errors.New("unauthorized")
	}

	return user, userToken.ExpiredAt, nil
}

func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
