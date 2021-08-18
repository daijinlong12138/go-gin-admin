package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
)

type AdminRoleUsers struct {
	gorm.Model
	RoleId int `gorm:"type:bigint(20) unsigned;not null;uniqueIndex:admin_user_roles"`
	UserId int `gorm:"type:bigint(20) unsigned;not null;uniqueIndex:admin_user_roles"`
}

func AddRoleUsers(db *gorm.DB, RoleUsers []AdminRoleUsers) error {
	return BatchSaveAdminRoleUsers(db, RoleUsers)
}

// 批量插入数据
func BatchSaveAdminRoleUsers(db *gorm.DB, RoleUsers []AdminRoleUsers) error {
	var buffer bytes.Buffer
	sql := "insert into `admin_role_users` (`role_id`,`user_id`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, RoleUser := range RoleUsers {
		if i == len(RoleUsers)-1 {
			buffer.WriteString(fmt.Sprintf("('%d','%d');", RoleUser.RoleId, RoleUser.UserId))
		} else {
			buffer.WriteString(fmt.Sprintf("('%d','%d'),", RoleUser.RoleId, RoleUser.UserId))
		}
	}
	return db.Exec(buffer.String()).Error
}
