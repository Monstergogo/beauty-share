package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Monstergogo/beauty-share/init/logger"
	myMinio "github.com/Monstergogo/beauty-share/init/minio"
	"github.com/Monstergogo/beauty-share/init/nacos"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/minio/minio-go/v7"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"mime/multipart"
	"strings"
	"sync"
)

var contentType = map[string]string{
	"gif":  "image/gif",
	"jpeg": "image/jpeg",
	"jpg":  "image/jpg",
	"png":  "image/png",
	"pdf":  "application/pdf",
}

type OssServiceImpl struct {
}

func getFileExtensionByFilename(filename string) string {
	filenameSplit := strings.Split(filename, ".")
	fileType := filenameSplit[len(filenameSplit)-1]
	return fileType
}

func getUploadContentType(fileType string) string {
	var uploadContentType string
	if _, ok := contentType[fileType]; ok {
		uploadContentType = contentType[fileType]
	} else {
		uploadContentType = "application/octet-stream"
	}
	return uploadContentType
}

// 根据雪花算法生成相对有序的文件名
func genFilenameAscBySnow(fileNum int) ([]int64, error) {
	snowNode, err := util.NewWorker(1)
	if err != nil {
		logger.GetLogger().Error("init snow worker err", zap.Any("err_msg", err))
		return nil, errors.New("get snow work node er")
	}
	res := make([]int64, fileNum)
	for i := 0; i < fileNum; i++ {
		res[i] = snowNode.GetId()
	}
	return res, err
}

// 获取endpoint和bucket配置信息
func getMinioEndpointAndBucketName(ctx context.Context) (endpoint, bucketName string, err error) {
	endpoint, err = nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
		DataId: util.MinioEndpointDataID,
	})
	if err != nil {
		logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "upload object to bucket err", zap.Any("err_msg", err))
		return
	}
	bucketName, err = nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
		DataId: util.MinioShareBucketDataID,
	})
	if err != nil {
		logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "upload object to bucket err", zap.Any("err_msg", err))
		return
	}
	return
}

func (o *OssServiceImpl) ObjectUpload(ctx context.Context, files []*multipart.FileHeader) (fileUrl []string, err error) {
	filenameAsc, err := genFilenameAscBySnow(len(files))
	if err != nil {
		logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "upload file failed", zap.Any("err_msg", err))
		return nil, errors.New("upload file failed")
	}

	var wg sync.WaitGroup
	wg.Add(len(files))
	fileUrl = make([]string, 0)
	for index, _ := range files {
		go func(index int) {
			defer wg.Done()

			var file multipart.File
			file, err = files[index].Open()
			defer file.Close()
			if err != nil {
				logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "file open err", zap.Any("err_msg", err))
				return
			}
			fileSize := files[index].Size
			filename := files[index].Filename
			fileExtension := getFileExtensionByFilename(filename)
			uploadFilename := fmt.Sprintf("%d.%s", filenameAsc[index], fileExtension)
			uploadContentType := getUploadContentType(fileExtension)

			minioEndpoint, bucketName, err := getMinioEndpointAndBucketName(ctx)
			if err != nil {
				return
			}
			_, err = myMinio.GetMinioClient().PutObject(ctx, bucketName, uploadFilename, file, fileSize, minio.PutObjectOptions{ContentType: uploadContentType})
			if err != nil {
				logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "upload object to bucket err", zap.Any("err_msg", err))
				return
			}
			fileUrl = append(fileUrl, fmt.Sprintf("%s:%s/%s/%s", util.MinioNetProtocol, minioEndpoint, bucketName, uploadFilename))
		}(index)
	}
	wg.Wait()
	if err != nil {
		return fileUrl, errors.New("upload file failed")
	}
	return
}
