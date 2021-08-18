package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
)

type AdminRolePermissions struct {
	gorm.Model
	RoleId       int `gorm:"type:bigint(20) unsigned;not null;uniqueIndex:admin_role_permissions"`
	PermissionId int `gorm:"type:bigint(20) unsigned;not null;uniqueIndex:admin_role_permissions"`
}

func AddRolePermissions(db *gorm.DB, RolePermissions []AdminRolePermissions) error {
	return BatchSaveAdminRolePermissions(db, RolePermissions)
}

// 批量插入数据
func BatchSaveAdminRolePermissions(db *gorm.DB, RolePermissions []AdminRolePermissions) error {
	var buffer bytes.Buffer
	sql := "insert into `admin_role_permissions` (`role_id`,`permission_id`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, RolePermission := range RolePermissions {
		if i == len(RolePermissions)-1 {
			buffer.WriteString(fmt.Sprintf("('%d','%d');", RolePermission.RoleId, RolePermission.PermissionId))
		} else {
			buffer.WriteString(fmt.Sprintf("('%d','%d'),", RolePermission.RoleId, RolePermission.PermissionId))
		}
	}
	return db.Exec(buffer.String()).Error
}
