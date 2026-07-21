package repository

import (
	"context"
	"time"

	"monolith/internal/domains/order/domain/entity"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) (int64, error)
	GetOrderByID(ctx context.Context, id int64) (*entity.Order, error)
	ListOrdersByCustomer(ctx context.Context, customerID int64) ([]entity.Order, error)
	ListOrdersByDriver(ctx context.Context, driverID int64) ([]entity.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status entity.OrderStatus, acceptedAt, pickedUpAt, completedAt, cancelledAt *time.Time, cancellationReason *string) error
	UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	DeleteActiveOrder(ctx context.Context, customerID int64) error
}
