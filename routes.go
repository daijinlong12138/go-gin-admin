package main

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/contronller"
)

func CollectRoute(r *gin.Engine) *gin.Engine  {

	r.POST("/api/auth/register", contronller.Register)

	return r
}