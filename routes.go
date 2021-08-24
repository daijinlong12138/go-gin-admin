package main

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/contronller"
	"go-gin-admin/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {

	//版本1--使用公共中间件，使用接口鉴权--登录不需要
	v1 := r.Group("auth").Use(middleware.ApiInterfaceAuthCheck())
	v1.POST("/admin/roles/new", contronller.RolesNew)                   //角色新增接口
	v1.POST("/admin/roles/detail", contronller.RolesDetail)             //角色详情接口
	v1.POST("/admin/roles/delete", contronller.RolesDelete)             //角色删除接口
	v1.POST("/admin/roles/info", contronller.RolesInfo)                 //角色管理展示接口
	v1.POST("/admin/roles/edit", contronller.RolesEdit)                 //角色编辑接口
	v1.POST("/admin/roles/allinfo", contronller.RolesAllinfo)           //全部角色展示接口
	v1.POST("/admin/permission/new", contronller.PermissionNew)         //权限新增接口
	v1.POST("/admin/permission/detail", contronller.PermissionDetail)   //权限详情接口
	v1.POST("/admin/permission/delete", contronller.PermissionDelete)   //权限删除接口
	v1.POST("/admin/permission/info", contronller.PermissionInfo)       //权限管理展示接口
	v1.POST("/admin/permission/edit", contronller.PermissionEdit)       //权限编辑接口
	v1.POST("/admin/permission/allinfo", contronller.PermissionAllinfo) //全部权限展示接口
	v1.POST("/admin/uploadimg", contronller.UploadImg)                  //上传文件
	v1.POST("/admin/manager/new", contronller.ManagerNew)               //用户新增接口
	v1.POST("/admin/manager/detail", contronller.ManagerDetail)         //用户详情接口
	v1.POST("/admin/manager/delete", contronller.ManagerDelete)         //用户删除接口
	v1.POST("/admin/manager/info", contronller.ManagerInfo)             //用户管理展示接口
	v1.POST("/admin/manager/edit", contronller.ManagerEdit)             //用户编辑接口
	v1.POST("/admin/menu/new", contronller.MenuNew)                     //菜单新增接口
	v1.POST("/admin/menu/detail", contronller.MenuDetail)               //菜单详情接口
	v1.POST("/admin/menu/delete", contronller.MenuDelete)               //菜单删除接口
	v1.POST("/admin/menu/info", contronller.MenuInfo)                   //菜单管理展示接口
	v1.POST("/admin/menu/edit", contronller.MenuEdit)                   //菜单编辑接口
	v1.POST("/admin/oplog/info", contronller.OpLogsInfo)                //日志列表获取接口

	r.POST("/api/auth/register", contronller.Register)
	r.POST("/api/auth/login", contronller.Login)
	r.POST("/api/auth/info", middleware.AuthMiddleware(), contronller.Info)

	r.POST("/admin/roles/new", contronller.RolesNew)         //角色新增接口
	r.POST("/admin/roles/detail", contronller.RolesDetail)   //角色详情接口
	r.POST("/admin/roles/delete", contronller.RolesDelete)   //角色删除接口
	r.POST("/admin/roles/info", contronller.RolesInfo)       //角色管理展示接口
	r.POST("/admin/roles/edit", contronller.RolesEdit)       //角色编辑接口
	r.POST("/admin/roles/allinfo", contronller.RolesAllinfo) //全部角色展示接口

	r.POST("/admin/permission/new", contronller.PermissionNew)         //权限新增接口
	r.POST("/admin/permission/detail", contronller.PermissionDetail)   //权限详情接口
	r.POST("/admin/permission/delete", contronller.PermissionDelete)   //权限删除接口
	r.POST("/admin/permission/info", contronller.PermissionInfo)       //权限管理展示接口
	r.POST("/admin/permission/edit", contronller.PermissionEdit)       //权限编辑接口
	r.POST("/admin/permission/allinfo", contronller.PermissionAllinfo) //全部权限展示接口

	r.POST("/admin/uploadimg", contronller.UploadImg)          //上传文件
	r.POST("/admin/manager/new", contronller.ManagerNew)       //用户新增接口
	r.POST("/admin/manager/detail", contronller.ManagerDetail) //用户详情接口
	r.POST("/admin/manager/delete", contronller.ManagerDelete) //用户删除接口
	r.POST("/admin/manager/info", contronller.ManagerInfo)     //用户管理展示接口
	r.POST("/admin/manager/edit", contronller.ManagerEdit)     //用户编辑接口
	r.POST("/admin/manager/login", contronller.ManagerLogin)   //用户登录接口

	r.POST("/admin/menu/new", contronller.MenuNew)       //菜单新增接口
	r.POST("/admin/menu/detail", contronller.MenuDetail) //菜单详情接口
	r.POST("/admin/menu/delete", contronller.MenuDelete) //菜单删除接口
	r.POST("/admin/menu/info", contronller.MenuInfo)     //菜单管理展示接口
	r.POST("/admin/menu/edit", contronller.MenuEdit)     //菜单编辑接口

	r.POST("/admin/oplog/info", contronller.OpLogsInfo) //日志列表获取接口
	return r
}
