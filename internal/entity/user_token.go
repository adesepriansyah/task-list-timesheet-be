package entity

import "time"

// UserToken represents a session token for a user.
type UserToken struct {
	ID        int       `json:"id" db:"id"`
	Token     string    `json:"token" db:"token"`
	UserID    int       `json:"user_id" db:"user_id"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
