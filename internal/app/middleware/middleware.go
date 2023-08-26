package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func JSONMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=UTF-8")
		c.Next()
	})
}

func TimeoutMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})
}

// LoggingMiddleware - a handy middleware function that logs out incoming requests
func LoggingMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		log.WithFields(
			log.Fields{
				"Method": c.Request.Method,
				"Path":   c.Request.URL.Path,
			}).
			Info("handled request")
		c.Next()
	})
}
