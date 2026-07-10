package http

import (
	"bytes"
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
)

type DriverHandler struct {
	createDriver       *command.CreateDriverHandler
	getDriverByUserID  *query.GetDriverByUserIDHandler
	getFilteredDrivers *query.GetFilteredDriversHandler
	updateProfileImage *command.UpdateProfileImageHandler
	staticDir          string
}

func NewDriverHandler(
	createDriver *command.CreateDriverHandler,
	getDriverByUserID *query.GetDriverByUserIDHandler,
	getFilteredDrivers *query.GetFilteredDriversHandler,
	updateProfileImage *command.UpdateProfileImageHandler,
	staticDir string,
) *DriverHandler {
	return &DriverHandler{
		createDriver:       createDriver,
		getDriverByUserID:  getDriverByUserID,
		getFilteredDrivers: getFilteredDrivers,
		updateProfileImage: updateProfileImage,
		staticDir:          staticDir,
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
