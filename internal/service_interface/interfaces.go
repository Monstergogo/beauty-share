package service_interface

import (
	"context"
	"mime/multipart"
)

type OSSService interface {
	ObjectUpload(ctx context.Context, files []*multipart.FileHeader) (fileUrl []string, err error)
}
