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
