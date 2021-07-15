package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
)


func main() {

	common.InitDB()

	// 1.创建路由
	router := gin.Default()

	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	router = CollectRoute(router)

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	fmt.Printf("start server (port:%s)", "8000")
	router.Run(":8000")

}


