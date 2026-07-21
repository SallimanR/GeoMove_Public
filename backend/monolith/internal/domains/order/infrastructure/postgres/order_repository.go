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

type orderRow struct {
	ID                   int64
	CreatedAt            time.Time
	UpdatedAt            time.Time
	CustomerID           int64
	DriverID             *int64
	FromLon              float32
	FromLat              float32
	FromAddress          string
	ToLon                float32
	ToLat                float32
	ToAddress            string
	TotalDistanceMeters  *int32
	HowManyWheelsBlocked int16
	PriceRubles          *int32
	CarWeightKg          int32
	CarLengthMeters      float32
	CarType              sqlc.CarType
	CarName              string
	CarPhotoUrl          *string
	CustomerMessage      *string
	Status               sqlc.OrderStatus
	AcceptedAt           *time.Time
	PickedUpAt           *time.Time
	CompletedAt          *time.Time
	CancelledAt          *time.Time
	CancellationReason   *string
}

func toOrderEntity(r orderRow) entity.Order {
	return entity.Order{
		ID:                   r.ID,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
		CustomerID:           r.CustomerID,
		DriverID:             r.DriverID,
		FromLat:              r.FromLat,
		FromLon:              r.FromLon,
		FromAddress:          r.FromAddress,
		ToLat:                r.ToLat,
		ToLon:                r.ToLon,
		ToAddress:            r.ToAddress,
		TotalDistanceMeters:  r.TotalDistanceMeters,
		HowManyWheelsBlocked: r.HowManyWheelsBlocked,
		PriceRubles:          r.PriceRubles,
		CarWeightKg:          r.CarWeightKg,
		CarLengthMeters:      r.CarLengthMeters,
		CarType:              string(r.CarType),
		CarName:              r.CarName,
		CarPhotoUrl:          r.CarPhotoUrl,
		CustomerMessage:      r.CustomerMessage,
		Status:               entity.OrderStatus(r.Status),
		AcceptedAt:           r.AcceptedAt,
		PickedUpAt:           r.PickedUpAt,
		CompletedAt:          r.CompletedAt,
		CancelledAt:          r.CancelledAt,
		CancellationReason:   r.CancellationReason,
	}
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
		CarWeightKg:          order.CarWeightKg,
		CarLengthMeters:      order.CarLengthMeters,
		CarType:              sqlc.CarType(order.CarType),
		CarName:              order.CarName,
		CarPhotoUrl:          order.CarPhotoUrl,
		CustomerMessage:      order.CustomerMessage,
	})
	return id, err
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, id int64) (*entity.Order, error) {
	row, err := r.queries.GetOrderByID(ctx, id)
	if err != nil {
		return nil, err
	}
	e := toOrderEntity(orderRow{
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
		CarWeightKg:          row.CarWeightKg,
		CarLengthMeters:      row.CarLengthMeters,
		CarType:              row.CarType,
		CarName:              row.CarName,
		CarPhotoUrl:          row.CarPhotoUrl,
		CustomerMessage:      row.CustomerMessage,
		Status:               row.Status,
		AcceptedAt:           row.AcceptedAt,
		PickedUpAt:           row.PickedUpAt,
		CompletedAt:          row.CompletedAt,
		CancelledAt:          row.CancelledAt,
		CancellationReason:   row.CancellationReason,
	})
	return &e, nil
}

func (r *OrderRepository) ListOrdersByCustomer(ctx context.Context, customerID int64) ([]entity.Order, error) {
	rows, err := r.queries.ListOrdersByCustomer(ctx, customerID)
	if err != nil {
		return nil, err
	}
	result := make([]entity.Order, 0, len(rows))
	for _, row := range rows {
		result = append(result, toOrderEntity(orderRow{
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
			CarWeightKg:          row.CarWeightKg,
			CarLengthMeters:      row.CarLengthMeters,
			CarType:              row.CarType,
			CarName:              row.CarName,
			CarPhotoUrl:          row.CarPhotoUrl,
			CustomerMessage:      row.CustomerMessage,
			Status:               row.Status,
			AcceptedAt:           row.AcceptedAt,
			PickedUpAt:           row.PickedUpAt,
			CompletedAt:          row.CompletedAt,
			CancelledAt:          row.CancelledAt,
			CancellationReason:   row.CancellationReason,
		}))
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
		result = append(result, toOrderEntity(orderRow{
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
			CarWeightKg:          row.CarWeightKg,
			CarLengthMeters:      row.CarLengthMeters,
			CarType:              row.CarType,
			CarName:              row.CarName,
			CarPhotoUrl:          row.CarPhotoUrl,
			CustomerMessage:      row.CustomerMessage,
			Status:               row.Status,
			AcceptedAt:           row.AcceptedAt,
			PickedUpAt:           row.PickedUpAt,
			CompletedAt:          row.CompletedAt,
			CancelledAt:          row.CancelledAt,
			CancellationReason:   row.CancellationReason,
		}))
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
		CarWeightKg:          order.CarWeightKg,
		CarLengthMeters:      order.CarLengthMeters,
		CarType:              sqlc.CarType(order.CarType),
		CarName:              order.CarName,
		CarPhotoUrl:          order.CarPhotoUrl,
		CustomerMessage:      order.CustomerMessage,
	})
	if err != nil {
		return nil, err
	}
	e := toOrderEntity(orderRow{
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
		CarWeightKg:          row.CarWeightKg,
		CarLengthMeters:      row.CarLengthMeters,
		CarType:              row.CarType,
		CarName:              row.CarName,
		CarPhotoUrl:          row.CarPhotoUrl,
		CustomerMessage:      row.CustomerMessage,
		Status:               row.Status,
		AcceptedAt:           row.AcceptedAt,
		PickedUpAt:           row.PickedUpAt,
		CompletedAt:          row.CompletedAt,
		CancelledAt:          row.CancelledAt,
		CancellationReason:   row.CancellationReason,
	})
	return &e, nil
}

func (r *OrderRepository) DeleteActiveOrder(ctx context.Context, customerID int64) error {
	return r.queries.DeleteActiveOrder(ctx, customerID)
}

func (r *OrderRepository) ListAvailableOrders(ctx context.Context) ([]entity.Order, error) {
	rows, err := r.queries.ListAvailableOrders(ctx)
	if err != nil {
		return nil, err
	}
	orders := make([]entity.Order, len(rows))
	for i, row := range rows {
		orders[i] = toOrderEntity(orderRow{
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
			CarWeightKg:          row.CarWeightKg,
			CarLengthMeters:      row.CarLengthMeters,
			CarType:              row.CarType,
			CarName:              row.CarName,
			CarPhotoUrl:          row.CarPhotoUrl,
			CustomerMessage:      row.CustomerMessage,
			Status:               row.Status,
			AcceptedAt:           row.AcceptedAt,
			PickedUpAt:           row.PickedUpAt,
			CompletedAt:          row.CompletedAt,
			CancelledAt:          row.CancelledAt,
			CancellationReason:   row.CancellationReason,
		})
	}
	return orders, nil
}

func (r *OrderRepository) SetOrderDriver(ctx context.Context, orderID, driverID int64) error {
	_, err := r.queries.SetOrderDriver(ctx, sqlc.SetOrderDriverParams{ID: orderID, DriverID: &driverID})
	return err
}
