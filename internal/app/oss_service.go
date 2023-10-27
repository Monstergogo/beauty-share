package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Monstergogo/beauty-share/util"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"strings"
	"sync"
)

var minIOClient *minio.Client

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

func (receiver *OssServiceImpl) ObjectUpload(ctx context.Context, files []*multipart.FileHeader) (fileUrl []string, err error) {
	// Initialize minio client object.
	minioClient, err := minio.New(util.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(util.MinioID, util.MinioSecret, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
		return fileUrl, err
	}

	snowNode, err := util.NewWorker(1)
	if err != nil {
		log.Fatalln(err)
		return fileUrl, errors.New("get snow work node er")
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
				log.Fatalln(err)
				return
			}
			fileSize := files[index].Size
			filename := files[index].Filename
			fileExtension := getFileExtensionByFilename(filename)
			uploadContentType := getUploadContentType(fileExtension)
			// 文件唯一id: 雪花id.文件后缀
			uploadFilename := fmt.Sprintf("%d.%s", snowNode.GetId(), fileExtension)

			_, err = minioClient.PutObject(ctx, util.BucketName, uploadFilename, file, fileSize, minio.PutObjectOptions{ContentType: uploadContentType})
			if err != nil {
				log.Fatalln(err)
				return
			}
			fileUrl = append(fileUrl, fmt.Sprintf("%s:%s/%s/%s", util.MinioNetProtocol, util.MinioEndpoint, util.BucketName, uploadFilename))
		}(index)
	}
	wg.Wait()
	return
}
