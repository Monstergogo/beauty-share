package rest

import (
	"github.com/Monstergogo/beauty-share/cmd/server/interface"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
)

func toGinHandler(handler func(*gin.Context) (interface{}, error)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			util.GinErrResponse(ctx, err)
			return
		}
		util.GinSuccessResponse(ctx, res)
	}
}

func InitRouter(server _interface.MicroServer) {
	server.RegisterGinRouter(_interface.HttpServerRouter{
		Path:    "/v1/oss/upload",
		Method:  util.HttpMethodPost,
		Handler: toGinHandler(OssUpload),
	})
}
