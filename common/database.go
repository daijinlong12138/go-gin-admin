package common

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"go-gin-admin/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {

	host := viper.GetString("datasoure.host")
	port := viper.GetString("datasoure.port")
	username := viper.GetString("datasoure.username")
	password := viper.GetString("datasoure.password")
	database := viper.GetString("datasoure.database")
	charset := viper.GetString("datasoure.charset")

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	sqlDB, err := sql.Open("mysql", args)
	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		Log.Error("failed to connect database " + err.Error())
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&model.AdminRoles{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminRoles " + err.Error())
		panic("failed to AutoMigrate AdminRoles")
	}
	err = db.AutoMigrate(&model.AdminUsers{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminUsers " + err.Error())
		panic("failed to AutoMigrate AdminUsers")
	}
	err = db.AutoMigrate(&model.AdminRoleUsers{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminRoleUsers " + err.Error())
		panic("failed to AutoMigrate AdminRoleUsers")
	}
	err = db.AutoMigrate(&model.AdminPermissions{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminPermissions " + err.Error())
		panic("failed to AutoMigrate AdminPermissions")
	}
	err = db.AutoMigrate(&model.AdminMenu{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminMenu " + err.Error())
		panic("failed to AutoMigrate AdminMenu")
	}
	err = db.AutoMigrate(&model.AdminRoleMenu{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminRoleMenu " + err.Error())
		panic("failed to AutoMigrate AdminRoleMenu")
	}
	err = db.AutoMigrate(&model.AdminRolePermissions{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminRolePermissions " + err.Error())
		panic("failed to AutoMigrate AdminRolePermissions")
	}
	err = db.AutoMigrate(&model.AdminUserPermissions{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminUserPermissions " + err.Error())
		panic("failed to AutoMigrate AdminUserPermissions")
	}
	err = db.AutoMigrate(&model.AdminOperationLog{})
	if err != nil {
		Log.Error("failed to AutoMigrate AdminOperationLog " + err.Error())
		panic("failed to AutoMigrate AdminOperationLog")
	}

	//检查用户是否有id=1的，没有表插入默认数据  123456
	var user model.AdminUsers
	err = db.Where("id = 1").Find(&user).Error
	if err != nil {
		Log.Error("查询失败: " + err.Error())
		panic(" 查询失败 ")
	}
	if user.ID == 0 {
		hasedPassword, err := bcrypt.GenerateFromPassword([]byte("123qwm@"), bcrypt.DefaultCost)
		if err != nil {
			Log.Error("加密错误: " + err.Error())
			panic(" 加密错误 ")
		}
		//创建用户
		err = db.Exec("INSERT INTO `admin_users` (`id`, `created_at`, `username`, `password`, `name`) VALUES "+
			"(1, NOW(),  'admin', ?, 'admin')", string(hasedPassword)).Error
		if err != nil {
			Log.Error("创建失败: " + err.Error())
			panic(" 创建失败 ")
		}
	}

	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
