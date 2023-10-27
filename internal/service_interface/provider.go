package service_interface

import "github.com/Monstergogo/beauty-share/internal/app"

func OssServiceProvider() OSSService {
	return &app.OssServiceImpl{}
}
