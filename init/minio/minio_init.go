package minio

import (
	"github.com/Monstergogo/beauty-share/conf"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func InitMinio() {
	//minioEndpoint, err := nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
	//	DataId: util.MinioEndpointDataID,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//minioID, err := nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
	//	DataId: util.MinioIDDataID,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//minioSecret, err := nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
	//	DataId: util.MinioSecretDataID,
	//})
	//if err != nil {
	//	panic(err)
	//}

	var err error
	// Initialize minio client object.
	minioClient, err = minio.New(conf.ServerConf.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.ServerConf.Minio.ID, conf.ServerConf.Minio.Secret, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
}

func GetMinioClient() *minio.Client {
	return minioClient
}
