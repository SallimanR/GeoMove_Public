package notification

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"monolith/internal/auth"
	"monolith/internal/notification/sqlc"
)

type Handler struct {
	store       *Store
	sender      *WebPushService
	vapidPubKey string
}

func NewHandler(store *Store, sender *WebPushService, vapidPubKey string) *Handler {
	return &Handler{
		store:       store,
		sender:      sender,
		vapidPubKey: vapidPubKey,
	}
}

func (h *Handler) GetVapidPublicKey(ctx context.Context, request GetVapidPublicKeyRequestObject) (GetVapidPublicKeyResponseObject, error) {
	return GetVapidPublicKey200JSONResponse{
		PublicKey: h.vapidPubKey,
	}, nil
}

func (h *Handler) Subscribe(ctx context.Context, request SubscribeRequestObject) (SubscribeResponseObject, error) {
	gctx := ctx.(*gin.Context)
	userVal, exists := gctx.Get("user")
	if !exists {
		return nil, fmt.Errorf("unauthorized")
	}
	user := userVal.(*auth.User)

	deviceType := "web"
	if request.Body.DeviceType != nil {
		deviceType = *request.Body.DeviceType
	}

	err := h.store.UpsertSubscription(ctx, sqlc.UpsertSubscriptionParams{
		UserID:          user.ID,
		Endpoint:        request.Body.Endpoint,
		DevicePublicKey: request.Body.DevicePublicKey,
		AuthSecret:      request.Body.AuthSecret,
		DeviceType:      deviceType,
	})
	if err != nil {
		log.Printf("subscribe: save error: %v", err)
		return nil, fmt.Errorf("failed to save subscription: %w", err)
	}

	return Subscribe200Response{}, nil
}

func (h *Handler) Unsubscribe(ctx context.Context, request UnsubscribeRequestObject) (UnsubscribeResponseObject, error) {
	err := h.store.Delete(ctx, request.Body.Endpoint)
	if err != nil {
		log.Printf("unsubscribe: delete error: %v", err)
		return nil, fmt.Errorf("failed to delete subscription: %w", err)
	}

	return Unsubscribe200Response{}, nil
}

type TestPushRequest struct {
	UserID int64  `json:"userId" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
}

func (h *Handler) TestPush(c *gin.Context) {
	var req TestPushRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.sender.SendToUser(c.Request.Context(), req.UserID, NotificationPayload{
		Title: req.Title,
		Body:  req.Body,
	})
	if err != nil {
		log.Printf("test push error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "sent"})
}

func RegisterNotificationRoutes(api *gin.RouterGroup, h *Handler, authMiddleware gin.HandlerFunc) {
	strictHandler := NewStrictHandler(h, nil)
	wrapper := ServerInterfaceWrapper{Handler: strictHandler}

	api.GET("/notifications/vapid-public-key", wrapper.GetVapidPublicKey)

	protected := api.Group("/notifications")
	protected.Use(authMiddleware)
	{
		protected.POST("/subscribe", wrapper.Subscribe)
		protected.POST("/unsubscribe", wrapper.Unsubscribe)
		protected.POST("/test-push", h.TestPush)
	}
}

func (h *Handler) SubscribeErrorHandler(c *gin.Context, err error, statusCode int) {
	log.Printf("subscribe error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}

func (h *Handler) UnsubscribeErrorHandler(c *gin.Context, err error, statusCode int) {
	log.Printf("unsubscribe error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}

func (h *Handler) GetVapidPublicKeyErrorHandler(c *gin.Context, err error, statusCode int) {
	log.Printf("vapid public key error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}
