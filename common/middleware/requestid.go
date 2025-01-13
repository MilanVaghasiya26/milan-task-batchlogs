package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/team-scaletech/common/logging"
)

// DefaultRequestId is a middleware that injects a requestid into the context and header of each request.
func DefaultRequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := uuid.NewString()
		c.Request.Header.Set(logging.RequestIDKey, reqId)
		c.Set(logging.RequestIDKey, reqId)
		c.Next()
	}
}
