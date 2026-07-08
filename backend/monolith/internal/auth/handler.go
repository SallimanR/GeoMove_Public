package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"monolith/internal/auth/openapi"
)

type Handler struct {
	service      *Service
	cookieSecure bool
	cookieDomain string
}

func NewHandler(service *Service) openapi.ServerInterface {
	return &Handler{
		service:      service,
		cookieSecure: os.Getenv("COOKIE_SECURE") == "true",
		cookieDomain: os.Getenv("COOKIE_DOMAIN"), // empty = current domain
	}
}

func (h *Handler) PostAuthProviderCallback(c *gin.Context, params openapi.PostAuthProviderCallbackParamsProvider) {
	provider := string(params)

	var req struct {
		Code string `json:"code"`
	}
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user": gin.H{
			"id":    userID,
			"email": email,
		},
	})
}

func (h *Handler) PostAuthLogout(c *gin.Context) {
	h.logout(c)
}

func (h *Handler) GetAuthLogout(c *gin.Context) {
	h.logout(c)
}

func (h *Handler) logout(c *gin.Context) {
	cookie, err := c.Cookie("session")
	if err == nil {
		_ = h.service.DeleteSession(c.Request.Context(), cookie)
	}
	h.clearAuthCookie(c)
	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}

func (h *Handler) PostProfileImage(c *gin.Context) {
	sessionVal, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	session := sessionVal.(*Session)

	var req struct {
		Image string `json:"image"`
	}
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

	c.JSON(http.StatusOK, gin.H{"image_url": imageURL})
}

func (h *Handler) setAuthCookie(c *gin.Context, token string, maxAge int) {
	c.SetCookie(
		"session",
		token,
		maxAge,
		"/",
		h.cookieDomain,
		h.cookieSecure,
		true,
	)
}

func (h *Handler) clearAuthCookie(c *gin.Context) {
	c.SetCookie(
		"session",
		"",
		-1,
		"/",
		h.cookieDomain,
		h.cookieSecure,
		true,
	)
}
