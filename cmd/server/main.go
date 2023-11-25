package main

import (
	"github.com/Monstergogo/beauty-share/init/conf"
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/init/gin"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/init/minio"
	"github.com/Monstergogo/beauty-share/init/server"
)

func main() {
	logger.Init(conf.Log.LogPath, conf.Log.ErrPath)
	db.InitMongoDB()
	minio.Init()

	ginServer := gin.Init()
	// 初始化gin路由
	gin.InitRouter()
	s := server.MicroServer{GinServer: ginServer}
	s.RunServer()
}
