package model

import (
	"gorm.io/gorm"
	"strings"
)

type AdminPermissions struct {
	gorm.Model
	Id         int    `gorm:"-;primary_key;AUTO_INCREMENT"`
	Name       string `gorm:"type:varchar(50);not null;comment:权限名称"`
	Slug       string `gorm:"type:varchar(50);not null;unique;comment:标志(需要保证唯一)"`
	HttpMethod string `gorm:"type:varchar(255);comment:方法(为空默认为所有方法)"`
	HttpPath   string `gorm:"type:text;not null;comment:路径(路径不包含全局路由前缀，换行多条)"`
}

type AdminPermissionsTodo struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	HttpMethod []string `json:"http_method"`
	HttpPath   []string `json:"http_path"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
}

func ToAdminPermissionsTodo(permissions AdminPermissions) AdminPermissionsTodo {

	var HttpPath []string
	var HttpMethod []string
	if len(permissions.HttpMethod) != 0 {
		HttpMethod = strings.Split(permissions.HttpMethod, ",")
	}
	if len(permissions.HttpPath) != 0 {
		HttpPath = strings.Split(permissions.HttpPath, "\n")
	}
	return AdminPermissionsTodo{
		Id:         int(permissions.ID),
		Name:       permissions.Name,
		Slug:       permissions.Slug,
		HttpMethod: HttpMethod,
		HttpPath:   HttpPath,
		CreatedAt:  permissions.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  permissions.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
