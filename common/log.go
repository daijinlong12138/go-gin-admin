package common

import (
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
