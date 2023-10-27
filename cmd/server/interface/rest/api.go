package rest

import (
	"github.com/Monstergogo/beauty-share/internal/injector"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
)

func OssUpload(ctx *gin.Context) (interface{}, error) {
	form, _ := ctx.MultipartForm()
	files := form.File[util.FileUploadRequestName]
	return injector.GetOssServer().ObjectUpload(ctx, files)
}
