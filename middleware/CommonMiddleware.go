package middleware

import (
	"github.com/gin-gonic/gin"
)

func CommonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		//请求处理
		c.Next()


	}
}