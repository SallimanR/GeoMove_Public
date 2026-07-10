package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterDriverRoutes(router *gin.RouterGroup, h *DriverHandler, authMiddleware gin.HandlerFunc) {
	driver := router.Group("/driver")
	{
		driver.POST("/", h.CreateDriver)
		driver.POST("/filter", h.GetFilteredDrivers)
		driver.GET("/:user_id", h.GetDriverByUserID)

		profile := driver.Group("/profile")
		profile.Use(authMiddleware)
		{
			profile.POST("/", h.CreateDriverProfile)
			profile.GET("/", h.GetMyDriverProfile)
			profile.POST("/image", h.UploadProfileImage)
		}
	}
}
