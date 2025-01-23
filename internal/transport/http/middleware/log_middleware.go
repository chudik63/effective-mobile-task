package middleware

import (
	"effective-mobile-task/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LogMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		queryParams := c.Request.URL.Query()

		logger.Debug(c.Request.Context(), "Incoming request", zap.String("method", method), zap.String("path", path), zap.Any("query params", queryParams))

		c.Next()

		logger.Debug(c.Request.Context(), "Request completed", zap.Int("status", c.Writer.Status()))
	}
}
