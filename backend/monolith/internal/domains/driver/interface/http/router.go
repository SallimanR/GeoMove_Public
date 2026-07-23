package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterDriverRoutes(router *gin.RouterGroup, h *DriverHandler, authMiddleware gin.HandlerFunc) {
	driver := router.Group("/driver")
	{
		driver.GET("/:user_id", h.GetDriverByUserID)
		driver.GET("/:user_id/freely-available", h.GetFreelyAvailableByID)
		driver.POST("/filter", h.GetFilteredDrivers)
		driver.POST("/freely-available/search", h.GetFreelyAvailableDrivers)

		profile := driver.Group("/profile")
		profile.Use(authMiddleware)
		{
			profile.POST("/", h.CreateDriverProfile)
			profile.GET("/", h.GetMyDriverProfile)
			profile.PUT("/", h.UpdateDriverProfile)
			profile.POST("/image", h.UploadProfileImage)
			profile.POST("/car-photo", h.UploadCarPhoto)
		}

		freely := driver.Group("/freely-available")
		freely.Use(authMiddleware)
		{
			freely.POST("/", h.CreateFreelyAvailable)
			freely.PUT("/", h.UpdateFreelyAvailable)
			freely.DELETE("/", h.DeleteFreelyAvailable)
		}
	}
}
