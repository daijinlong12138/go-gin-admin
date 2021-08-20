package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-admin/common"
	"os"
)

func main() {
	InitConfig()
	common.InitDB()

	// 1.创建路由
	router := gin.Default()

	//router.Use(middleware.CommonMiddleware())

	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	router = CollectRoute(router)

	prot := viper.GetString("server.port")
	if prot != "" {
		// Run("里面不指定端口号默认为8080")
		fmt.Printf("start server (port:%s)", prot)
		panic(router.Run(":" + prot))
	}
	// 3.监听端口，默认在8080
	panic(router.Run())
}

func InitConfig() {
	workDir, _ := os.Getwd()
	//设置要读取的文件名
	viper.SetConfigName("application")
	//设置文件读取类型
	viper.SetConfigType("yml")
	//设置文件的路径
	viper.AddConfigPath(workDir + "/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
