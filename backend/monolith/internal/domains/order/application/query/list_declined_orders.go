package query

import (
	"context"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
)

type ListDeclinedOrdersHandler struct {
	repo repository.OrderRepository
}

func NewListDeclinedOrdersHandler(repo repository.OrderRepository) *ListDeclinedOrdersHandler {
	return &ListDeclinedOrdersHandler{repo: repo}
}

func (h *ListDeclinedOrdersHandler) Handle(ctx context.Context, driverID int64) ([]entity.Order, error) {
	return h.repo.ListDeclinedOrders(ctx, driverID)
}
