package auth

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"go-gin-admin/model"
	"go-gin-admin/response"
	"gorm.io/gorm"
	"strconv"
)

func OpLogsInfo(c *gin.Context) {
	var OpLog model.AdminOperationLog
	var total int64
	OpLogs := make([]model.AdminOperationLog, 0)
	var _OpLog []model.AdminOperationLogTodo

	db := common.GetDB()
	//获取参数
	page, _ := strconv.Atoi(c.DefaultPostForm("page", "1"))          //页数
	pageSize, _ := strconv.Atoi(c.DefaultPostForm("pageSize", "10")) //每页数

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	//查询
	db = db.Model(OpLog) //查询对应的数据库表

	if username, isExist := c.GetPostForm("user_name"); isExist == true && username != "" {
		db = db.Where("user_name = ?", username)
	}
	if path, isExist := c.GetPostForm("path"); isExist == true && path != "" {
		db = db.Where("path = ?", path)
	}
	if method, isExist := c.GetPostForm("method"); isExist == true && method != "" {
		db = db.Where("method = ?", method)
	}
	if ip, isExist := c.GetPostForm("ip"); isExist == true && ip != "" {
		db = db.Where("ip = ?", ip)
	}

	if err := db.Count(&total).Error; err != nil {
		common.LogError(c, "查询数据异常 : "+err.Error())
		response.Fail(c, "查询数据异常", nil)
		return
	}

	db = db.Limit(pageSize).Offset((page - 1) * pageSize)
	if err := db.Find(&OpLogs).Error; err != nil {
		common.LogError(c, "查询数据异常 : "+err.Error())
		response.Fail(c, "查询数据异常", nil)
		return
	}

	// 格式化
	if len(OpLogs) > 0 {
		for _, item := range OpLogs {
			//fmt.Println(item)
			_OpLog = append(_OpLog, model.ToAdminOperationLogTodo(item))
		}
	}

	//fmt.Println(roles)
	response.Success(c, gin.H{
		"data":     _OpLog,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}, "成功")
}

func InsertLogs(db *gorm.DB, operationLog model.AdminOperationLog) bool {

	err := db.Create(&operationLog).Error
	if err != nil || operationLog.ID == 0 {
		return false
	}
	return true
}
