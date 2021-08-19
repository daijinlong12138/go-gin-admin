package common

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"go-gin-admin/model"
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
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Users{})

	db.AutoMigrate(&model.AdminRoles{})
	db.AutoMigrate(&model.AdminUsers{})
	db.AutoMigrate(&model.AdminRoleUsers{})
	db.AutoMigrate(&model.AdminPermissions{})
	db.AutoMigrate(&model.AdminMenu{})
	db.AutoMigrate(&model.AdminRoleMenu{})
	db.AutoMigrate(&model.AdminRolePermissions{})
	db.AutoMigrate(&model.AdminUserPermissions{})
	db.AutoMigrate(&model.AdminOperationLog{})

	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
