package common

import (
	"database/sql"
	"fmt"
	"go-gin-admin/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	host := "47.107.140.71"
	port := "3306"
	username := "aliyunget"
	password := "yuan@get940314"
	database := "test"
	charset := "utf8"

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
	DB = db
	return db
}

func GetDB() *gorm.DB  {
	return DB
}
