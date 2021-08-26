package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-admin/common"
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
		method := c.Request.Method
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		Header := c.Request.Header
		Log.Info("| %3d | %13v | %15s | %s  %s |  %s|",
			statusCode,
			latency,
			clientIP,
			method, URLpath,
			Header,
		)

		//日志中打印参数
		switch method {
		case "POST":
			tmp := make(map[string]interface{})
			for k, v := range c.Request.PostForm {
				tmp[k] = v
			}
			strs, err := json.Marshal(tmp)
			if err != nil {
				common.Log.Error(URLpath + " " + err.Error())
			}
			common.Log.Info(URLpath + " Param: " + string(strs))
			break
		case "GET":
			param := c.Request.URL.RawQuery
			common.Log.Info(URLpath + " Param: " + param)
			break
		default:
			common.Log.Info(URLpath + " 方法未配置无法获取参数: " + method)
			break
		}

	}
}
