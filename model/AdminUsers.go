package model

import (
	"gorm.io/gorm"
)

type AdminUsers struct {
	gorm.Model
	Id            int                `gorm:"-;primary_key;AUTO_INCREMENT"`
	Username      string             `gorm:"type:varchar(100);not null;unique;comment:用户名(用于登录)"`
	Password      string             `gorm:"type:varchar(100);not null;default:'';comment:密码"`
	Name          string             `gorm:"type:varchar(100);not null;comment:昵称(用来展示)"`
	Avatar        string             `gorm:"type:varchar(255);comment:头像"`
	RememberToken string             `gorm:"type:varchar(100);"`
	Roles         []AdminRoles       `gorm:"-"`
	Permissions   []AdminPermissions `gorm:"-"`
}

type AdminUsersTodo struct {
	Id          int                    `json:"id"`
	Username    string                 `json:"username"`
	Name        string                 `json:"name"`
	Avatar      string                 `json:"avatar"`
	Roles       []AdminRolesTodo       `json:"roles"`
	Permissions []AdminPermissionsTodo `json:"permissions"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

func ToAdminUsersTodo(Users AdminUsers) AdminUsersTodo {

	var AdminRolesTodos []AdminRolesTodo
	var AdminPermissionsTodos []AdminPermissionsTodo

	for i, _ := range Users.Roles {
		AdminRolesTodos = append(AdminRolesTodos, ToAdminRolesTodo(Users.Roles[i]))
	}
	for i, _ := range Users.Permissions {
		AdminPermissionsTodos = append(AdminPermissionsTodos, ToAdminPermissionsTodo(Users.Permissions[i]))
	}

	return AdminUsersTodo{
		Id:          int(Users.ID),
		Username:    Users.Username,
		Name:        Users.Name,
		Avatar:      Users.Avatar,
		Roles:       AdminRolesTodos,
		Permissions: AdminPermissionsTodos,
		CreatedAt:   Users.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   Users.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
