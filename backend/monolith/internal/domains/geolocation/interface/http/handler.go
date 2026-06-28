package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"monolith/internal/domains/geolocation/application/dto"
	"monolith/internal/domains/geolocation/application/query"
)

type GeolocationHandler struct {
	findClosestDriversRealtime query.FindClosestDriversRealtimeHandler
}

func NewGeolocationHandler(
	findClosestDriversRealtime query.FindClosestDriversRealtimeHandler,
) GeolocationHandler {
	return GeolocationHandler{
		findClosestDriversRealtime: findClosestDriversRealtime,
	}
}

func (h *GeolocationHandler) FindClosestDriversRealtime(ctx *gin.Context) {
	var req FindNearbyDriversParams
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := query.FindClosestDriversRealtimeQuery{
		Location: dto.Location{
			Latitude:  req.Lat,
			Longitude: req.Lon,
		},
	}

	drivers, err := h.findClosestDriversRealtime.Handle(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make(FindNearbyDrivers200JSONResponse, len(drivers))
	for i := 0; i < len(drivers); i++ {
		resp[i].DriverId = int(drivers[i].DriverID)
		resp[i].DistanceMeters = int(drivers[i].DistanceMeters)
	}
	ctx.JSON(http.StatusOK, resp)
}
