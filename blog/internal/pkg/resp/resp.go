package resp

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(200, Response{Code: 0, Message: "ok", Data: data})
}

func Fail(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, Response{Code: 1, Message: msg})
}
