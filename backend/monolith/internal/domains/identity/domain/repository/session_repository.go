package repository

import (
	"context"
	"time"

	"monolith/internal/domains/identity/domain/value"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, userID value.UserID, token string, expiresAt time.Time) error
	GetSessionByID(ctx context.Context, token string) (value.UserID, time.Time, error)
	GetSessionIDByTokenHash(ctx context.Context, token string) (value.UserID, time.Time, error)
	DeleteSession(ctx context.Context, token string) error
}
