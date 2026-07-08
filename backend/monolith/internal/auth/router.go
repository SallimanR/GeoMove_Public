package auth

import (
	"monolith/internal/auth/openapi"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPRoutes(r *gin.RouterGroup, service *Service) {
	handler := NewHandler(service)
	openapi.RegisterHandlers(r, handler)
}
