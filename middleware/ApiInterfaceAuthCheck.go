package middleware

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"go-gin-admin/contronller/auth"
	"go-gin-admin/model"
	"go-gin-admin/response"
	"strings"
)

//操作日志数据库缺少参数数据
func ApiInterfaceAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		//获取authorization header
		tokenString := c.GetHeader("Authorization")
		method := c.Request.Method
		url := c.Request.URL.Path
		if url != "/admin/manager/login" {
			//validate token formate
			if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
				response.AuthFail(c, nil, "Authorization None")
				c.Abort()
				return
			}

			tokenString = tokenString[7:]
			token, claims, err := common.ParseToken(tokenString)
			if err != nil || !token.Valid {
				response.AuthFail(c, nil, "Authorization Error")
				c.Abort()
				return
			}

			//验证通过后获取 claims的userId
			userId := claims.UserID
			DB := common.GetDB()
			var user model.AdminUsers
			err = DB.First(&user, userId).Error
			if err != nil {
				response.AuthFail(c, nil, "账户 Fail")
				c.Abort()
				return
			}
			if user.ID == 0 {
				response.AuthFail(c, nil, "账户不存在")
				c.Abort()
				return
			}
			auth.FindUserDetailInfoById(DB, &user)

			if userId != 1 && !findUserPermission(method, url, user) {
				//默认id为1的用户为超级用户，不需要权限检验
				response.AuthFail(c, nil, "权限不足")
				c.Abort()
				return
			}
			//记录到数据库中
			IP := c.ClientIP()
			newAdminOperationLog := model.AdminOperationLog{
				UserId:   int(userId),
				UserName: user.Name,
				Method:   method,
				Path:     url,
				Ip:       IP,
			}
			auth.InsertLogs(DB, newAdminOperationLog)

			c.Set("user", user)
			c.Set("userId", userId)
		}

		c.Next()
	}
}

func findUserPermission(method string, url string, user model.AdminUsers) bool {
	var needPermissions []model.AdminPermissions
	if len(user.Permissions) < 1 {
		return false
	}
	if len(url) < 1 {
		return false
	}
	//1.循环查找，满足当前url路径permissions有多少
	for _, item := range user.Permissions {
		if strings.Contains(item.HttpPath, url) {
			needPermissions = append(needPermissions, item)
		}
	}
	if len(needPermissions) == 0 {
		return false
	}
	//2.在这些permissions查找是否满足方法method
	var isPass = false
	for _, item := range needPermissions {
		if len(item.HttpMethod) == 0 {
			//空表示全部方法
			isPass = true
			break
		}
		if strings.Contains(item.HttpMethod, method) {
			isPass = true
			break
		}
	}
	return isPass
}
