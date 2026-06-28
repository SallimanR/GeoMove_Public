package port

import (
	"context"

	"monolith/internal/domains/identity/domain/value"
)

type SessionService interface {
	Validate(ctx context.Context, token string) (value.UserID, error)
	GetUserRoles(ctx context.Context, userID value.UserID) ([]string, error)
}
