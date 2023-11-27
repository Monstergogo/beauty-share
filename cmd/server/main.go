package main

import (
	"github.com/Monstergogo/beauty-share/init/conf"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/init/gin"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/init/minio"
	"github.com/Monstergogo/beauty-share/init/server"
	"github.com/Monstergogo/beauty-share/init/tracing"
	"github.com/Monstergogo/beauty-share/internal/app"
)

func main() {
	logger.Init(conf.Log.LogPath, conf.Log.ErrPath)
	db.InitMongoDB()
	minio.Init()
	tracing.InitProvider()

	ginServer := gin.Init()
	// 初始化gin路由
	gin.InitRouter()
	consulServer := new(app.ConsulServiceImpl)
	s := server.MicroServer{GinServer: ginServer, ConsulServer: consulServer}
	s.RunServer()
}
