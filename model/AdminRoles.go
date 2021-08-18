package model

import (
	"gorm.io/gorm"
)

type AdminRoles struct {
	gorm.Model
	Id   int    `gorm:"-;primary_key;AUTO_INCREMENT"`
	Name string `gorm:"type:varchar(50);not null;comment:角色名称"`
	Slug string `gorm:"type:varchar(50);not null;unique;comment:标志(需要保证唯一)"`
}

type AdminRolesTodo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ToAdminRolesTodo(roles AdminRoles) AdminRolesTodo {
	return AdminRolesTodo{
		Id:        int(roles.ID),
		Name:      roles.Name,
		Slug:      roles.Slug,
		CreatedAt: roles.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: roles.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
