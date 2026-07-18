package query

import (
	"context"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
)

type ListOrdersByUserHandler struct {
	repo repository.OrderRepository
}

func NewListOrdersByUserHandler(repo repository.OrderRepository) *ListOrdersByUserHandler {
	return &ListOrdersByUserHandler{repo: repo}
}

func (h *ListOrdersByUserHandler) Handle(ctx context.Context, userID int64, role string) ([]entity.Order, error) {
	switch role {
	case "customer":
		return h.repo.ListOrdersByCustomer(ctx, userID)
	case "driver":
		return h.repo.ListOrdersByDriver(ctx, userID)
	default:
		return nil, nil
	}
}
