package contronller

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"go-gin-admin/dto"
	"go-gin-admin/model"
	"go-gin-admin/response"
	"go-gin-admin/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := c.PostForm("name")
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	//校验
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码必须大于6位")
		return
	}
	if len(name) == 0 {
		name = util.RandGetName(10)
	}

	if isTelephoneExist(DB, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已经存在")
		return
	}

	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "加密错误")
		return
	}

	newUser := model.Users{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	log.Println(name, telephone, password)

	response.Success(c, nil, "注册成功")
}

func Login(c *gin.Context) {

	DB := common.GetDB()

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

	var user model.Users
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Fail(c, "用户不存在", nil)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(c, "密码错误", nil)
		return
	}

	/*token,err := common.ReleaseToken(user)
	if err!=nil{
		response.Fail(c,"系统异常",nil)
		log.Printf("token generate error : %v",err)
		return
	}

	response.Success(c,gin.H{"token":token},"登录成功")*/
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	response.Success(c, dto.ToUserDto(user.(model.Users)), "成功")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var users model.Users
	db.Where("telephone = ?", telephone).First(&users)
	if users.ID != 0 {
		return true

	}
	return false
}
