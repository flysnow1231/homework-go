package middleware

import (
	"net/http"
	"strings"
	"time"

	"blog/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const CtxUserIDKey = "user_id"

func JWTAuth(secret string) gin.HandlerFunc {
	if secret == "" {
		secret = "change-me"
	}
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			resp.Fail(c, http.StatusUnauthorized, "missing_token")
			c.Abort()
			return
		}
		tokenStr := strings.TrimSpace(h[len("Bearer "):])

		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			resp.Fail(c, http.StatusUnauthorized, "invalid_token")
			c.Abort()
			return
		}

		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			resp.Fail(c, http.StatusUnauthorized, "invalid_claims")
			c.Abort()
			return
		}
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				resp.Fail(c, http.StatusUnauthorized, "token_expired")
				c.Abort()
				return
			}
		}
		if uid, ok := claims["sub"].(string); ok {
			c.Set(CtxUserIDKey, uid)
		}
		c.Next()
	}
}
