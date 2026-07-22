package http

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"monolith/internal/auth"
	"monolith/internal/domains/driver/application/command"
	"monolith/internal/domains/driver/application/query"
	"monolith/internal/domains/driver/domain/entity"
	"monolith/pkg/geo"
)

type UserRoleManager interface {
	AddRole(ctx context.Context, userID int64, role string) error
}

type DriverHandler struct {
	createDriver              *command.CreateDriverHandler
	updateDriver              *command.UpdateDriverHandler
	getDriverByUserID         *query.GetDriverByUserIDHandler
	getFilteredDrivers        *query.GetFilteredDriversHandler
	updateProfileImage        *command.UpdateProfileImageHandler
	createFreelyAvailable     *command.CreateFreelyAvailableHandler
	updateFreelyAvailable     *command.UpdateFreelyAvailableHandler
	deleteFreelyAvailable     *command.DeleteFreelyAvailableHandler
	getFreelyAvailableByID    *query.GetFreelyAvailableByUserIDHandler
	getFreelyAvailableDrivers *query.GetFreelyAvailableDriversHandler
	staticDir                 string
	roleManager               UserRoleManager
}

func NewDriverHandler(
	createDriver *command.CreateDriverHandler,
	updateDriver *command.UpdateDriverHandler,
	getDriverByUserID *query.GetDriverByUserIDHandler,
	getFilteredDrivers *query.GetFilteredDriversHandler,
	updateProfileImage *command.UpdateProfileImageHandler,
	createFreelyAvailable *command.CreateFreelyAvailableHandler,
	updateFreelyAvailable *command.UpdateFreelyAvailableHandler,
	deleteFreelyAvailable *command.DeleteFreelyAvailableHandler,
	getFreelyAvailableByID *query.GetFreelyAvailableByUserIDHandler,
	getFreelyAvailableDrivers *query.GetFreelyAvailableDriversHandler,
	staticDir string,
	roleManager UserRoleManager,
) *DriverHandler {
	return &DriverHandler{
		createDriver:              createDriver,
		updateDriver:              updateDriver,
		getDriverByUserID:         getDriverByUserID,
		getFilteredDrivers:        getFilteredDrivers,
		updateProfileImage:        updateProfileImage,
		createFreelyAvailable:     createFreelyAvailable,
		updateFreelyAvailable:     updateFreelyAvailable,
		deleteFreelyAvailable:     deleteFreelyAvailable,
		getFreelyAvailableByID:    getFreelyAvailableByID,
		getFreelyAvailableDrivers: getFreelyAvailableDrivers,
		staticDir:                 staticDir,
		roleManager:               roleManager,
	}
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
		Phone               *string  `json:"phone"`
		WorkStarts          *string  `json:"work_starts"`
		WorkEnds            *string  `json:"work_ends"`
		MaxCarWeightKg      *int32   `json:"max_car_weight_kg"`
		MaxCarLengthMeters  *float32 `json:"max_car_length_meters"`
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

	address, err := geo.ReverseGeocode(ctx.Request.Context(), req.Lat, req.Lon)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to resolve address: " + err.Error()})
		return
	}

	cmd := command.CreateDriverCommand{
		UserID:             session.UserID,
		Name:               req.Name,
		Phone:              req.Phone,
		WorkStarts:         workStarts,
		WorkEnds:           workEnds,
		Latitude:           req.Lat,
		Longitude:          req.Lon,
		MaxCarWeightKg:     req.MaxCarWeightKg,
		MaxCarLengthMeters: req.MaxCarLengthMeters,
		Address:            address,
	}
	err = h.createDriver.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if h.roleManager != nil {
		_ = h.roleManager.AddRole(ctx.Request.Context(), session.UserID, "tow_driver")
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func driverToResponse(driver *entity.Driver) Driver {
	var maxCarWeightKg *int
	if driver.MaxCarWeightKg > 0 {
		v := int(driver.MaxCarWeightKg)
		maxCarWeightKg = &v
	}
	var maxCarLengthMeters *float32
	if driver.MaxCarLengthMeters > 0 {
		maxCarLengthMeters = &driver.MaxCarLengthMeters
	}
	phone := &driver.Phone
	if driver.Phone == "" {
		phone = nil
	}
	rating := &driver.Rating
	if driver.Rating == 0 {
		rating = nil
	}
	address := &driver.Address
	if driver.Address == "" {
		address = nil
	}
	resp := Driver{
		UserId:             driver.UserID,
		Name:               driver.Name,
		Lat:                driver.Location.Lat,
		Lon:                driver.Location.Lon,
		ProfileImage:       driver.ProfileImage,
		IsAvailable:        &driver.IsAvailable,
		WorkStarts:         driver.WorkStarts,
		WorkEnds:           driver.WorkEnds,
		Rating:             rating,
		Phone:              phone,
		MaxCarWeightKg:     maxCarWeightKg,
		MaxCarLengthMeters: maxCarLengthMeters,
		Address:            address,
	}
	return resp
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

	resp := driverToResponse(driver)
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

	resp := driverToResponse(driver)
	ctx.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) UpdateDriverProfile(ctx *gin.Context) {
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
		Phone               *string  `json:"phone"`
		WorkStarts          *string  `json:"work_starts"`
		WorkEnds            *string  `json:"work_ends"`
		MaxCarWeightKg      *int32   `json:"max_car_weight_kg"`
		MaxCarLengthMeters  *float32 `json:"max_car_length_meters"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var workStarts, workEnds *time.Time
	if req.WorkStarts != nil {
		t, err := time.Parse("15:04", *req.WorkStarts)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_starts format, expected HH:MM"})
			return
		}
		workStarts = &t
	}
	if req.WorkEnds != nil {
		t, err := time.Parse("15:04", *req.WorkEnds)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid work_ends format, expected HH:MM"})
			return
		}
		workEnds = &t
	}

	address, err := geo.ReverseGeocode(ctx.Request.Context(), req.Lat, req.Lon)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to resolve address: " + err.Error()})
		return
	}

	cmd := command.UpdateDriverCommand{
		UserID:             session.UserID,
		Name:               req.Name,
		Phone:              req.Phone,
		WorkStarts:         workStarts,
		WorkEnds:           workEnds,
		Latitude:           req.Lat,
		Longitude:          req.Lon,
		MaxCarWeightKg:     req.MaxCarWeightKg,
		MaxCarLengthMeters: req.MaxCarLengthMeters,
		Address:            address,
	}
	err = h.updateDriver.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *DriverHandler) UploadProfileImage(ctx *gin.Context) {
	sessionVal, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*auth.Session)

	var req struct {
		Image string `json:"image"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	imageURL, err := saveProfileImage(req.Image, h.staticDir)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.UpdateProfileImageCommand{
		UserID:   session.UserID,
		ImageURL: imageURL,
	}
	if err := h.updateProfileImage.Handle(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile image"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"image_url": imageURL})
}

func saveProfileImage(imageBase64, staticDir string) (string, error) {
	parts := bytes.SplitN([]byte(imageBase64), []byte(","), 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid image data")
	}
	data, err := base64.StdEncoding.DecodeString(string(parts[1]))
	if err != nil {
		return "", fmt.Errorf("invalid base64")
	}
	const maxSize = 5 * 1024 * 1024
	if len(data) > maxSize {
		return "", fmt.Errorf("image too large")
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("invalid image format")
	}
	if format != "jpeg" && format != "png" {
		return "", fmt.Errorf("unsupported format: %s", format)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
		return "", fmt.Errorf("encoding image: %w", err)
	}

	uploadDir := staticDir
	if uploadDir == "" {
		uploadDir = "../../data/static/avatars"
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	filename := uuid.New().String() + ".jpg"
	filePath := filepath.Join(uploadDir, filename)
	if err := os.WriteFile(filePath, buf.Bytes(), 0o644); err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}
	return "/static/avatars/" + filename, nil
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
		driversRes = append(driversRes, driverToResponse(&driver))
	}

	resp := GetFilteredDrivers200JSONResponse{
		Drivers: &driversRes,
	}
	ctx.JSON(http.StatusCreated, resp)
}

func (h *DriverHandler) CreateFreelyAvailable(ctx *gin.Context) {
	sessionVal, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*auth.Session)

	var req FreelyAvailable
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toLocations := make([]entity.LocationWithAddress, 0)
	if req.ToLocations != nil {
		for _, l := range *req.ToLocations {
			toLocations = append(toLocations, entity.LocationWithAddress{Lat: l.Lat, Lon: l.Lon})
		}
	}

	cmd := command.CreateFreelyAvailableCommand{
		UserID:       session.UserID,
		FromDate:     req.FromDate,
		ToDate:       req.ToDate,
		FromLocation: entity.LocationWithAddress{Lat: req.FromLocation.Lat, Lon: req.FromLocation.Lon},
		ToLocations:  toLocations,
		EnRouteOrder: req.EnRouteOrder,
		TariffPerKm:  req.TariffPerKm,
	}
	if err := h.createFreelyAvailable.Handle(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (h *DriverHandler) UpdateFreelyAvailable(ctx *gin.Context) {
	sessionVal, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*auth.Session)

	var req FreelyAvailable
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toLocations := make([]entity.LocationWithAddress, 0)
	if req.ToLocations != nil {
		for _, l := range *req.ToLocations {
			toLocations = append(toLocations, entity.LocationWithAddress{Lat: l.Lat, Lon: l.Lon})
		}
	}

	cmd := command.UpdateFreelyAvailableCommand{
		UserID:       session.UserID,
		FromDate:     req.FromDate,
		ToDate:       req.ToDate,
		FromLocation: entity.LocationWithAddress{Lat: req.FromLocation.Lat, Lon: req.FromLocation.Lon},
		ToLocations:  toLocations,
		EnRouteOrder: req.EnRouteOrder,
		TariffPerKm:  req.TariffPerKm,
	}
	if err := h.updateFreelyAvailable.Handle(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *DriverHandler) DeleteFreelyAvailable(ctx *gin.Context) {
	sessionVal, exists := ctx.Get("session")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*auth.Session)

	cmd := command.DeleteFreelyAvailableCommand{UserID: session.UserID}
	if err := h.deleteFreelyAvailable.Handle(ctx.Request.Context(), cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *DriverHandler) GetFreelyAvailableByID(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	fa, err := h.getFreelyAvailableByID.Handle(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "freely available entry not found"})
		return
	}

	toLocations := make([]Location, 0, len(fa.ToLocations))
	for _, l := range fa.ToLocations {
		toLocations = append(toLocations, Location{Lat: l.Lat, Lon: l.Lon, Address: &l.Address})
	}

	fromAddress := fa.FromLocation.Address
	resp := FreelyAvailableResponse{
		UserId:       fa.UserID,
		FromDate:     fa.FromDate,
		ToDate:       fa.ToDate,
		FromLocation: Location{Lat: fa.FromLocation.Lat, Lon: fa.FromLocation.Lon, Address: &fromAddress},
		ToLocations:  &toLocations,
		EnRouteOrder: fa.EnRouteOrder,
		TariffPerKm:  fa.TariffPerKm,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *DriverHandler) GetFreelyAvailableDrivers(ctx *gin.Context) {
	var req GetFreelyAvailableDriversRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qry := query.GetFreelyAvailableDriversQuery{
		UserLat:      req.UserLat,
		UserLon:      req.UserLon,
		EnRouteOrder: req.EnRouteOrder,
		MinTariff:    req.MinTariff,
		MaxTariff:    req.MaxTariff,
	}
	drivers, err := h.getFreelyAvailableDrivers.Handle(ctx, qry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	items := make([]FreelyAvailableDriverResponse, 0, len(drivers))
	for _, d := range drivers {
		items = append(items, FreelyAvailableDriverResponse{
			UserId:       d.UserID,
			FromDate:     d.FromDate,
			ToDate:       d.ToDate,
			FromLocation: Location{Lat: d.FromLocation.Lat, Lon: d.FromLocation.Lon, Address: &d.FromLocation.Address},
			EnRouteOrder: d.EnRouteOrder,
			TariffPerKm:  d.TariffPerKm,
			Name:         d.Name,
			Rating:       d.Rating,
			ProfileImage: d.ProfileImage,
			Distance:     d.Distance,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"drivers": items})
}
