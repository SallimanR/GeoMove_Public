package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterDriverRoutes(router *gin.RouterGroup, h *DriverHandler, authMiddleware gin.HandlerFunc) {
	driver := router.Group("/driver")
	{
		driver.POST("/", h.CreateDriver)
		driver.GET("/:id", h.GetDriverByID)
		driver.POST("/filter", h.GetFilteredDrivers)

		profile := driver.Group("/profile")
		profile.Use(authMiddleware)
		{
			profile.POST("/", h.CreateDriverProfile)
			profile.GET("/", h.GetMyDriverProfile)
		}
	}
}
