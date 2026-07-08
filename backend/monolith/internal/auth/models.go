package auth

import "time"

type User struct {
	ID           int64
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	Phone        *string
	Email        *string
	ProfileImage *string
}

type Session struct {
	UserID    int64     `json:"user_id"`
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Roles     []string  `json:"roles"`
}
