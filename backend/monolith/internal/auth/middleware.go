package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "войдите в аккаунт"})
			return
		}

		session, err := s.ValidateToken(c.Request.Context(), cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "сессия недействительна или истекла"})
			return
		}

		user, err := s.GetUserByID(c.Request.Context(), session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден"})
			return
		}

		c.Set("session", session)
		c.Set("user", user)

		c.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "не авторизован"})
			return
		}
		user := userVal.(*User)

		for _, r := range user.Roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "доступно только водителям"})
	}
}
