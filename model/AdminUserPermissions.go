package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
)

type AdminUserPermissions struct {
	gorm.Model
	UserId       int `gorm:"type:bigint(20) unsigned;not null;uniqueIndex:admin_user_permissions"`
	PermissionId int `gorm:"type:bigint(20) unsigned;not null;uniqueIndex:admin_user_permissions"`
}

func AddUserPermissions(db *gorm.DB, UserPermissions []AdminUserPermissions) error {
	return BatchSaveAdminUserPermissions(db, UserPermissions)
}

// 批量插入数据--后期可优化
func BatchSaveAdminUserPermissions(db *gorm.DB, UserPermissions []AdminUserPermissions) error {
	var buffer bytes.Buffer
	sql := "insert into `admin_user_permissions` (`user_id`,`permission_id`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, UserPermission := range UserPermissions {
		if i == len(UserPermissions)-1 {
			buffer.WriteString(fmt.Sprintf("('%d','%d');", UserPermission.UserId, UserPermission.PermissionId))
		} else {
			buffer.WriteString(fmt.Sprintf("('%d','%d'),", UserPermission.UserId, UserPermission.PermissionId))
		}
	}
	return db.Exec(buffer.String()).Error
}
