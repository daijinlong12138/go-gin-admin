package auth

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"go-gin-admin/model"
	"go-gin-admin/response"
	"go-gin-admin/util"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func PermissionNew(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")              //权限名称
	slug := ctx.PostForm("slug")              //标志(需要保证唯一)
	HttpMethod := ctx.PostForm("http_method") //方法(为空默认为所有方法)
	HttpPath := ctx.PostForm("http_path")     //路径(路径不包含全局路由前缀，换行多条)

	//校验
	if len(name) == 0 {
		response.Fail(ctx, "权限名称不能为空", nil)
		return
	}
	if len(slug) == 0 {
		response.Fail(ctx, "标志不能为空", nil)
		return
	}
	if len(HttpMethod) != 0 {
		HttpMethodArr := strings.Split(HttpMethod, `,`)
		var AllMethodArr = []string{"GET", "POST", "PUT", "DELETE"}
		for _, method := range HttpMethodArr {
			if !util.IsContainStr(AllMethodArr, method) {
				response.Fail(ctx, "方法不存在："+method, nil)
				return
			}
		}
	}
	HttpPath = strings.Trim(HttpPath, " ")

	if isHavePermissionSlug(DB, slug) {
		response.Fail(ctx, "标志已存在", nil)
		return
	}

	newPermissions := model.AdminPermissions{
		Name:       name,
		Slug:       slug,
		HttpMethod: HttpMethod,
		HttpPath:   HttpPath,
	}
	if err := DB.Create(&newPermissions).Error; err != nil {
		response.Fail(ctx, "权限增加失败", nil)
		return
	}
	response.Success(ctx, nil, "权限增加成功")
}

func PermissionDetail(ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := ctx.PostForm("id") //权限ID
	var permission = model.AdminPermissions{}

	DB.Where("id = ?", id).Find(&permission)
	if permission.ID == 0 {
		response.Fail(ctx, "不存在", nil)
		return
	}
	response.Success(ctx, model.ToAdminPermissionsTodo(permission), "成功")

}

func PermissionDelete(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id") //ID
	var permission = model.AdminPermissions{}
	Id, _ := strconv.Atoi(id)
	permission = findPermissionById(DB, Id)
	if permission.ID != uint(Id) {
		response.Fail(c, "不存在", nil)
		return
	}

	tx := DB.Begin()
	err := tx.Where("id = ?", Id).Unscoped().Delete(&model.AdminPermissions{}).Error
	if err != nil {
		common.LogError(c, "删除 AdminPermissions 失败: "+err.Error())
		response.Fail(c, "删除 AdminPermissions 失败", nil)
		tx.Rollback()
		return
	}
	err = tx.Where("permission_id = ?", Id).Unscoped().Delete(&model.AdminUserPermissions{}).Error
	if err != nil {
		common.LogError(c, "删除 AdminUserPermissions 失败: "+err.Error())
		response.Fail(c, "删除 AdminUserPermissions 失败", nil)
		tx.Rollback()
		return
	}

	err = tx.Where("permission_id = ?", Id).Unscoped().Delete(&model.AdminRolePermissions{}).Error
	if err != nil {
		common.LogError(c, "删除 AdminRolePermissions 失败: "+err.Error())
		response.Fail(c, "删除 AdminRolePermissions 失败", nil)
		tx.Rollback()
		return
	}
	tx.Commit()

	response.Success(c, nil, "成功")
}

func PermissionInfo(c *gin.Context) {
	var permission model.AdminPermissions
	var total int64
	var permissions []model.AdminPermissions
	var _permissions []model.AdminPermissionsTodo

	db := common.GetDB()
	//获取参数
	page, _ := strconv.Atoi(c.PostForm("page"))         //页数
	pageSize, _ := strconv.Atoi(c.PostForm("pageSize")) //每页数

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	db = db.Model(permission) //查询对应的数据库表

	if name, isExist := c.GetPostForm("name"); isExist == true && name != "" {
		db = db.Where("name = ?", name)
	}
	if slug, isExist := c.GetPostForm("slug"); isExist == true && slug != "" {
		db = db.Where("slug = ?", slug)
	}

	if err := db.Count(&total).Error; err != nil {
		response.Fail(c, "查询数据异常", nil)
		return
	}

	db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	if err := db.Find(&permissions).Error; err != nil {
		common.LogError(c, "查询数据异常 : "+err.Error())
		response.Fail(c, "查询数据异常", nil)
		return
	}

	// 格式化
	for _, item := range permissions {
		_permissions = append(_permissions, model.ToAdminPermissionsTodo(item))
	}

	//fmt.Println(roles)
	response.Success(c, gin.H{
		"data":     _permissions,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, "成功")
}

func PermissionEdit(c *gin.Context) {
	DB := common.GetDB()
	var permissions = model.AdminPermissions{}
	//获取参数
	id := c.PostForm("id")                  //角色ID
	name := c.PostForm("name")              //权限
	slug := c.PostForm("slug")              //标志(需要保证唯一)
	HttpMethod := c.PostForm("http_method") //方法(为空默认为所有方法)
	HttpPath := c.PostForm("http_path")     //路径(路径不包含全局路由前缀，换行多条)

	if len(name) == 0 {
		response.Fail(c, "权限不能为空", nil)
		return
	}
	if len(slug) == 0 {
		response.Fail(c, "标志不能为空", nil)
		return
	}
	if len(HttpMethod) != 0 {
		HttpMethodArr := strings.Split(HttpMethod, `,`)
		var AllMethodArr = []string{"GET", "POST", "PUT", "DELETE"}
		for _, method := range HttpMethodArr {
			if !util.IsContainStr(AllMethodArr, method) {
				response.Fail(c, "方法不存在："+method, nil)
				return
			}
		}
	}
	HttpPath = strings.Trim(HttpPath, " ")

	Id, _ := strconv.Atoi(id)
	permissions = findPermissionById(DB, Id)
	if permissions.ID != uint(Id) {
		response.Fail(c, "不存在权限", nil)
		return
	}
	if permissions.Slug != slug && isHavePermissionSlug(DB, slug) {
		response.Fail(c, "标志已存在", nil)
		return
	}
	data := make(map[string]interface{})
	data["name"] = name
	data["slug"] = slug
	data["http_method"] = HttpMethod
	data["http_path"] = HttpPath

	if err := DB.Model(&model.AdminPermissions{}).Where("id = ?", id).Updates(&data).Error; err != nil {
		common.LogError(c, "更新失败 : "+err.Error())
		response.Fail(c, "更新失败", nil)
		return
	}
	response.Success(c, nil, "成功")
}

func PermissionAllInfo(ctx *gin.Context) {
	db := common.GetDB()
	var data []model.AdminPermissions
	var _data []model.AdminPermissionsTodo
	db.Find(&data)
	// 格式化
	for _, item := range data {
		_data = append(_data, model.ToAdminPermissionsTodo(item))
	}
	response.Success(ctx, _data, "成功")
}

func isHavePermissionSlug(db *gorm.DB, slug string) bool {
	var permissions model.AdminPermissions
	err := db.Where("slug = ?", slug).Find(&permissions).Error
	if err == nil && permissions.ID != 0 {
		return true
	}
	return false
}

func findPermissionById(db *gorm.DB, id int) model.AdminPermissions {
	var permissions model.AdminPermissions
	db.Where("id = ?", id).Find(&permissions)
	return permissions
}
