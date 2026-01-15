package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-Id"
const CtxRequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(RequestIDHeader)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set(CtxRequestIDKey, rid)
		c.Writer.Header().Set(RequestIDHeader, rid)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	v, ok := c.Get(CtxRequestIDKey)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
