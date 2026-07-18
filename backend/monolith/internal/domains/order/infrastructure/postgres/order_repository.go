package postgres

import (
	"context"
	"time"

	"monolith/internal/domains/order/domain/entity"
	"monolith/internal/domains/order/domain/repository"
	"monolith/internal/domains/order/infrastructure/postgres/sqlc"
)

type OrderRepository struct {
	queries sqlc.Queries
}

func NewOrderRepository(queries *sqlc.Queries) repository.OrderRepository {
	return &OrderRepository{queries: *queries}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *entity.Order) (int64, error) {
	id, err := r.queries.CreateOrder(ctx, sqlc.CreateOrderParams{
		CustomerID:           order.CustomerID,
		FromLat:              order.FromLat,
		FromLon:              order.FromLon,
		FromAddress:          order.FromAddress,
		ToLat:                order.ToLat,
		ToLon:                order.ToLon,
		ToAddress:            order.ToAddress,
		TotalDistanceMeters:  order.TotalDistanceMeters,
		HowManyWheelsBlocked: order.HowManyWheelsBlocked,
		PriceRubles:          order.PriceRubles,
	})
	return id, err
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id int64) (*entity.Order, error) {
	row, err := r.queries.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.Order{
		ID:                   row.ID,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		CustomerID:           row.CustomerID,
		DriverID:             row.DriverID,
		FromLat:              row.FromLat,
		FromLon:              row.FromLon,
		FromAddress:          row.FromAddress,
		ToLat:                row.ToLat,
		ToLon:                row.ToLon,
		ToAddress:            row.ToAddress,
		TotalDistanceMeters:  row.TotalDistanceMeters,
		HowManyWheelsBlocked: row.HowManyWheelsBlocked,
		PriceRubles:          row.PriceRubles,
		Status:               entity.OrderStatus(row.Status),
		AcceptedAt:           row.AcceptedAt,
		PickedUpAt:           row.PickedUpAt,
		CompletedAt:          row.CompletedAt,
		CancelledAt:          row.CancelledAt,
		CancellationReason:   row.CancellationReason,
	}, nil
}

func (r *OrderRepository) ListOrdersByCustomer(ctx context.Context, customerID int64) ([]entity.Order, error) {
	rows, err := r.queries.ListOrdersByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Order, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.Order{
			ID:                   row.ID,
			CreatedAt:            row.CreatedAt,
			UpdatedAt:            row.UpdatedAt,
			CustomerID:           row.CustomerID,
			DriverID:             row.DriverID,
			FromLat:              row.FromLat,
			FromLon:              row.FromLon,
			FromAddress:          row.FromAddress,
			ToLat:                row.ToLat,
			ToLon:                row.ToLon,
			ToAddress:            row.ToAddress,
			TotalDistanceMeters:  row.TotalDistanceMeters,
			HowManyWheelsBlocked: row.HowManyWheelsBlocked,
			PriceRubles:          row.PriceRubles,
			Status:               entity.OrderStatus(row.Status),
			AcceptedAt:           row.AcceptedAt,
			PickedUpAt:           row.PickedUpAt,
			CompletedAt:          row.CompletedAt,
			CancelledAt:          row.CancelledAt,
			CancellationReason:   row.CancellationReason,
		})
	}
	return result, nil
}

func (r *OrderRepository) ListOrdersByDriver(ctx context.Context, driverID int64) ([]entity.Order, error) {
	rows, err := r.queries.ListOrdersByDriver(ctx, &driverID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Order, 0, len(rows))
	for _, row := range rows {
		result = append(result, entity.Order{
			ID:                   row.ID,
			CreatedAt:            row.CreatedAt,
			UpdatedAt:            row.UpdatedAt,
			CustomerID:           row.CustomerID,
			DriverID:             row.DriverID,
			FromLat:              row.FromLat,
			FromLon:              row.FromLon,
			FromAddress:          row.FromAddress,
			ToLat:                row.ToLat,
			ToLon:                row.ToLon,
			ToAddress:            row.ToAddress,
			TotalDistanceMeters:  row.TotalDistanceMeters,
			HowManyWheelsBlocked: row.HowManyWheelsBlocked,
			PriceRubles:          row.PriceRubles,
			Status:               entity.OrderStatus(row.Status),
			AcceptedAt:           row.AcceptedAt,
			PickedUpAt:           row.PickedUpAt,
			CompletedAt:          row.CompletedAt,
			CancelledAt:          row.CancelledAt,
			CancellationReason:   row.CancellationReason,
		})
	}
	return result, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, id int64, status entity.OrderStatus, acceptedAt, pickedUpAt, completedAt, cancelledAt *time.Time, cancellationReason *string) error {
	params := sqlc.UpdateOrderStatusParams{
		ID:                 id,
		Status:             sqlc.OrderStatus(string(status)),
		AcceptedAt:         acceptedAt,
		PickedUpAt:         pickedUpAt,
		CompletedAt:        completedAt,
		CancelledAt:        cancelledAt,
		CancellationReason: cancellationReason,
	}
	return r.queries.UpdateOrderStatus(ctx, params)
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	row, err := r.queries.UpdateOrderDetails(ctx, sqlc.UpdateOrderDetailsParams{
		ID:                   order.ID,
		FromLat:              order.FromLat,
		FromLon:              order.FromLon,
		FromAddress:          order.FromAddress,
		ToLat:                order.ToLat,
		ToLon:                order.ToLon,
		ToAddress:            order.ToAddress,
		TotalDistanceMeters:  order.TotalDistanceMeters,
		HowManyWheelsBlocked: order.HowManyWheelsBlocked,
		PriceRubles:          order.PriceRubles,
	})
	if err != nil {
		return nil, err
	}
	return &entity.Order{
		ID:                   row.ID,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
		CustomerID:           row.CustomerID,
		DriverID:             row.DriverID,
		FromLat:              row.FromLat,
		FromLon:              row.FromLon,
		FromAddress:          row.FromAddress,
		ToLat:                row.ToLat,
		ToLon:                row.ToLon,
		ToAddress:            row.ToAddress,
		TotalDistanceMeters:  row.TotalDistanceMeters,
		HowManyWheelsBlocked: row.HowManyWheelsBlocked,
		PriceRubles:          row.PriceRubles,
		Status:               entity.OrderStatus(row.Status),
		AcceptedAt:           row.AcceptedAt,
		PickedUpAt:           row.PickedUpAt,
		CompletedAt:          row.CompletedAt,
		CancelledAt:          row.CancelledAt,
		CancellationReason:   row.CancellationReason,
	}, nil
}
