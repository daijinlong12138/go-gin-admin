package response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"net/http"
)

const (
	JSON_SUCCESS int = 200
	JSON_ERROR   int = 400
	JSON_AUTH    int = 401
)

func Response(ctx *gin.Context, httpStatus int, code int, data interface{}, msg string) {
	h := gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	}
	ctx.JSON(httpStatus, h)
	str, err := json.Marshal(h)
	if err != nil {
		common.Log.Error(err.Error())
	}
	URLpath := ctx.Request.URL.Path
	common.Log.Info(URLpath + " Response: " + string(str))
}

func AuthFail(ctx *gin.Context, data interface{}, msg string) {
	Response(ctx, http.StatusUnauthorized, JSON_AUTH, data, msg)
}

func Success(ctx *gin.Context, data interface{}, msg string) {
	Response(ctx, http.StatusOK, JSON_SUCCESS, data, msg)
}

func Fail(ctx *gin.Context, msg string, data interface{}) {
	Response(ctx, http.StatusOK, JSON_ERROR, data, msg)
}
