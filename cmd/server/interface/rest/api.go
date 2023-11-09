package rest

import (
	"fmt"
	"github.com/Monstergogo/beauty-share/internal/injector"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
	"time"
)

// OssUpload 上传文件
func OssUpload(ctx *gin.Context) (interface{}, error) {
	form, _ := ctx.MultipartForm()
	files := form.File[util.FileUploadRequestName]
	return injector.GetOssServer().ObjectUpload(ctx, files)
}

func Ping(ctx *gin.Context) (interface{}, error) {
	return fmt.Sprintf("Ping: [%s]", time.Now().Format("2006-01-02 15:04:05")), nil
}
