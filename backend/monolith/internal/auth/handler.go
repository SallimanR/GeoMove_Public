package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"monolith/internal/auth/openapi"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service      *Service
	cookieSecure bool
	cookieDomain string
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service:      service,
		cookieSecure: os.Getenv("COOKIE_SECURE") == "true",
		cookieDomain: os.Getenv("COOKIE_DOMAIN"),
	}
}

func (h *Handler) PostAuthProviderCallback(c *gin.Context) {
	provider := c.Param("provider")

	var req openapi.PostAuthProviderCallbackJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	userID, email, err := h.service.ExchangeOAuthCode(c.Request.Context(), provider, req.Code)
	if err != nil {
		if errors.Is(err, ErrProviderUnsupported) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported provider"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	roles := []string{"user"}
	token, _, err := h.service.CreateSession(c.Request.Context(), userID, roles, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "session creation failed"})
		return
	}

	h.setAuthCookie(c, token, int(24*time.Hour.Seconds()))

	c.JSON(http.StatusOK, openapi.PostAuthProviderCallback200JSONResponse{
		Status: ptr("success"),
		User: &struct {
			Email *string `json:"email,omitempty"`
			Id    *int64  `json:"id,omitempty"`
		}{
			Id:    &userID,
			Email: &email,
		},
	})
}

func (h *Handler) GetAuthLogout(c *gin.Context) {
	h.logout(c)
	c.JSON(http.StatusOK, struct{}{})
}

func (h *Handler) PostAuthLogout(c *gin.Context) {
	h.logout(c)
	c.JSON(http.StatusOK, openapi.PostAuthLogout200JSONResponse{Status: ptr("logged out")})
}

func (h *Handler) logout(c *gin.Context) {
	cookie, err := c.Cookie("session")
	if err == nil {
		_ = h.service.DeleteSession(c.Request.Context(), cookie)
	}
	h.clearAuthCookie(c)
}

func (h *Handler) GetMe(c *gin.Context) {
	sessionVal, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*Session)

	user, err := h.service.GetUserByID(c.Request.Context(), session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	var roles *[]string
	if len(session.Roles) > 0 {
		roles = &session.Roles
	}
	c.JSON(http.StatusOK, openapi.GetAuthMe200JSONResponse{
		User: &struct {
			Email        *string `json:"email"`
			Id           *int64  `json:"id,omitempty"`
			Phone        *string `json:"phone"`
			ProfileImage *string `json:"profile_image"`
		}{
			Id:           &user.ID,
			Email:        user.Email,
			Phone:        user.Phone,
			ProfileImage: user.ProfileImage,
		},
		Session: &struct {
			ExpiresAt *time.Time `json:"expires_at,omitempty"`
			Roles     *[]string  `json:"roles,omitempty"`
		}{
			ExpiresAt: &session.ExpiresAt,
			Roles:     roles,
		},
	})
}

func (h *Handler) PostProfileImage(c *gin.Context) {
	sessionVal, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*Session)

	var req openapi.PostProfileImageJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	imageURL, err := h.service.UploadProfileImage(c.Request.Context(), req.Image)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrInvalidImage) || errors.Is(err, ErrImageTooLarge) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateUserProfileImage(c.Request.Context(), session.UserID, imageURL); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, openapi.PostProfileImage200JSONResponse{ImageUrl: &imageURL})
}

func (h *Handler) setAuthCookie(c *gin.Context, token string, maxAge int) {
	c.SetCookie("session", token, maxAge, "/", h.cookieDomain, h.cookieSecure, true)
}

func (h *Handler) clearAuthCookie(c *gin.Context) {
	c.SetCookie("session", "", -1, "/", h.cookieDomain, h.cookieSecure, true)
}

func ptr[T any](v T) *T {
	return &v
}
