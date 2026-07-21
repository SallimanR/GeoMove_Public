package command

import (
	"context"
	"fmt"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
)

type DeleteActiveOrderHandler struct {
	repo repository.OrderRepository
}

func NewDeleteActiveOrderHandler(repo repository.OrderRepository) *DeleteActiveOrderHandler {
	return &DeleteActiveOrderHandler{repo: repo}
}

func (h *DeleteActiveOrderHandler) Handle(ctx context.Context, customerID int64) error {
	orders, err := h.repo.ListOrdersByCustomer(ctx, customerID)
	if err != nil {
		return fmt.Errorf("ошибка получения заказов: %w", err)
	}

	var activeOrder *entity.Order
	for i := range orders {
		o := &orders[i]
		if o.Status == entity.OrderStatusForming || o.Status == entity.OrderStatusPending {
			activeOrder = o
			break
		}
	}

	if activeOrder == nil {
		return fmt.Errorf("нет активного заказа для удаления")
	}

	return h.repo.DeleteActiveOrder(ctx, customerID)
}
