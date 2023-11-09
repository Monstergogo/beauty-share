package rest

import (
	"github.com/Monstergogo/beauty-share/init/server"
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

func InitRouter(s server.MicroServer) {
	s.RegisterGinRouter(server.HttpServerRouter{
		Path:    "/v1/oss/upload",
		Method:  util.HttpMethodPost,
		Handler: toGinHandler(OssUpload),
	})
	s.RegisterGinRouter(server.HttpServerRouter{
		Path:    "v1/ping",
		Method:  util.HttpMethodGet,
		Handler: toGinHandler(Ping),
	})
}
