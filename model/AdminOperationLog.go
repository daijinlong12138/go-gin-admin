package model

import (
	"gorm.io/gorm"
)

type AdminOperationLog struct {
	gorm.Model
	Id     int    `gorm:"-;primary_key;AUTO_INCREMENT"`
	UserId int    `gorm:"type:bigint(20) unsigned;not null;index;comment:管理员ID"`
	Path   string `gorm:"type:varchar(255);not null;comment:路径"`
	Method string `gorm:"type:varchar(10);not null;comment:方法"`
	Ip     string `gorm:"type:varchar(15);not null;comment:ip"`
	Input  string `gorm:"type:text;not null;comment:内容"`
}
