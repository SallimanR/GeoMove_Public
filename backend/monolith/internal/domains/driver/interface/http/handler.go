package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/application/query"
	"monolith/internal/domains/driver/domain/entity"
	geoQuery "monolith/internal/domains/geolocation/application/query"
)

type DriverHandler struct {
	createDriver               *command.CreateDriverHandler
	getDriverByID              *query.GetDriverByIDHandler
	getFilteredDrivers         *query.GetFilteredDriversHandler
	findClosestDriversRealtime *geoQuery.FindClosestDriversRealtimeQuery
}

func NewDriverHandler(
	createDriver *command.CreateDriverHandler,
	findDriverByID *query.GetDriverByIDHandler,
	getFilteredDrivers *query.GetFilteredDriversHandler,
) *DriverHandler {
	return &DriverHandler{
		createDriver:       createDriver,
		getDriverByID:      findDriverByID,
		getFilteredDrivers: getFilteredDrivers,
	}
}

func (h *DriverHandler) CreateDriver(ctx *gin.Context) {
	var req CreateDriverJSONBody
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// start, _ := time.Parse(time.RFC3339, *req.WorkStarts)
	// end, _ := time.Parse(time.RFC3339, *req.WorkEnds)

	cmd := command.CreateDriverCommand{
		WorkStarts: &req.WorkStarts,
		WorkEnds:   &req.WorkEnds,
		Latitude:   req.Lat,
		Longitude:  req.Lon,
	}
	driver, err := h.createDriver.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := CreateDriver201JSONResponse{Id: int(driver.ID)}
	ctx.JSON(http.StatusCreated, resp)
}

func (h *DriverHandler) GetDriverByID(ctx *gin.Context) {
	var req GetDriverByIdRequestObject
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	cmd := query.GetDriverByIDQuery{
		DriverID: entity.DriverID(req.Id),
	}
	driver, err := h.getDriverByID.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	resp := GetDriverById200JSONResponse{
		WorkStarts: driver.WorkStarts,
		WorkEnds:   driver.WorkEnds,
		Rating:     driver.Rating,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) GetFilteredDrivers(ctx *gin.Context) {
	var req GetFilteredDriversJSONBody
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// start, _ := time.Parse(time.RFC3339, *req.WorkStarts)
	// end, _ := time.Parse(time.RFC3339, *req.WorkEnds)

	qry := query.GetFilteredDriversQuery{
		UserLat:    req.UserLat,
		UserLon:    req.UserLon,
		WorkStarts: req.WorkStarts,
		WorkEnds:   req.WorkEnds,
		MinRating:  req.MinRating,
	}

	drivers, err := h.getFilteredDrivers.Handle(ctx, qry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	driversRes := make([]Driver, 0, len(drivers))
	for _, driver := range drivers {
		driversRes = append(driversRes, Driver{
			Id:         int32(driver.ID),
			Name:       driver.Name,
			Lat:        driver.Location.Lat,
			Lon:        driver.Location.Lon,
			WorkStarts: driver.WorkStarts,
			WorkEnds:   driver.WorkEnds,
			Rating:     driver.Rating,
		})
	}

	resp := GetFilteredDrivers200JSONResponse{
		Drivers: &driversRes,
	}
	ctx.JSON(http.StatusCreated, resp)
}
