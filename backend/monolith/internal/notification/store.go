package notification

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"monolith/internal/notification/sqlc"
)

type Store struct {
	queries *sqlc.Queries
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{queries: sqlc.New(db)}
}

func (s *Store) UpsertSubscription(ctx context.Context, params sqlc.UpsertSubscriptionParams) error {
	err := s.queries.UpsertSubscription(ctx, params)
	if err != nil {
		return fmt.Errorf("upsert subscription: %w", err)
	}
	return nil
}

func (s *Store) GetByUserID(ctx context.Context, userID int64) ([]sqlc.PushSubscription, error) {
	subs, err := s.queries.GetSubscriptionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get subscriptions by user id: %w", err)
	}
	return subs, nil
}

func (s *Store) Delete(ctx context.Context, endpoint string) error {
	err := s.queries.DeleteSubscription(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	return nil
}
