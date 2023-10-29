//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
	"github.com/Monstergogo/beauty-share/internal/service_interface"
	"github.com/google/wire"
)

func GetOssServer() service_interface.OSSService {
	wire.Build(service_interface.OssServiceProvider)
	return nil
}

func GetShareServer() service_interface.ShareService {
	wire.Build(repo_interface.MongoRepoProvider, service_interface.ShareServiceProvider)
	return nil
}
