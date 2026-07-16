package http

import "github.com/gin-gonic/gin"

func RegisterGeolocationRoutes(router *gin.RouterGroup, h GeolocationHandler) {
	si := NewStrictHandler(&h, nil)
	wrapper := ServerInterfaceWrapper{Handler: si}

	geolocation := router.Group("/driver")
	{
		geolocation.GET("/moving/closest", wrapper.GetClosestWithinRadiusMovingDriversByIDs)
		geolocation.POST("/moving", wrapper.GetMovingDriversByIDs)
	}
}
