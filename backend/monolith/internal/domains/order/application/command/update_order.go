package command

import (
	"context"
	"fmt"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
)

var editableStatuses = map[entity.OrderStatus]bool{
	entity.OrderStatusForming: true,
	entity.OrderStatusPending: true,
}

type UpdateOrderCommand struct {
	OrderID              int64
	FromLat              float32
	FromLon              float32
	FromAddress          string
	ToLat                float32
	ToLon                float32
	ToAddress            string
	HowManyWheelsBlocked int16
	TotalDistanceMeters  *int32
	PriceRubles          *int32
}

type UpdateOrderHandler struct {
	repo repository.OrderRepository
}

func NewUpdateOrderHandler(repo repository.OrderRepository) *UpdateOrderHandler {
	return &UpdateOrderHandler{repo: repo}
}

func (h *UpdateOrderHandler) Handle(ctx context.Context, cmd UpdateOrderCommand) (*entity.Order, error) {
	order, err := h.repo.GetOrderByID(ctx, cmd.OrderID)
	if err != nil {
		return nil, err
	}

	if !editableStatuses[order.Status] {
		return nil, fmt.Errorf("заказ не может быть изменён на данном этапе: %s", order.Status)
	}

	order.FromLat = cmd.FromLat
	order.FromLon = cmd.FromLon
	order.FromAddress = cmd.FromAddress
	order.ToLat = cmd.ToLat
	order.ToLon = cmd.ToLon
	order.ToAddress = cmd.ToAddress
	order.HowManyWheelsBlocked = cmd.HowManyWheelsBlocked
	order.TotalDistanceMeters = cmd.TotalDistanceMeters
	order.PriceRubles = cmd.PriceRubles

	return h.repo.UpdateOrder(ctx, order)
}
