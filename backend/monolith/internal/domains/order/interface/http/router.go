package http

import "github.com/gin-gonic/gin"

func RegisterOrderRoutes(router *gin.RouterGroup, h *OrderHandler, authMiddleware gin.HandlerFunc) {
	si := NewStrictHandler(h, nil)
	wrapper := ServerInterfaceWrapper{Handler: si}

	order := router.Group("/order")
	order.Use(authMiddleware)
	{
		order.POST("", wrapper.CreateOrder)
		order.GET("/my", wrapper.ListMyOrders)
		order.DELETE("/my/active", wrapper.DeleteMyActiveOrder)
		order.GET("/available", wrapper.ListAvailableOrders)
		order.GET("/:order_id", wrapper.GetOrder)
		order.PUT("/:order_id", wrapper.UpdateOrder)
		order.PATCH("/:order_id/status", wrapper.UpdateOrderStatus)
	}
}
