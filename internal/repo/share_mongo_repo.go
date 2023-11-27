package repo

import (
	"context"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/internal/entity"
	"github.com/Monstergogo/beauty-share/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MongoRepoImpl struct {
	DB *mongo.Client
}

func (s *MongoRepoImpl) GetShareByPage(ctx context.Context, lastId primitive.ObjectID, pageSize int64) (int64, []*entity.ShareInfo, error) {
	var (
		res   []*entity.ShareInfo
		total int64
	)

	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	findOptions.SetSort(bson.M{"created_at": -1})

	collection := s.DB.Database(util.MongoShareDBName).Collection(util.MongoShareCollectName)
	filter := bson.D{}
	if !lastId.IsZero() {
		filter = bson.D{{"_id", bson.M{"$gt": lastId}}}
	}

	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "get share info by page repo err", zap.Any("err_msg", err))
		return total, res, err
	}

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "get share info by page err", zap.Any("err_smg", err))
		return total, res, err
	}
	defer cur.Close(ctx)

	res = make([]*entity.ShareInfo, 0)
	for cur.Next(ctx) {
		var temp *entity.ShareInfo
		if err = cur.Decode(&temp); err != nil {
			logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "repo decode to share info err", zap.Any("err_smg", err))
			return total, res, err
		}
		res = append(res, temp)
	}
	return total, res, err
}

func (s *MongoRepoImpl) AddShare(ctx context.Context, shareInfo entity.ShareInfo) error {
	collection := s.DB.Database(util.MongoShareDBName).Collection(util.MongoShareCollectName)
	_, err := collection.InsertOne(ctx, shareInfo)
	if err != nil {
		logger.LogWithTraceId(ctx, zapcore.ErrorLevel, "mongo insert one err", zap.Any("err_msg", err))
		return err
	}
	return err
}
