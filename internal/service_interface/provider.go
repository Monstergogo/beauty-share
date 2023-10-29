package service_interface

import (
	"github.com/Monstergogo/beauty-share/internal/app"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
)

func OssServiceProvider() OSSService {
	return &app.OssServiceImpl{}
}

func ShareServiceProvider(mongoDB repo_interface.MongoRepo) ShareService {
	return &app.ShareServiceImpl{
		MongoRepo: mongoDB,
	}
}
