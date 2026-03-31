package auth

import (
	"context"
	"testing"
	"time"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
	"github.com/golang-jwt/jwt/v5"
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
	svc := NewService(repo, "testsecret")

	t.Run("Success", func(t *testing.T) {
		token, err := svc.Login(context.Background(), "test@example.com", "password123")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected non-empty token")
		}
	})

	t.Run("WrongPassword", func(t *testing.T) {
		_, err := svc.Login(context.Background(), "test@example.com", "wrong")
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})
}

func TestRegister(t *testing.T) {
	repo := &mockRepo{
		users: make(map[string]*entity.User),
	}
	svc := NewService(repo, "testsecret")

	t.Run("Success", func(t *testing.T) {
		err := svc.Register(context.Background(), "Test User", "register@example.com", "password123")
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})
}

func TestGetUserInfo(t *testing.T) {
	secret := "testsecret"
	repo := &mockRepo{
		users: map[string]*entity.User{
			"info@example.com": {ID: 1, Name: "Info User", Email: "info@example.com"},
		},
		tokens: make(map[string]*entity.UserToken),
	}
	svc := NewService(repo, secret)

	// Helper to generate JWT
	genToken := func(id int, exp time.Time) string {
		tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": id,
			"exp": exp.Unix(),
		})
		s, _ := tkn.SignedString([]byte(secret))
		return s
	}

	validToken := genToken(1, time.Now().Add(time.Hour))
	repo.tokens[validToken] = &entity.UserToken{UserID: 1, Token: validToken, ExpiredAt: time.Now().Add(time.Hour)}

	expiredToken := genToken(1, time.Now().Add(-time.Hour))
	repo.tokens[expiredToken] = &entity.UserToken{UserID: 1, Token: expiredToken, ExpiredAt: time.Now().Add(-time.Hour)}

	t.Run("Success", func(t *testing.T) {
		user, _, err := svc.GetUserInfo(context.Background(), validToken)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user == nil || user.Email != "info@example.com" {
			t.Error("expected correct user info")
		}
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		_, _, err := svc.GetUserInfo(context.Background(), expiredToken)
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		_, _, err := svc.GetUserInfo(context.Background(), "invalid-token-format")
		if err == nil || err.Error() != "unauthorized" {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})
}
