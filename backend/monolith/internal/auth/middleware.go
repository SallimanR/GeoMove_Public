package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing session cookie"})
			return
		}

		session, err := s.ValidateToken(c.Request.Context(), cookie)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			return
		}

		user, err := s.GetUserByID(c.Request.Context(), session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set("session", session)
		c.Set("user", user)

		c.Next()
	}
}
