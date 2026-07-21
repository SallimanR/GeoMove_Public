package http

import (
	"context"

	"github.com/gin-gonic/gin"

	"monolith/internal/auth"
	"monolith/internal/domains/order/application/command"
	"monolith/internal/domains/order/application/query"
	"monolith/internal/domains/order/domain/entity"
)

func err400(msg string) CreateOrder400JSONResponse {
	return CreateOrder400JSONResponse{Error: &msg}
}

func err400Status(msg string) UpdateOrderStatus400JSONResponse {
	return UpdateOrderStatus400JSONResponse{Error: &msg}
}

func err400Update(msg string) UpdateOrder400JSONResponse {
	return UpdateOrder400JSONResponse{Error: &msg}
}

type OrderHandler struct {
	createOrder         *command.CreateOrderHandler
	updateOrder         *command.UpdateOrderHandler
	updateOrderStatus   *command.UpdateOrderStatusHandler
	deleteActiveOrder   *command.DeleteActiveOrderHandler
	getOrderByID        *query.GetOrderByIDHandler
	listOrdersByUser    *query.ListOrdersByUserHandler
	listAvailableOrders *query.ListAvailableOrdersHandler
}

func NewOrderHandler(
	createOrder *command.CreateOrderHandler,
	updateOrder *command.UpdateOrderHandler,
	updateOrderStatus *command.UpdateOrderStatusHandler,
	deleteActiveOrder *command.DeleteActiveOrderHandler,
	getOrderByID *query.GetOrderByIDHandler,
	listOrdersByUser *query.ListOrdersByUserHandler,
	listAvailableOrders *query.ListAvailableOrdersHandler,
) *OrderHandler {
	return &OrderHandler{
		createOrder:         createOrder,
		updateOrder:         updateOrder,
		updateOrderStatus:   updateOrderStatus,
		deleteActiveOrder:   deleteActiveOrder,
		getOrderByID:        getOrderByID,
		listOrdersByUser:    listOrdersByUser,
		listAvailableOrders: listAvailableOrders,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, request CreateOrderRequestObject) (CreateOrderResponseObject, error) {
	session := getSession(ctx)
	if session == nil {
		return CreateOrder401Response{}, nil
	}

	body := request.Body
	if body == nil {
		return err400("неверный запрос"), nil
	}

	var totalDistanceMeters *int32
	if body.TotalDistanceMeters != nil {
		v := int32(*body.TotalDistanceMeters)
		totalDistanceMeters = &v
	}
	var priceRubles *int32
	if body.PriceRubles != nil {
		v := int32(*body.PriceRubles)
		priceRubles = &v
	}

	cmd := command.CreateOrderCommand{
		CustomerID:           session.UserID,
		FromLat:              body.FromLat,
		FromLon:              body.FromLon,
		FromAddress:          stringPtrValue(body.FromAddress),
		ToLat:                body.ToLat,
		ToLon:                body.ToLon,
		ToAddress:            stringPtrValue(body.ToAddress),
		HowManyWheelsBlocked: int16(body.HowManyWheelsBlocked),
		TotalDistanceMeters:  totalDistanceMeters,
		PriceRubles:          priceRubles,
		CarWeightKg:          int32(body.CarWeightKg),
		CarLengthMeters:      body.CarLengthMeters,
		CarType:              string(body.CarType),
		CarName:              body.CarName,
		CarPhotoUrl:          body.CarPhotoUrl,
		CustomerMessage:      body.CustomerMessage,
	}

	order, err := h.createOrder.Handle(ctx, cmd)
	if err != nil {
		return err400(err.Error()), nil
	}

	return CreateOrder201JSONResponse(toAPIOrder(order)), nil
}

func (h *OrderHandler) ListMyOrders(ctx context.Context, request ListMyOrdersRequestObject) (ListMyOrdersResponseObject, error) {
	session := getSession(ctx)
	if session == nil {
		return ListMyOrders401Response{}, nil
	}

	orders, err := h.listOrdersByUser.Handle(ctx, session.UserID, string(request.Params.Role))
	if err != nil {
		return ListMyOrders401Response{}, nil
	}

	apiOrders := make([]Order, 0, len(orders))
	for i := range orders {
		apiOrders = append(apiOrders, toAPIOrder(&orders[i]))
	}

	return ListMyOrders200JSONResponse{Orders: &apiOrders}, nil
}

func (h *OrderHandler) ListAvailableOrders(ctx context.Context, request ListAvailableOrdersRequestObject) (ListAvailableOrdersResponseObject, error) {
	session := getSession(ctx)
	if session == nil {
		return ListAvailableOrders401Response{}, nil
	}

	orders, err := h.listAvailableOrders.Handle(ctx)
	if err != nil {
		return ListAvailableOrders200JSONResponse{}, nil
	}

	apiOrders := make([]Order, len(orders))
	for i := range orders {
		apiOrders[i] = toAPIOrder(&orders[i])
	}
	return ListAvailableOrders200JSONResponse{Orders: &apiOrders}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, request GetOrderRequestObject) (GetOrderResponseObject, error) {
	session := getSession(ctx)
	if session == nil {
		return GetOrder401Response{}, nil
	}

	order, err := h.getOrderByID.Handle(ctx, request.OrderId)
	if err != nil {
		return GetOrder404Response{}, nil
	}

	if order.CustomerID != session.UserID && (order.DriverID == nil || *order.DriverID != session.UserID) {
		return GetOrder404Response{}, nil
	}

	return GetOrder200JSONResponse(toAPIOrder(order)), nil
}

func (h *OrderHandler) UpdateOrder(ctx context.Context, request UpdateOrderRequestObject) (UpdateOrderResponseObject, error) {
	session := getSession(ctx)
	if session == nil {
		return UpdateOrder401Response{}, nil
	}

	body := request.Body
	if body == nil {
		return err400Update("неверный запрос"), nil
	}

	var totalDistanceMeters *int32
	if body.TotalDistanceMeters != nil {
		v := int32(*body.TotalDistanceMeters)
		totalDistanceMeters = &v
	}
	var priceRubles *int32
	if body.PriceRubles != nil {
		v := int32(*body.PriceRubles)
		priceRubles = &v
	}

	cmd := command.UpdateOrderCommand{
		OrderID:              request.OrderId,
		FromLat:              body.FromLat,
		FromLon:              body.FromLon,
		FromAddress:          stringPtrValue(body.FromAddress),
		ToLat:                body.ToLat,
		ToLon:                body.ToLon,
		ToAddress:            stringPtrValue(body.ToAddress),
		HowManyWheelsBlocked: int16(body.HowManyWheelsBlocked),
		TotalDistanceMeters:  totalDistanceMeters,
		PriceRubles:          priceRubles,
		CarWeightKg:          int32(body.CarWeightKg),
		CarLengthMeters:      body.CarLengthMeters,
		CarType:              string(body.CarType),
		CarName:              body.CarName,
		CarPhotoUrl:          body.CarPhotoUrl,
		CustomerMessage:      body.CustomerMessage,
	}

	order, err := h.updateOrder.Handle(ctx, cmd)
	if err != nil {
		return err400Update(err.Error()), nil
	}

	return UpdateOrder200JSONResponse(toAPIOrder(order)), nil
}

func (h *OrderHandler) UpdateOrderStatus(ctx context.Context, request UpdateOrderStatusRequestObject) (UpdateOrderStatusResponseObject, error) {
	session := getSession(ctx)
	if session == nil {
		return UpdateOrderStatus401Response{}, nil
	}

	body := request.Body
	if body == nil {
		return err400Status("неверный запрос"), nil
	}

	cmd := command.UpdateOrderStatusCommand{
		OrderID:            request.OrderId,
		Status:             entity.OrderStatus(body.Status),
		DriverID:           &session.UserID,
		CancellationReason: body.CancellationReason,
	}

	order, err := h.updateOrderStatus.Handle(ctx, cmd)
	if err != nil {
		return err400Status(err.Error()), nil
	}

	return UpdateOrderStatus200JSONResponse(toAPIOrder(order)), nil
}

func err400Delete(msg string) DeleteMyActiveOrder400JSONResponse {
	return DeleteMyActiveOrder400JSONResponse{Error: &msg}
}

func (h *OrderHandler) DeleteMyActiveOrder(ctx context.Context, request DeleteMyActiveOrderRequestObject) (DeleteMyActiveOrderResponseObject, error) {
	session := getSession(ctx)
	if session == nil {
		return DeleteMyActiveOrder401Response{}, nil
	}

	if err := h.deleteActiveOrder.Handle(ctx, session.UserID); err != nil {
		return err400Delete(err.Error()), nil
	}

	return DeleteMyActiveOrder204Response{}, nil
}

func getSession(ctx context.Context) *auth.Session {
	ginCtx, ok := ctx.(*gin.Context)
	if !ok {
		return nil
	}
	sessionVal, exists := ginCtx.Get("session")
	if !exists {
		return nil
	}
	session, ok := sessionVal.(*auth.Session)
	if !ok {
		return nil
	}
	return session
}

func int32PtrToIntPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	r := int(*v)
	return &r
}

func stringPtrValue(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func strPtr(s string) *string {
	return &s
}

func toAPIOrder(o *entity.Order) Order {
	return Order{
		Id:                   o.ID,
		CreatedAt:            o.CreatedAt,
		UpdatedAt:            o.UpdatedAt,
		CustomerId:           o.CustomerID,
		DriverId:             o.DriverID,
		FromLat:              o.FromLat,
		FromLon:              o.FromLon,
		FromAddress:          strPtr(o.FromAddress),
		ToLat:                o.ToLat,
		ToLon:                o.ToLon,
		ToAddress:            strPtr(o.ToAddress),
		TotalDistanceMeters:  int32PtrToIntPtr(o.TotalDistanceMeters),
		HowManyWheelsBlocked: int(o.HowManyWheelsBlocked),
		PriceRubles:          int32PtrToIntPtr(o.PriceRubles),
		Status:               OrderStatus(o.Status),
		AcceptedAt:           o.AcceptedAt,
		PickedUpAt:           o.PickedUpAt,
		CompletedAt:          o.CompletedAt,
		CancelledAt:          o.CancelledAt,
		CancellationReason:   o.CancellationReason,
		CarWeightKg:          int(o.CarWeightKg),
		CarLengthMeters:      o.CarLengthMeters,
		CarType:              OrderCarType(o.CarType),
		CarName:              o.CarName,
		CarPhotoUrl:          o.CarPhotoUrl,
		CustomerMessage:      o.CustomerMessage,
	}
}
