package auth

import (
	"context"
	"errors"
	"time"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
	authRepo "github.com/adesepriansyah/task-list-timesheet-be/internal/repository"
	"github.com/golang-jwt/jwt/v5"
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
	repo      authRepo.AuthRepository
	jwtSecret []byte
}

// NewService creates a new auth service.
func NewService(repo authRepo.AuthRepository, jwtSecret string) Service {
	return &service{repo, []byte(jwtSecret)}
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

	// Generate token (JWT)
	expiredAt := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expiredAt.Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	if err := s.repo.CreateToken(ctx, user.ID, tokenString, expiredAt); err != nil {
		return "", err
	}

	return tokenString, nil
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

func (s *service) GetUserInfo(ctx context.Context, tokenString string) (*entity.User, time.Time, error) {
	// 1. Parse and validate JWT signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, time.Time{}, errors.New("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, time.Time{}, errors.New("unauthorized")
	}

	// 2. Double check with DB (for revocation/logout)
	userToken, err := s.repo.FindToken(ctx, tokenString)
	if err != nil || userToken == nil {
		return nil, time.Time{}, errors.New("unauthorized")
	}

	// 3. Get user info
	userID := int(claims["sub"].(float64))
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, time.Time{}, errors.New("unauthorized")
	}

	return user, userToken.ExpiredAt, nil
}
