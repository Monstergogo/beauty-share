package app

import (
	"context"
	"errors"
	pb "github.com/Monstergogo/beauty-share/api/protobuf-spec"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/internal/entity"
	"github.com/Monstergogo/beauty-share/internal/repo_interface"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"time"
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
	shareInfo := entity.ShareInfo{
		Text:      in.PostContent.GetText(),
		Images:    in.PostContent.GetImg(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.MongoRepo.AddShare(ctx, shareInfo)
	if err != nil {
		return resp, err
	}
	resp = new(pb.AddShareResp)
	resp.Message = "success"
	return resp, err
}

func (s ShareServiceImpl) GetShareByPage(ctx context.Context, in *pb.GetShareByPageReq) (*pb.GetShareByPageResp, error) {
	var resp *pb.GetShareByPageResp
	if in.GetLastId() == "" || in.GetPageSize() == 0 {
		return resp, errors.New("params err")
	}
	lastId, err := primitive.ObjectIDFromHex(in.LastId)
	if err != nil {
		logger.GetLogger().Error("trans to object id err", zap.Any("err_msg", err))
		return resp, errors.New("lastId err")
	}
	total, queryItem, err := s.MongoRepo.GetShareByPage(ctx, lastId, in.PageSize)
	if err != nil {
		return resp, err
	}
	resp = new(pb.GetShareByPageResp)
	resp.Total = total
	resp.Data = make([]*pb.PostItem, len(queryItem))
	timeLayout := "2006-01-02 15:04:05"
	for index, item := range queryItem {
		resp.Data[index] = &pb.PostItem{
			Text:      item.Text,
			Img:       item.Images,
			CreatedAt: item.CreatedAt.Local().Format(timeLayout),
			UpdatedAt: item.UpdatedAt.Local().Format(timeLayout),
		}
	}
	if len(queryItem) > 0 {
		resp.LastId = queryItem[0].ID.Hex()
	}
	return resp, err
}
