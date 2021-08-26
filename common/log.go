package common

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-admin/package/mylogger"
)

var Log *mylogger.FileLogger

func InitLog() *mylogger.FileLogger {
	path := viper.GetString("Logs.path")
	name := viper.GetString("Logs.name")
	Log = mylogger.NewFileLogger(path, name)
	return Log
}

func GetLog() *mylogger.FileLogger {
	return Log
}

//方法中描述--加入访问地址--不返回到界面的
func LogInfo(ctx *gin.Context, msg string, arg ...interface{}) {
	URLpath := ctx.Request.URL.Path + " "
	Log.Info(URLpath+msg, arg...)
}

func LogError(ctx *gin.Context, msg string, arg ...interface{}) {
	URLpath := ctx.Request.URL.Path + " "
	Log.Error(URLpath+msg, arg...)
}
