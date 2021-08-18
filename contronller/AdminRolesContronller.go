package contronller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"go-gin-admin/model"
	"go-gin-admin/response"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

func RolesNew(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := c.PostForm("name") //角色
	slug := c.PostForm("slug") //标志(需要保证唯一)

	//校验
	if len(name) == 0 {
		response.Fail(c, "角色不能为空", nil)
		return
	}
	if len(slug) == 0 {
		response.Fail(c, "标志不能为空", nil)
		return
	}

	if isHaveRolesSlug(DB, slug) {
		response.Fail(c, "标志已存在", nil)
		return
	}

	newRoles := model.AdminRoles{
		Name: name,
		Slug: slug,
	}
	DB.Create(&newRoles)

	log.Println(name, slug)

	response.Success(c, nil, "角色增加成功")
}

func RolesDetail(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id") //角色ID
	var roles = model.AdminRoles{}

	DB.Where("id = ?", id).Find(&roles)
	if roles.ID == 0 {
		response.Fail(c, "不存在", nil)
		return
	}
	log.Println(roles)
	response.Success(c, model.ToAdminRolesTodo(roles), "成功")

}

func RolesDelete(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id") //角色ID
	var roles = model.AdminRoles{}
	Id, _ := strconv.Atoi(id)
	roles = findRolesById(DB, Id)
	if roles.ID != uint(Id) {
		response.Fail(c, "不存在角色", nil)
		return
	}
	tx := DB.Begin()
	err := tx.Where("id = ?", Id).Unscoped().Delete(&model.AdminRoles{}).Error
	if err != nil {
		response.Fail(c, "角色删除失败", nil)
		tx.Rollback()
		return
	}
	err = tx.Where("role_id = ?", Id).Unscoped().Delete(&model.AdminRolePermissions{}).Error
	if err != nil {
		response.Fail(c, "权限关联删除失败", nil)
		tx.Rollback()
		return
	}
	tx.Commit()
	response.Success(c, nil, "成功")
}

func RolesInfo(c *gin.Context) {
	var role model.AdminRoles
	var total int64
	roles := make([]model.AdminRoles, 0)
	var _roles []model.AdminRolesTodo

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
	db = db.Model(role) //查询对应的数据库表

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
	if err := db.Find(&roles).Error; err != nil {
		fmt.Println(err.Error())
	}

	// 格式化
	for _, item := range roles {
		_roles = append(_roles, model.ToAdminRolesTodo(item))
	}

	//fmt.Println(roles)
	response.Success(c, gin.H{
		"data":     _roles,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, "成功")
}

func RolesEdit(c *gin.Context) {
	DB := common.GetDB()
	var roles = model.AdminRoles{}
	//获取参数
	id := c.PostForm("id")                         //角色ID
	name := c.PostForm("name")                     //角色
	slug := c.PostForm("slug")                     //标志(需要保证唯一)
	permission_ids := c.PostForm("permission_ids") //权限以逗号分割

	if len(name) == 0 {
		response.Fail(c, "角色不能为空", nil)
		return
	}
	if len(slug) == 0 {
		response.Fail(c, "标志不能为空", nil)
		return
	}

	var permissionIdArr []string
	if len(permission_ids) != 0 {
		permissionIdArr = strings.Split(permission_ids, `,`)
	}

	Id, _ := strconv.Atoi(id)
	roles = findRolesById(DB, Id)
	if roles.ID != uint(Id) {
		response.Fail(c, "不存在角色", nil)
		return
	}
	if roles.Slug != slug && isHaveRolesSlug(DB, slug) {
		response.Fail(c, "标志已存在", nil)
		return
	}
	if len(permissionIdArr) > 0 {
		for _, permissionId := range permissionIdArr {
			permissionIdint, _ := strconv.Atoi(permissionId)
			permissions := findPermissionById(DB, permissionIdint)
			if permissions.ID == 0 {
				response.Fail(c, permissionId+"不存在权限", nil)
				return
			}
		}
	}

	tx := DB.Begin()

	data := make(map[string]interface{})
	data["name"] = name
	data["slug"] = slug

	if err := tx.Model(&model.AdminRoles{}).Where("id = ?", id).Updates(&data).Error; err != nil {
		response.Fail(c, "更新失败", nil)
		tx.Rollback()
		return
	}

	//先批量删除角色对应权限
	err := tx.Where("role_id = ?", Id).Unscoped().Delete(&model.AdminRolePermissions{}).Error
	if err != nil {
		response.Fail(c, "删除失败", nil)
		tx.Rollback()
		return
	}

	if len(permissionIdArr) > 0 {
		//批量添加
		var arr []model.AdminRolePermissions
		idint, _ := strconv.Atoi(id)
		arr = make([]model.AdminRolePermissions, len(permissionIdArr))
		for k, permissionId := range permissionIdArr {
			permissionIdint, _ := strconv.Atoi(permissionId)
			arr[k].RoleId = idint
			arr[k].PermissionId = permissionIdint
		}
		err = model.AddRolePermissions(tx, arr)
		if err != nil {
			response.Fail(c, "更新权限失败", nil)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	response.Success(c, nil, "成功")
}

func RolesAllinfo(ctx *gin.Context) {
	db := common.GetDB()
	var roles []model.AdminRoles
	var _roles []model.AdminRolesTodo
	db.Find(&roles)
	// 格式化
	for _, item := range roles {
		_roles = append(_roles, model.ToAdminRolesTodo(item))
	}
	response.Success(ctx, _roles, "成功")
}

func findRolesById(db *gorm.DB, id int) model.AdminRoles {
	var roles model.AdminRoles
	db.Where("id = ?", id).Find(&roles)
	return roles
}

func isHaveRolesSlug(db *gorm.DB, slug string) bool {
	var roles model.AdminRoles
	err := db.Where("slug = ?", slug).Find(&roles).Error
	if err == nil && roles.ID != 0 {
		return true
	}
	return false
}
