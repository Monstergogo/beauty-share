package service_interface

import (
	"context"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"mime/multipart"
)

type OSSService interface {
	ObjectUpload(ctx context.Context, files []*multipart.FileHeader) (fileUrl []string, err error)
}

type ShareService interface {
	AddShare(ctx context.Context, in *pb.AddShareReq) (*pb.AddShareResp, error)
	GetShareByPage(ctx context.Context, in *pb.GetShareByPageReq) (*pb.GetShareByPageResp, error)
}
