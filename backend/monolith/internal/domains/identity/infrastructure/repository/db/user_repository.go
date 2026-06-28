package db

import (
	"context"

	"monolith/internal/domains/identity/domain/repository"
	"monolith/internal/domains/identity/domain/value"
	"monolith/internal/domains/identity/infrastructure/repository/sqlc"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepositry(queries *sqlc.Queries) repository.UserRepository {
	return &UserRepository{
		queries: queries,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context) (value.UserID, error) {
	userID, err := r.queries.CreateUser(ctx)
	if err != nil {
		return 0, err
	}
	return value.UserID(userID), nil
}

func (r *UserRepository) GetUserRolesByID(ctx context.Context, userID value.UserID) ([]int32, error) {
	result, err := r.queries.GetUserRolesByID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID value.UserID) error {
	err := r.queries.DeleteUser(ctx, int64(userID))
	if err != nil {
		return err
	}
	return err
}
