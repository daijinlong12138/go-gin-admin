package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
)

type AdminRoleMenu struct {
	gorm.Model
	RoleId int `gorm:"type:bigint(20);not null;index:role_id_menu_id;unsigned"`
	MenuId int `gorm:"type:bigint(20);not null;index:role_id_menu_id;unsigned"`
}

func AddRoleMenus(db *gorm.DB, RoleMenus []AdminRoleMenu) error {
	return BatchSaveAdminRoleMenus(db, RoleMenus)
}

// 批量插入数据
func BatchSaveAdminRoleMenus(db *gorm.DB, RoleMenus []AdminRoleMenu) error {
	var buffer bytes.Buffer
	sql := "insert into `admin_role_menus` (`role_id`,`menu_id`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, RoleMenu := range RoleMenus {
		if i == len(RoleMenus)-1 {
			buffer.WriteString(fmt.Sprintf("('%d','%d');", RoleMenu.RoleId, RoleMenu.MenuId))
		} else {
			buffer.WriteString(fmt.Sprintf("('%d','%d'),", RoleMenu.RoleId, RoleMenu.MenuId))
		}
	}
	return db.Exec(buffer.String()).Error
}
