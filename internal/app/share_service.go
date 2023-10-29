package app

import (
	"context"
	"errors"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
)

type ShareServiceImpl struct {
	*pb.UnimplementedShareServiceServer
	MongoRepo repo_interface.MongoRepo
}

func (s ShareServiceImpl) AddShare(ctx context.Context, in *pb.AddShareReq) (*pb.AddShareResp, error) {
	var resp *pb.AddShareResp
	if in.PostContent.Text == "" && len(in.PostContent.Img) == 0 {
		return resp, errors.New("params err")
	}
	err := s.MongoRepo.AddShare(ctx, in.PostContent)
	if err != nil {
		return resp, err
	}
	return resp, err
}

func (s ShareServiceImpl) GetShareByPage(ctx context.Context, in *pb.GetShareByPageReq) (*pb.GetShareByPageResp, error) {
	var resp *pb.GetShareByPageResp
	total, queryItem, err := s.MongoRepo.GetShareByPage(ctx, in.PageIndex, in.PageSize)
	if err != nil {
		return resp, err
	}
	resp = new(pb.GetShareByPageResp)
	resp.Total = total
	resp.Data = make([]*pb.PostItem, len(queryItem))
	for index, item := range queryItem {
		resp.Data[index] = &pb.PostItem{
			Text: item.Text,
			Img:  item.Images,
		}
	}
	return resp, err
}
