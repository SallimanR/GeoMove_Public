package http

import "github.com/gin-gonic/gin"

func RegisterOrderRoutes(router *gin.RouterGroup, h *OrderHandler, authMiddleware gin.HandlerFunc, driverOnly gin.HandlerFunc) {
	si := NewStrictHandler(h, nil)
	wrapper := ServerInterfaceWrapper{Handler: si}

	order := router.Group("/order")
	order.Use(authMiddleware)
	{
		order.POST("", wrapper.CreateOrder)
		order.POST("/estimate", h.EstimateOrderPrice)
		order.GET("/my", wrapper.ListMyOrders)
		order.DELETE("/my/active", wrapper.DeleteMyActiveOrder)
		order.GET("/:order_id", wrapper.GetOrder)
		order.PUT("/:order_id", wrapper.UpdateOrder)
	}

	order.GET("/available", driverOnly, wrapper.ListAvailableOrders)
	order.GET("/declined", driverOnly, wrapper.ListDeclinedOrders)
	order.PATCH("/:order_id/status", driverOnly, wrapper.UpdateOrderStatus)
	order.PATCH("/:order_id/decline", driverOnly, wrapper.DeclineOrder)
	order.DELETE("/:order_id/decline", driverOnly, wrapper.UndoDeclineOrder)
}
