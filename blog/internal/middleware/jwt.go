package middleware

import (
	"blog/internal/pkg/resp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

const Token = "token"

var NoAuthUri = map[string]bool{
	"/api/v1/register":                  true,
	"/api/v1/login/:username/:password": true,
	"/healthz":                          true,
}

func JWTAuth(log *zap.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		path := c.FullPath()
		log.Info("jwt middleware start", zap.String("path", path))
		token := c.GetHeader(Token)
		log.Info("path:", zap.Any("bingo:", NoAuthUri[path]))
		if !NoAuthUri[path] {

			if len(token) == 0 {
				resp.Fail(c, http.StatusUnauthorized, "missing_token")
				c.Abort()
				return
			}
		}
		log.Info("token auth pass", zap.String("token", token))
		//tok, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		//	return []byte(secret), nil
		//})
		//if err != nil || !tok.Valid {
		//	resp.Fail(c, http.StatusUnauthorized, "invalid_token")
		//	c.Abort()
		//	return
		//}
		//
		//claims, ok := tok.Claims.(jwt.MapClaims)
		//if !ok {
		//	resp.Fail(c, http.StatusUnauthorized, "invalid_claims")
		//	c.Abort()
		//	return
		//}
		//if exp, ok := claims["exp"].(float64); ok {
		//	if time.Now().Unix() > int64(exp) {
		//		resp.Fail(c, http.StatusUnauthorized, "token_expired")
		//		c.Abort()
		//		return
		//	}
		//}
		//if uid, ok := claims["sub"].(string); ok {
		//	c.Set(CtxUserIDKey, uid)
		//}
		c.Next()

	}
}
