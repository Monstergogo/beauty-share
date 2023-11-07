package minio

import (
	"github.com/Monstergogo/beauty-share/init/nacos"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var minioClient *minio.Client

func InitMinio() {
	minioEndpoint, err := nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
		DataId: util.MinioEndpointDataID,
	})
	if err != nil {
		panic(err)
	}

	minioID, err := nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
		DataId: util.MinioIDDataID,
	})
	if err != nil {
		panic(err)
	}

	minioSecret, err := nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
		DataId: util.MinioSecretDataID,
	})
	if err != nil {
		panic(err)
	}

	// Initialize minio client object.
	minioClient, err = minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioID, minioSecret, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
}

func GetMinioClient() *minio.Client {
	return minioClient
}
