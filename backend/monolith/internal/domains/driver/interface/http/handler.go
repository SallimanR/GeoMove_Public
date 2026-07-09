package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"monolith/internal/auth"
	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/application/query"
)

type DriverHandler struct {
	createDriver       *command.CreateDriverHandler
	getDriverByUserID  *query.GetDriverByUserIDHandler
	getFilteredDrivers *query.GetFilteredDriversHandler
}

func NewDriverHandler(
	createDriver *command.CreateDriverHandler,
	getDriverByUserID *query.GetDriverByUserIDHandler,
	getFilteredDrivers *query.GetFilteredDriversHandler,
) *DriverHandler {
	return &DriverHandler{
		createDriver:       createDriver,
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
	err = h.createDriver.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, nil)
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
		t, err := time.Parse("15:04", *req.WorkStarts)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_starts format, expected HH:MM:SS"})
			return
		}
		workStarts = &t
	}
	if req.WorkEnds != nil {
		t, err := time.Parse("15:04", *req.WorkEnds)
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
	err := h.createDriver.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (h *DriverHandler) GetDriverByUserID(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	driver, err := h.getDriverByUserID.Handle(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "driver profile not found"})
		return
	}

	resp := Driver{
		UserId:       driver.UserID,
		Name:         driver.Name,
		Lat:          driver.Location.Lat,
		Lon:          driver.Location.Lon,
		ProfileImage: driver.ProfileImage,
		IsAvailable:  &driver.IsAvailable,
		WorkStarts:   driver.WorkStarts,
		WorkEnds:     driver.WorkEnds,
		Rating:       driver.Rating,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) GetMyDriverProfile(ctx *gin.Context) {
	sessionVal, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*auth.Session)

	driver, err := h.getDriverByUserID.Handle(ctx.Request.Context(), session.UserID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "driver profile not found"})
		return
	}

	resp := Driver{
		UserId:       driver.UserID,
		Name:         driver.Name,
		Lat:          driver.Location.Lat,
		Lon:          driver.Location.Lon,
		ProfileImage: driver.ProfileImage,
		IsAvailable:  &driver.IsAvailable,
		WorkStarts:   driver.WorkStarts,
		WorkEnds:     driver.WorkEnds,
		Rating:       driver.Rating,
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
			UserId:       driver.UserID,
			Name:         driver.Name,
			Lat:          driver.Location.Lat,
			Lon:          driver.Location.Lon,
			ProfileImage: driver.ProfileImage,
			IsAvailable:  &driver.IsAvailable,
			WorkStarts:   driver.WorkStarts,
			WorkEnds:     driver.WorkEnds,
			Rating:       driver.Rating,
		})
	}

	resp := GetFilteredDrivers200JSONResponse{
		Drivers: &driversRes,
	}
	ctx.JSON(http.StatusCreated, resp)
}
