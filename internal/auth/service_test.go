package auth

import (
	"context"
	"testing"
	"time"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	users  map[string]*entity.User
	tokens map[string]*entity.UserToken
}

func (m *mockRepo) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.users[email], nil
}

func (m *mockRepo) CreateToken(ctx context.Context, userID int, token string, expiredAt time.Time) error {
	m.tokens[token] = &entity.UserToken{UserID: userID, Token: token, ExpiredAt: expiredAt}
	return nil
}

func (m *mockRepo) DeleteToken(ctx context.Context, token string) error {
	delete(m.tokens, token)
	return nil
}

func (m *mockRepo) CreateUser(ctx context.Context, name, email, hashedPassword string) error {
	m.users[email] = &entity.User{ID: len(m.users) + 1, Name: name, Email: email, Password: hashedPassword}
	return nil
}

func (m *mockRepo) FindUserByID(ctx context.Context, id int) (*entity.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockRepo) FindToken(ctx context.Context, token string) (*entity.UserToken, error) {
	return m.tokens[token], nil
}

func TestLogin(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	repo := &mockRepo{
		users: map[string]*entity.User{
			"test@example.com": {ID: 1, Email: "test@example.com", Password: string(hashedPassword)},
		},
		tokens: make(map[string]*entity.UserToken),
	}
	service := NewService(repo)

	t.Run("Success", func(t *testing.T) {
		token, err := service.Login(context.Background(), "test@example.com", "password123")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
	})

	t.Run("WrongPassword", func(t *testing.T) {
		_, err := service.Login(context.Background(), "test@example.com", "wrong")
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		_, err := service.Login(context.Background(), "unknown@example.com", "password123")
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})
}

func TestRegister(t *testing.T) {
	repo := &mockRepo{
		users: make(map[string]*entity.User),
	}
	service := NewService(repo)

	t.Run("Success", func(t *testing.T) {
		err := service.Register(context.Background(), "Test User", "register@example.com", "password123")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		user, _ := repo.FindUserByEmail(context.Background(), "register@example.com")
		if user == nil || user.Name != "Test User" {
			t.Error("expected user to be created with correct name")
		}
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		err := service.Register(context.Background(), "Another User", "register@example.com", "password123")
		if err == nil || err.Error() != "email already registered" {
			t.Errorf("expected email already registered error, got %v", err)
		}
	})
}

func TestGetUserInfo(t *testing.T) {
	repo := &mockRepo{
		users: map[string]*entity.User{
			"info@example.com": {ID: 1, Name: "Info User", Email: "info@example.com"},
		},
		tokens: map[string]*entity.UserToken{
			"valid-token":   {UserID: 1, Token: "valid-token", ExpiredAt: time.Now().Add(time.Hour)},
			"expired-token": {UserID: 1, Token: "expired-token", ExpiredAt: time.Now().Add(-time.Hour)},
		},
	}
	service := NewService(repo)

	t.Run("Success", func(t *testing.T) {
		user, expiredAt, err := service.GetUserInfo(context.Background(), "valid-token")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil || user.Email != "info@example.com" {
			t.Error("expected correct user info")
		}
		if expiredAt.IsZero() {
			t.Error("expected non-zero expiredAt")
		}
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		_, _, err := service.GetUserInfo(context.Background(), "expired-token")
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, _, err := service.GetUserInfo(context.Background(), "invalid-token")
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})
}
