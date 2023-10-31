package db

import (
	"context"
	"github.com/Monstergogo/beauty-share/init/logger"
	"github.com/Monstergogo/beauty-share/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

var mongoDB *mongo.Client

func InitMongoDB() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoDB, err = mongo.Connect(ctx, options.Client().ApplyURI(util.MongoURI))
	if err != nil {
		panic(err)
	}
	err = mongoDB.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	logger.GetLogger().Info("mongo db connected success")
}

func GetMongoDB() *mongo.Client {
	return mongoDB
}

func DisconnectMongoDB() error {
	if err := mongoDB.Disconnect(context.Background()); err != nil {
		logger.GetLogger().Error("disconnect mongo db err:%v", zap.Any("err", err))
		return err
	}
	logger.GetLogger().Info("disconnect mongo db success")
	return nil
}
