package http

import "github.com/gin-gonic/gin"

func RegisterGeolocationRoutes(router *gin.RouterGroup, h GeolocationHandler) {
	geolocation := router.Group("/drivers")
	{
		geolocation.GET("/closest/lat=...&lon=...", h.FindClosestDriversRealtime)
	}
}
