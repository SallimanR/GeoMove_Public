package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

func Success(c *gin.Context, data any) {
	c.JSON(200, Response{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Success: false,
		Data:    message,
	})
}
