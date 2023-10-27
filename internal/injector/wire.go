//+build wireinject

package injector

import (
	"github.com/Monstergogo/beauty-share/internal/service_interface"
	"github.com/google/wire"
)

func GetOssServer() service_interface.OSSService {
	wire.Build(service_interface.OssServiceProvider)
	return nil
}
