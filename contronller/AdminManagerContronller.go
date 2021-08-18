package contronller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-admin/common"
	"go-gin-admin/model"
	"go-gin-admin/response"
	"go-gin-admin/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ManagerNew(c *gin.Context) {
	UploadImg := viper.GetString("Upload.img")
	UploadTmp := viper.GetString("Upload.tmp")
	DB := common.GetDB()
	//获取参数
	username := c.PostForm("username")             //用户名(用于登录)
	password := c.PostForm("password")             //密码
	name := c.PostForm("name")                     //昵称(用来展示)
	avatar := c.PostForm("avatar")                 //头像
	role_ids := c.PostForm("role_ids")             //角色,逗号分割
	permission_ids := c.PostForm("permission_ids") //权限,逗号分割

	//校验
	if len(name) == 0 {
		response.Fail(c, "昵称不能为空", nil)
		return
	}
	if len(username) == 0 {
		response.Fail(c, "用户名不能为空", nil)
		return
	}
	if len(password) < 6 || len(password) > 20 {
		response.Fail(c, "密码必须大于6位且小于二十位", nil)
		return
	}

	var roleIdArr []string
	if len(permission_ids) != 0 {
		roleIdArr = strings.Split(role_ids, `,`)
	}
	var permissionIdArr []string
	if len(permission_ids) != 0 {
		permissionIdArr = strings.Split(permission_ids, `,`)
	}
	var fildDir string
	if len(avatar) > 0 {
		IsExist, _ := util.IsFileExist(UploadTmp + avatar)
		if !IsExist {
			response.Fail(c, "头像文件不存在", nil)
			return
		}
		fildDir = fmt.Sprintf("%d%s/", time.Now().Year(), time.Now().Month().String())
		isExist, _ := util.IsFileExist(UploadImg + fildDir)
		if !isExist {
			if err := os.Mkdir(UploadImg+fildDir, os.ModePerm); err != nil {
				response.Fail(c, "创建文件夹失败", nil)
				return
			}
		}
	}

	if isHaveManagerUsername(DB, username) {
		response.Fail(c, "用户名已存在", nil)
		return
	}

	if len(roleIdArr) > 0 {
		for _, roleId := range roleIdArr {
			roleIdint, _ := strconv.Atoi(roleId)
			roles := findRolesById(DB, roleIdint)
			if roles.ID == 0 {
				response.Fail(c, roleId+"不存在角色", nil)
				return
			}
		}
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

	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, "加密错误", nil)
		return
	}

	tx := DB.Begin()
	//创建用户

	newUser := model.AdminUsers{
		Username: username,
		Name:     name,
		Password: string(hasedPassword),
		Avatar:   fildDir + avatar,
	}
	err = tx.Create(&newUser).Error
	if err != nil || newUser.ID == 0 {
		response.Fail(c, "创建失败", nil)
		tx.Rollback()
		return
	}

	if len(roleIdArr) > 0 {
		//批量添加
		var arr []model.AdminRoleUsers
		arr = make([]model.AdminRoleUsers, len(roleIdArr))
		for k, roleId := range roleIdArr {
			roleIdint, _ := strconv.Atoi(roleId)
			arr[k].RoleId = roleIdint
			arr[k].UserId = int(newUser.ID)
		}
		err = model.AddRoleUsers(tx, arr)
		if err != nil {
			response.Fail(c, "用户添加角色失败", nil)
			tx.Rollback()
			return
		}
	}

	if len(permissionIdArr) > 0 {
		//批量添加
		var arr []model.AdminUserPermissions
		arr = make([]model.AdminUserPermissions, len(permissionIdArr))
		for k, permissionId := range permissionIdArr {
			permissionIdint, _ := strconv.Atoi(permissionId)
			arr[k].PermissionId = permissionIdint
			arr[k].UserId = int(newUser.ID)
		}
		err = model.AddUserPermissions(tx, arr)
		if err != nil {
			response.Fail(c, "用户添加权限失败", nil)
			tx.Rollback()
			return
		}
	}

	if len(avatar) > 0 {
		// 移动文件
		if err := os.Rename(UploadTmp+avatar, UploadImg+fildDir+avatar); err != nil {
			response.Fail(c, "头像文件移动保存失败", nil)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	response.Success(c, nil, "用户创建成功")
}

func ManagerDetail(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id") //ID
	var user model.AdminUsers

	DB.Where("id = ?", id).Find(&user)
	if user.ID == 0 {
		response.Fail(c, "不存在", nil)
		return
	}
	FindUserDetailInfoById(DB, &user)
	response.Success(c, model.ToAdminUsersTodo(user), "成功")
}

func ManagerDelete(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id") //角色ID
	var user model.AdminUsers
	Id, _ := strconv.Atoi(id)
	DB.Where("id = ?", id).Find(&user)
	if user.ID == 0 {
		response.Fail(c, "不存在", nil)
		return
	}

	tx := DB.Begin()
	//用户表
	err := tx.Where("id = ?", Id).Unscoped().Delete(&model.AdminUsers{}).Error
	if err != nil {
		tx.Rollback()
		response.Fail(c, "删除失败", nil)
		return
	}
	//关联表
	err = tx.Where("user_id = ?", Id).Unscoped().Delete(&model.AdminUserPermissions{}).Error
	if err != nil {
		tx.Rollback()
		response.Fail(c, "删除用户权限关联表失败", nil)
		return
	}

	err = tx.Where("user_id = ?", Id).Unscoped().Delete(&model.AdminRoleUsers{}).Error
	if err != nil {
		tx.Rollback()
		response.Fail(c, "删除用户角色关联表失败", nil)
		return
	}

	tx.Commit()
	response.Success(c, nil, "用户删除成功")
}

func ManagerInfo(c *gin.Context) {
	var user model.AdminUsers
	var total int64
	users := make([]model.AdminUsers, 0)
	var _users []model.AdminUsersTodo

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

	if rolesname, isExist := c.GetPostForm("rolesname"); isExist == true && rolesname != "" {
		//以角色为主查询用户
		var role model.AdminRoles
		db.Where("name = ?", rolesname).Find(&role)
		/*if role.ID==0{
			response.Fail(c,"角色名称不存在",nil)
			return
		}*/
		db = db.Table("admin_role_users").
			Select("admin_users.id as id, admin_users.name as name, admin_users.username as username,admin_users.avatar as avatar,admin_users.created_at as created_at,admin_users.updated_at as updated_at ").
			Joins("left join admin_users on admin_users.id = admin_role_users.user_id").
			Where("admin_role_users.id = ?", role.ID)
	} else {
		//以用户查询
		db = db.Model(user) //查询对应的数据库表

		if username, isExist := c.GetPostForm("username"); isExist == true && username != "" {
			db = db.Where("username = ?", username)
		}

		if err := db.Count(&total).Error; err != nil {
			response.Fail(c, "查询数据异常", nil)
			return
		}

	}

	db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	if err := db.Find(&users).Error; err != nil {
		fmt.Println(err.Error())
		response.Fail(c, "查询数据异常", nil)
		return
	}

	// 格式化
	if len(users) > 0 {
		for _, item := range users {
			//处理角色名称
			FindUserDetailInfoById(db, &item)
			//fmt.Println(item)
			_users = append(_users, model.ToAdminUsersTodo(item))
		}
	}

	//fmt.Println(roles)
	response.Success(c, gin.H{
		"data":     _users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, "成功")
}

func ManagerEdit(c *gin.Context) {
	UploadImg := viper.GetString("Upload.img")
	UploadTmp := viper.GetString("Upload.tmp")
	DB := common.GetDB()
	var user = model.AdminUsers{}
	//获取参数
	id := c.PostForm("id")                         //用户id
	name := c.PostForm("name")                     //昵称(用来展示)
	username := c.PostForm("username")             //用户名(用于登录)
	avatar := c.PostForm("avatar")                 //头像
	password := c.PostForm("password")             //密码
	permission_ids := c.PostForm("permission_ids") //权限以逗号分割
	roles_ids := c.PostForm("roles_ids")           //角色以逗号分割

	if len(name) == 0 {
		response.Fail(c, "昵称不能为空", nil)
		return
	}
	if len(username) == 0 {
		response.Fail(c, "用户名不能为空", nil)
		return
	}

	var permissionIdArr []string
	if len(permission_ids) != 0 {
		permissionIdArr = strings.Split(permission_ids, `,`)
	}
	var rolesIdArr []string
	if len(roles_ids) != 0 {
		rolesIdArr = strings.Split(roles_ids, `,`)
	}
	Id, _ := strconv.Atoi(id)
	user = findUserById(DB, Id)
	if user.ID != uint(Id) {
		response.Fail(c, "不存在用户", nil)
		return
	}

	fildDir := ""
	if user.Avatar != avatar && len(avatar) > 0 {
		IsExist, _ := util.IsFileExist(UploadTmp + avatar)
		if !IsExist {
			response.Fail(c, "头像文件不存在", nil)
			return
		}
		fildDir = fmt.Sprintf("%d%s/", time.Now().Year(), time.Now().Month().String())
		isExist, _ := util.IsFileExist(UploadImg + fildDir)
		if !isExist {
			if err := os.Mkdir(UploadImg+fildDir, os.ModePerm); err != nil {
				response.Fail(c, "创建文件夹失败", nil)
				return
			}
		}
	}

	passwordString := ""
	if len(password) > 0 {
		if len(password) < 6 || len(password) > 20 {
			response.Fail(c, "密码必须大于6位且小于二十位", nil)
			return
		}
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.Fail(c, "加密错误", nil)
			return
		}

		passwordString = string(hasedPassword)
	}

	if user.Username != username && isHaveManagerUsername(DB, username) {
		response.Fail(c, "用户名已存在", nil)
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
	if len(rolesIdArr) > 0 {
		for _, rolesId := range rolesIdArr {
			rolesIdint, _ := strconv.Atoi(rolesId)
			role := findRolesById(DB, rolesIdint)
			if role.ID == 0 {
				response.Fail(c, rolesId+"不存在权限", nil)
				return
			}
		}
	}

	tx := DB.Begin()

	data := make(map[string]interface{})
	data["username"] = username
	if passwordString != "" {
		data["password"] = passwordString
	}
	data["name"] = name
	data["avatar"] = fildDir + avatar

	if err := tx.Model(&model.AdminUsers{}).Where("id = ?", id).Updates(&data).Error; err != nil {
		response.Fail(c, "更新失败", nil)
		tx.Rollback()
		return
	}

	//批量删除用户对应权限
	err := tx.Where("user_id = ?", Id).Unscoped().Delete(&model.AdminUserPermissions{}).Error
	if err != nil {
		response.Fail(c, "删除权限失败", nil)
		tx.Rollback()
		return
	}

	if len(permissionIdArr) > 0 {
		//批量添加
		var arr []model.AdminUserPermissions
		arr = make([]model.AdminUserPermissions, len(permissionIdArr))
		for k, permissionId := range permissionIdArr {
			permissionIdint, _ := strconv.Atoi(permissionId)
			arr[k].PermissionId = permissionIdint
			UserIdint, _ := strconv.Atoi(id)
			arr[k].UserId = UserIdint
		}
		err = model.AddUserPermissions(tx, arr)
		if err != nil {
			response.Fail(c, "更新权限失败", nil)
			tx.Rollback()
			return
		}
	}

	//批量删除用户对应角色
	err = tx.Where("user_id = ?", Id).Unscoped().Delete(&model.AdminRoleUsers{}).Error
	if err != nil {
		response.Fail(c, "删除角色失败", nil)
		tx.Rollback()
		return
	}
	if len(rolesIdArr) > 0 {
		//批量添加
		var arr []model.AdminRoleUsers
		arr = make([]model.AdminRoleUsers, len(rolesIdArr))
		for k, rolesId := range rolesIdArr {
			rolesIdint, _ := strconv.Atoi(rolesId)
			arr[k].RoleId = rolesIdint
			UserIdint, _ := strconv.Atoi(id)
			arr[k].UserId = UserIdint
		}
		err = model.AddRoleUsers(tx, arr)
		if err != nil {
			response.Fail(c, "更新角色失败", nil)
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	response.Success(c, nil, "成功")
}

func ManagerLogin(c *gin.Context) {

	DB := common.GetDB()
	var user = model.AdminUsers{}
	username := c.PostForm("username")
	password := c.PostForm("password")

	//校验
	if len(username) == 0 {
		response.Fail(c, "用户名不能为空", nil)
		return
	}
	if len(password) < 6 {
		response.Fail(c, "密码必须大于6位", nil)
		return
	}

	DB.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		response.Fail(c, "用户不存在", nil)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(c, "密码错误", nil)
		return
	}

	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Fail(c, "系统异常", nil)
		log.Printf("token generate error : %v", err)
		return
	}
	response.Success(c, gin.H{"token": token}, "登录成功")
}

func findUserById(db *gorm.DB, id int) model.AdminUsers {
	var user model.AdminUsers
	db.Where("id = ?", id).Find(&user)
	return user
}

func isHaveManagerUsername(db *gorm.DB, username string) bool {
	var user model.AdminUsers
	err := db.Where("username = ?", username).Find(&user).Error
	if err == nil && user.ID != 0 {
		return true
	}
	return false
}

func FindUserDetailInfoById(db *gorm.DB, user *model.AdminUsers) {
	//角色
	db.Raw("select admin_roles.id as id, admin_roles.name as name, admin_roles.slug as slug, admin_roles.created_at as created_at,admin_roles.updated_at as updated_at "+
		"from admin_role_users "+
		"left join admin_roles on admin_roles.id = admin_role_users.role_id where admin_role_users.user_id = ?", user.ID).Scan(&user.Roles)
	//权限
	db.Raw("SELECT admin_permissions.id AS id, admin_permissions. NAME AS NAME, admin_permissions.slug AS slug, admin_permissions.http_method AS http_method, admin_permissions.http_path AS http_path, admin_permissions.created_at AS created_at, admin_permissions.updated_at AS updated_at "+
		"FROM admin_permissions "+
		"LEFT JOIN admin_role_permissions ON admin_role_permissions.permission_id = admin_permissions.id "+
		"LEFT JOIN admin_role_users ON admin_role_permissions.role_id = admin_role_users.role_id "+
		"LEFT JOIN admin_user_permissions ON admin_permissions.id = admin_user_permissions.permission_id "+
		"LEFT JOIN admin_users ON admin_users.id = admin_role_users.user_id OR admin_users.id = admin_user_permissions.user_id "+
		"WHERE admin_users.id = ? GROUP BY admin_permissions.id", user.ID).Scan(&user.Permissions)
	return
}
