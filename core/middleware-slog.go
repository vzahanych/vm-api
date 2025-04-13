package core

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// Custom Gin Logger using slog for internal logging
func CustomGinLogger(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Before request: Log incoming request
		logger.Info("Incoming request",
			"method", c.Request.Method,
			"url", c.Request.URL.String(),
			"remote_ip", c.ClientIP(),
		)

		// Store the logger in the context
		c.Set("logger", logger)

		// Process request
		c.Next()

		// After request: Log the response status
		logger.Info("Response",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"url", c.Request.URL.String(),
		)
	}
}