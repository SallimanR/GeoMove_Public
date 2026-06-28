package db

import (
	"context"
	"time"

	"monolith/internal/domains/identity/domain/repository"
	"monolith/internal/domains/identity/domain/value"
	"monolith/internal/domains/identity/infrastructure/repository/sqlc"
)

type SessionRepository struct {
	queries *sqlc.Queries
}

func NewSessionRepository(queries *sqlc.Queries) repository.SessionRepository {
	return &SessionRepository{
		queries: queries,
	}
}

func (r *SessionRepository) CreateSession(ctx context.Context, userID value.UserID, token string, expiresAt time.Time) error {
	err := r.queries.CreateSession(ctx, int64(userID))
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, token string) (value.UserID, time.Time, error) {
}

// FIXME: hash instead of token
func (r *SessionRepository) GetSessionIDByToken(ctx context.Context, token string) (value.UserID, time.Time, error) {
	r.queries.GetSessionIDByTokenHash(ctx, token)
}

func (r *SessionRepository) DeleteSession(ctx context.Context, token string) error {}

func (r *SessionRepository) CreateAccessToken(ctx context.Context, userID value.UserID, token string, expiresAt time.Time) error {
}
