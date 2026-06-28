package repository

import (
	"context"

	"monolith/internal/domains/identity/domain/value"
)

type UserRolesRepository interface {
	AddRole(ctx context.Context, userID value.UserID, role string) error
}
