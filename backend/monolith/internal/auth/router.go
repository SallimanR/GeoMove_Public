package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterHTTPRoutes(r *gin.RouterGroup, service *Service, authMiddleware gin.HandlerFunc) {
	handler := NewHandler(service)

	auth := r.Group("/auth")
	auth.POST("/:provider/callback", handler.PostAuthProviderCallback)
	auth.POST("/logout", handler.PostAuthLogout)
	auth.GET("/logout", handler.GetAuthLogout)

	protected := r.Group("/")
	protected.Use(authMiddleware)
	protected.GET("/auth/me", handler.GetMe)
	protected.POST("/profile/image", handler.PostProfileImage)
}
