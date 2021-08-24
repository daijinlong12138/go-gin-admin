package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func CommonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//引入日志文件，详细记录访问ip，访问地址及参数
		//方法前
		fmt.Println("方法前")
		c.Next()
		//方法后
		fmt.Println("方法后")
	}
}
