package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/adesepriansyah/task-list-timesheet-be/internal/entity"
)

// AuthRepository interface for auth operations.
type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateToken(ctx context.Context, userID int, token string, expiredAt time.Time) error
	DeleteToken(ctx context.Context, token string) error
	FindToken(ctx context.Context, token string) (*entity.UserToken, error)
}

type authRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new auth repository.
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateToken(ctx context.Context, userID int, token string, expiredAt time.Time) error {
	query := `INSERT INTO user_tokens (user_id, token, expired_at) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiredAt)
	return err
}

func (r *authRepository) DeleteToken(ctx context.Context, token string) error {
	query := `DELETE FROM user_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *authRepository) FindToken(ctx context.Context, token string) (*entity.UserToken, error) {
	var userToken entity.UserToken
	query := `SELECT id, token, user_id, expired_at, created_at, updated_at FROM user_tokens WHERE token = $1`
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userToken.ID, &userToken.Token, &userToken.UserID, &userToken.ExpiredAt, &userToken.CreatedAt, &userToken.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &userToken, nil
}
