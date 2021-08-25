package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-admin/package/mylogger"
	"time"
)

func CommonMiddleware() gin.HandlerFunc {

	path := viper.GetString("Logs.path")
	Log := mylogger.NewFileLogger(path, "access.log")
	return func(c *gin.Context) {

		//引入日志文件，详细记录访问ip，访问地址及参数
		//方法前
		// 开始时间
		start := time.Now()
		c.Next()
		// 结束时间
		end := time.Now()
		//方法后
		//执行时间
		latency := end.Sub(start)

		URLpath := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		Log.Info("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, URLpath,
		)
	}
}
