package middleware

import (
	"net/http"
	"runtime/debug"

	"blog/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("request_id", GetRequestID(c)),
					zap.ByteString("stack", debug.Stack()),
				)
				resp.Fail(c, http.StatusInternalServerError, "internal_error")
			}
		}()
		c.Next()
	}
}
