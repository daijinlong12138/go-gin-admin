package contronller

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"go-gin-admin/model"
	"go-gin-admin/util"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(c *gin.Context)  {
	DB := common.GetDB()
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
		name = util.RandGetName(10)
	}

	if isTelephoneExist(DB, telephone) {
		c.JSON(http.StatusOK, gin.H{"code": 422, "msg": "用户已经存在"})
		return
	}

	//创建用户
	newUser := model.Users{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)

	log.Println(name, telephone, password)

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "注册成功"})
}



func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var users model.Users
	db.Where("telephone = ?", telephone).First(&users)
	if users.ID != 0 {
		return true

	}
	return false
}