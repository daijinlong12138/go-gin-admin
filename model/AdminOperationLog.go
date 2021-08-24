package model

import (
	"gorm.io/gorm"
)

type AdminOperationLog struct {
	gorm.Model
	Id       int    `gorm:"-;primary_key;AUTO_INCREMENT"`
	UserId   int    `gorm:"type:bigint(20) unsigned;not null;index;comment:管理员ID"`
	UserName string `gorm:"type:varchar(255);not null;comment:管理员名称"`
	Path     string `gorm:"type:varchar(255);not null;comment:路径"`
	Method   string `gorm:"type:varchar(10);not null;comment:方法"`
	Ip       string `gorm:"type:varchar(15);not null;comment:ip"`
	Input    string `gorm:"type:text;not null;comment:内容"`
}

type AdminOperationLogTodo struct {
	Id        int    `json:"id"`
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Ip        string `json:"ip"`
	Input     string `json:"input"`
	CreatedAt string `json:"created_at"`
}

func ToAdminOperationLogTodo(opLogs AdminOperationLog) AdminOperationLogTodo {
	return AdminOperationLogTodo{
		Id:        int(opLogs.ID),
		UserId:    opLogs.UserId,
		UserName:  opLogs.UserName,
		Path:      opLogs.Path,
		Method:    opLogs.Method,
		Ip:        opLogs.Ip,
		Input:     opLogs.Input,
		CreatedAt: opLogs.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
