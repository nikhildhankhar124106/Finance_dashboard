package middleware

import (
	"bytes"
	"io"
	"net/http"

	"backend/domain/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditMiddleware logs all write operations (POST, PUT, PATCH, DELETE) to the ActivityLog table.
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only log write methods natively over state changes
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodOptions || method == http.MethodHead {
			c.Next()
			return
		}

		// Read the body for detail mapping (need to restore it for the handler)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Process the request first
		c.Next()

		// Only log successful or specifically authorized attempts to keep audits clean
		if c.Writer.Status() >= 400 && c.Writer.Status() != http.StatusConflict {
			return
		}

		userIDVal, _ := c.Get("user_id")
		userID, _ := userIDVal.(uint)

		log := models.ActivityLog{
			UserID:    userID,
			Action:    method,
			Resource:  c.Request.URL.Path,
			Details:   string(bodyBytes),
			IPAddress: c.ClientIP(),
		}

		// Background saving to prevent blocking response times
		go func(l models.ActivityLog) {
			db.Create(&l)
		}(log)
	}
}
