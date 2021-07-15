package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Users struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(20);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {

	db := InitDB()

	// 1.创建路由
	router := gin.Default()

	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response

	router.POST("/api/auth/register", func(c *gin.Context) {
		//获取参数
		name := c.PostForm("name")
		telephone := c.PostForm("telephone")
		password := c.PostForm("password")

		//校验
		if len(telephone) != 11 {
			c.JSON(http.StatusOK, gin.H{"code": 422, "msg": "手机必须为11位"})
			return
		}
		if len(password) < 6 {
			c.JSON(http.StatusOK, gin.H{"code": 422, "msg": "密码必须大于6位"})
			return
		}
		if len(name) == 0 {
			name = randGetName(10)
		}

		if isTelephoneExist(db, telephone) {
			c.JSON(http.StatusOK, gin.H{"code": 422, "msg": "用户已经存在"})
			return
		}

		//创建用户
		newUser := Users{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		log.Println(name, telephone, password)

		c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})

	})

	// 3.监听端口，默认在8080
	// Run("里面不指定端口号默认为8080")
	fmt.Printf("start server (port:%s)", "8000")
	router.Run(":8000")

}

func randGetName(n int) string {
	var letters = []byte("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

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

	db.AutoMigrate(&Users{})

	return db
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var users Users
	db.Where("telephone = ?", telephone).First(&users)
	if users.ID != 0 {
		return true

	}
	return false
}
