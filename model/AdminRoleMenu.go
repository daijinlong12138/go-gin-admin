package model

import (
	"gorm.io/gorm"
)

type AdminRoleMenu struct {
	gorm.Model
	RoleId int `gorm:"type:bigint(20);not null;index:role_id_menu_id;unsigned"`
	MenuId int `gorm:"type:bigint(20);not null;index:role_id_menu_id;unsigned"`
}
