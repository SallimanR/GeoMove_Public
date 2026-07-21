package query

import (
	"context"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
)

type ListAvailableOrdersHandler struct {
	repo repository.OrderRepository
}

func NewListAvailableOrdersHandler(repo repository.OrderRepository) *ListAvailableOrdersHandler {
	return &ListAvailableOrdersHandler{repo: repo}
}

func (h *ListAvailableOrdersHandler) Handle(ctx context.Context) ([]entity.Order, error) {
	return h.repo.ListAvailableOrders(ctx)
}
