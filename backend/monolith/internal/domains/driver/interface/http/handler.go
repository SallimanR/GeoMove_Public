package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"monolith/internal/auth"
	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/application/query"
	"monolith/internal/domains/driver/domain/entity"
)

type DriverHandler struct {
	createDriver       *command.CreateDriverHandler
	getDriverByID      *query.GetDriverByIDHandler
	getDriverByUserID  *query.GetDriverByUserIDHandler
	getFilteredDrivers *query.GetFilteredDriversHandler
}

func NewDriverHandler(
	createDriver *command.CreateDriverHandler,
	findDriverByID *query.GetDriverByIDHandler,
	getDriverByUserID *query.GetDriverByUserIDHandler,
	getFilteredDrivers *query.GetFilteredDriversHandler,
) *DriverHandler {
	return &DriverHandler{
		createDriver:       createDriver,
		getDriverByID:      findDriverByID,
		getDriverByUserID:  getDriverByUserID,
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

func (h *DriverHandler) CreateDriverProfile(ctx *gin.Context) {
	sessionVal, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*auth.Session)

	var req struct {
		Name       string  `json:"name"`
		Lat        float32 `json:"lat"`
		Lon        float32 `json:"lon"`
		WorkStarts *string `json:"work_starts"`
		WorkEnds   *string `json:"work_ends"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var workStarts, workEnds *time.Time
	if req.WorkStarts != nil {
		t, err := time.Parse("15:04:05", *req.WorkStarts)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_starts format, expected HH:MM:SS"})
			return
		}
		workStarts = &t
	}
	if req.WorkEnds != nil {
		t, err := time.Parse("15:04:05", *req.WorkEnds)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_ends format, expected HH:MM:SS"})
			return
		}
		workEnds = &t
	}

	cmd := command.CreateDriverCommand{
		UserID:     session.UserID,
		Name:       req.Name,
		WorkStarts: workStarts,
		WorkEnds:   workEnds,
		Latitude:   req.Lat,
		Longitude:  req.Lon,
	}
	driver, err := h.createDriver.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": int(driver.ID)})
}

func (h *DriverHandler) GetMyDriverProfile(ctx *gin.Context) {
	sessionVal, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*auth.Session)

	qry := query.GetDriverByUserIDQuery{
		UserID: session.UserID,
	}
	driver, err := h.getDriverByUserID.Handle(ctx.Request.Context(), qry)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "driver profile not found"})
		return
	}

	resp := Driver{
		Id:         int32(driver.ID),
		Name:       driver.Name,
		Lat:        driver.Location.Lat,
		Lon:        driver.Location.Lon,
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
