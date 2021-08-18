package model

import (
	"gorm.io/gorm"
)

type AdminMenu struct {
	gorm.Model
	Id         int    `gorm:"-;primary_key;AUTO_INCREMENT"`
	ParentId   int    `gorm:"type:bigint(20) unsigned;not null;default:0;comment:父级ID"`
	Type       int    `gorm:"type:tinyint(4) unsigned;not null;default:0;comment:类型"`
	Order      int    `gorm:"type:int(11) unsigned;not null;default:0;comment:order"`
	Title      string `gorm:"type:varchar(50);not null;comment:菜单名"`
	Icon       string `gorm:"type:varchar(50);not null;comment:图标"`
	Uri        string `gorm:"type:varchar(3000);not null;default:'';comment:路径"`
	Header     string `gorm:"type:varchar(150);comment:header"`
	PluginName string `gorm:"type:varchar(150);not null;default:'';comment:插件名称"`
	Uuid       string `gorm:"type:varchar(150)"`
}
