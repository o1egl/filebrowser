package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/filebrowser/filebrowser/v3/log"
)

func Logger(c *gin.Context) {
	// Start timer
	start := time.Now()

	// Process request
	c.Next()

	log.WithContext(c.Request.Context()).WithFields(log.Fields{
		"client_ip":        c.ClientIP(),
		"request_method":   c.Request.Method,
		"request_path":     c.Request.URL.Path,
		"request_duration": time.Since(start).String(),
		"response_status":  c.Writer.Status(),
		"response_size":    c.Writer.Size(),
	}).Infof("")
}
