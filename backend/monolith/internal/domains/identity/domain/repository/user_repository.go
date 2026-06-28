package repository

import (
	"context"

	"monolith/internal/domains/identity/domain/value"
)

type UserRepository interface {
	CreateUser(ctx context.Context) (value.UserID, error)
	// GetUserByID(ctx context.Context, userID value.UserID) (entity.User, error)
	GetUserRolesByID(ctx context.Context, userID value.UserID) ([]int32, error)
	DeleteUser(ctx context.Context, userID value.UserID) error
}
