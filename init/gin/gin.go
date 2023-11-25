package gin

import (
	"github.com/Monstergogo/beauty-share/util"
	"github.com/gin-gonic/gin"
)

type HttpServerRouter struct {
	Path    string
	Method  util.HttpMethod
	Handler func(context *gin.Context)
}

var ginServer *gin.Engine

func Init() *gin.Engine {
	ginServer = gin.Default()
	return ginServer
}

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

func InitRouter() {
	registerGinRouter(HttpServerRouter{
		Path:    "/v1/oss/upload",
		Method:  util.HttpMethodPost,
		Handler: toGinHandler(OssUpload),
	})
	registerGinRouter(HttpServerRouter{
		Path:    "v1/ping",
		Method:  util.HttpMethodGet,
		Handler: toGinHandler(Ping),
	})
}

func registerGinRouter(router HttpServerRouter) {
	switch router.Method {
	case util.HttpMethodGet:
		ginServer.GET(router.Path, router.Handler)
	case util.HttpMethodPost:
		ginServer.POST(router.Path, router.Handler)
	case util.HttpMethodPut:
		ginServer.PUT(router.Path, router.Handler)
	case util.HttpMethodDelete:
		ginServer.DELETE(router.Path, router.Handler)
	case util.HttpMethodPatch:
		ginServer.PATCH(router.Path, router.Handler)
	}
}
