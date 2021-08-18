package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	JSON_SUCCESS int = 200
	JSON_ERROR   int = 400
)

func Response(ctx *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

func Success(ctx *gin.Context, data interface{}, msg string) {
	Response(ctx, http.StatusOK, JSON_SUCCESS, data, msg)
}

func Fail(ctx *gin.Context, msg string, data interface{}) {
	Response(ctx, http.StatusOK, JSON_ERROR, data, msg)
}
