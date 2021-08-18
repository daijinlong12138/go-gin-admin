package contronller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-gin-admin/response"
	"go-gin-admin/util"
	"path"
	"strings"
	"time"
)

func UploadImg(ctx *gin.Context) {
	UploadTmp := viper.GetString("Upload.tmp")

	f, err := ctx.FormFile("imgfile")
	if err != nil {
		response.Fail(ctx, "上传失败", nil)
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			response.Fail(ctx, "上传失败!只允许png,jpg,gif,jpeg文件", nil)
			return
		}
		fileName := util.GetMd5String(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))

		/*fildDir:=fmt.Sprintf("%s/",time.Now().Month().String())
		isExist,_:= util.IsFileExist(UploadTmp+fildDir)
		if !isExist{
			if err:=os.Mkdir(UploadTmp+fildDir,os.ModePerm);err != nil {
				response.Fail(ctx,"创建文件夹失败",nil)
				return
			}
		}*/
		filepath := fmt.Sprintf("%s%s", fileName, fileExt)
		err := ctx.SaveUploadedFile(f, UploadTmp+filepath)
		if err != nil {
			response.Fail(ctx, "保存文件失败", nil)
			return
		}
		response.Success(ctx, gin.H{
			"path": filepath,
		}, "上传成功")
	}

}
