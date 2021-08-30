package main

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/contronller"
	"go-gin-admin/contronller/auth"
	"go-gin-admin/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	//版本1--使用公共中间件，使用接口鉴权--登录不需要
	v1 := r.Group("admin").Use(middleware.ApiInterfaceAuthCheck())
	v1.POST("/roles/new", auth.RolesNew)                   //角色新增接口
	v1.POST("/roles/detail", auth.RolesDetail)             //角色详情接口
	v1.POST("/roles/delete", auth.RolesDelete)             //角色删除接口
	v1.POST("/roles/info", auth.RolesInfo)                 //角色管理展示接口
	v1.POST("/roles/edit", auth.RolesEdit)                 //角色编辑接口
	v1.POST("/roles/allinfo", auth.RolesAllinfo)           //全部角色展示接口
	v1.POST("/permission/new", auth.PermissionNew)         //权限新增接口
	v1.POST("/permission/detail", auth.PermissionDetail)   //权限详情接口
	v1.POST("/permission/delete", auth.PermissionDelete)   //权限删除接口
	v1.POST("/permission/info", auth.PermissionInfo)       //权限管理展示接口
	v1.POST("/permission/edit", auth.PermissionEdit)       //权限编辑接口
	v1.POST("/permission/allinfo", auth.PermissionAllInfo) //全部权限展示接口
	v1.POST("/uploadimg", contronller.UploadImg)           //上传文件
	v1.POST("/manager/new", auth.ManagerNew)               //用户新增接口
	v1.POST("/manager/detail", auth.ManagerDetail)         //用户详情接口
	v1.POST("/manager/delete", auth.ManagerDelete)         //用户删除接口
	v1.POST("/manager/info", auth.ManagerInfo)             //用户管理展示接口
	v1.POST("/manager/edit", auth.ManagerEdit)             //用户编辑接口
	v1.POST("/manager/login", auth.ManagerLogin)           //用户登录接口
	v1.POST("/menu/new", auth.MenuNew)                     //菜单新增接口
	v1.POST("/menu/detail", auth.MenuDetail)               //菜单详情接口
	v1.POST("/menu/delete", auth.MenuDelete)               //菜单删除接口
	v1.POST("/menu/info", auth.MenuInfo)                   //菜单管理展示接口
	v1.POST("/menu/edit", auth.MenuEdit)                   //菜单编辑接口
	v1.POST("/oplog/info", auth.OpLogsInfo)                //日志列表获取接口

	return r
}
