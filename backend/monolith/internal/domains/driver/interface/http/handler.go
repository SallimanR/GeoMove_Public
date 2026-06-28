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
	findDriverByID             *query.GetDriverByIDHandler
	findClosestDriversRealtime *geoQuery.FindClosestDriversRealtimeQuery
}

func NewDriverHandler(
	createDriver *command.CreateDriverHandler,
	findDriverByID *query.GetDriverByIDHandler,
) *DriverHandler {
	return &DriverHandler{
		createDriver:   createDriver,
		findDriverByID: findDriverByID,
	}
}

func (h *DriverHandler) CreateDriver(ctx *gin.Context) {
	var req CreateDriverJSONBody
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// start, _ := time.Parse(time.RFC3339, *req.WorkStarts) // error handling omitted
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

func (h *DriverHandler) FindDrivers(ctx *gin.Context) {
	var req CreateDriverJSONBody
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// start, _ := time.Parse(time.RFC3339, *req.WorkStarts) // error handling omitted
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
	var req FindDriverByIdRequestObject
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	cmd := query.GetDriverByIDQuery{
		DriverID: entity.DriverID(req.Id),
	}
	driver, err := h.findDriverByID.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	resp := FindDriverById200JSONResponse{
		WorkStarts: driver.WorkStarts,
		WorkEnds:   driver.WorkEnds,
		Rating:     driver.Rating,
	}
	ctx.JSON(http.StatusOK, resp)
}
