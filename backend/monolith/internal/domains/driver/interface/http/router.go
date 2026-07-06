package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterDriverRoutes(router *gin.RouterGroup, h *DriverHandler) {
	driver := router.Group("/driver")
	{
		driver.POST("/", h.CreateDriver)
		driver.GET("/:id", h.GetDriverByID)
		driver.POST("/filter", h.GetFilteredDrivers)
		// driver.PUT("/")
		// driver.DELETE("/")
	}
}
