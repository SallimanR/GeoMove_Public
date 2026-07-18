package order

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"monolith/internal/domains/order/application/command"
	"monolith/internal/domains/order/application/query"
	"monolith/internal/domains/order/infrastructure/postgres"
	"monolith/internal/domains/order/infrastructure/postgres/sqlc"
	orderHTTP "monolith/internal/domains/order/interface/http"
	"monolith/internal/notification"
)

type OrderDomain struct {
	Commands OrderCommands
	Queries  OrderQueries
}

type OrderCommands struct {
	CreateOrder       *command.CreateOrderHandler
	UpdateOrder       *command.UpdateOrderHandler
	UpdateOrderStatus *command.UpdateOrderStatusHandler
}

type OrderQueries struct {
	GetOrderByID     *query.GetOrderByIDHandler
	ListOrdersByUser *query.ListOrdersByUserHandler
}

func NewOrderDomain(db *pgxpool.Pool, notifSvc *notification.Service) *OrderDomain {
	orderRepo := postgres.NewOrderRepository(sqlc.New(db))

	createHandler := command.NewCreateOrderHandler(orderRepo, notifSvc)
	updateHandler := command.NewUpdateOrderHandler(orderRepo)
	updateStatusHandler := command.NewUpdateOrderStatusHandler(orderRepo, notifSvc)
	getOrderHandler := query.NewGetOrderByIDHandler(orderRepo)
	listOrdersHandler := query.NewListOrdersByUserHandler(orderRepo)

	return &OrderDomain{
		Commands: OrderCommands{
			CreateOrder:       createHandler,
			UpdateOrder:       updateHandler,
			UpdateOrderStatus: updateStatusHandler,
		},
		Queries: OrderQueries{
			GetOrderByID:     getOrderHandler,
			ListOrdersByUser: listOrdersHandler,
		},
	}
}

func (d *OrderDomain) RegisterHTTPRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	handler := orderHTTP.NewOrderHandler(
		d.Commands.CreateOrder,
		d.Commands.UpdateOrder,
		d.Commands.UpdateOrderStatus,
		d.Queries.GetOrderByID,
		d.Queries.ListOrdersByUser,
	)
	orderHTTP.RegisterOrderRoutes(router, handler, authMiddleware)
}
