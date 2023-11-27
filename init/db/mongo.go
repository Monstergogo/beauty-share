package db

import (
	"context"
	"github.com/Monstergogo/beauty-share/init/conf"
	"github.com/Monstergogo/beauty-share/init/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.uber.org/zap"
	"sync"
	"time"
)

type mongoClientStruct struct {
	Client *mongo.Client
	Lock   *sync.RWMutex
}

var mongoClient mongoClientStruct

func mongoDBConnectAndCheckHearty(ctx context.Context, mongoUri string) (*mongo.Client, error) {
	opt := options.Client()
	opt.Monitor = otelmongo.NewMonitor()
	opt.ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, opt)

	if err != nil {
		return client, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return client, err
	}
	return client, err
}

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//mongoUri, err := nacos.GetNacosConfigClient().GetConfig(vo.ConfigParam{
	//	DataId: util.MongoUriDataID,
	//})
	//if err != nil {
	//	panic(err)
	//}

	client, err := mongoDBConnectAndCheckHearty(ctx, conf.Mongo.Uri)
	if err != nil {
		panic(err)
	}

	mongoClient = mongoClientStruct{
		Client: client,
		Lock:   new(sync.RWMutex),
	}
	// 监听mongo_uri变化，建立新连接
	//nacos.GetNacosConfigClient().ListenConfig(vo.ConfigParam{
	//	DataId: util.MongoUriDataID,
	//	OnChange: func(namespace, group, dataId, data string) {
	//		c, err := mongoDBConnectAndCheckHearty(ctx, data)
	//		if err != nil {
	//			logger.GetLogger().Error("mongo uri changed but connected err", zap.Any("err_msg", err))
	//			return
	//		}
	//		DisconnectMongoDB()
	//		SetMongoDBClient(c)
	//	},
	//})
	logger.GetLogger().Info("mongo db connected success")
}

func GetMongoDB() *mongo.Client {
	mongoClient.Lock.RLock()
	defer mongoClient.Lock.RUnlock()

	return mongoClient.Client
}

// SetMongoDBClient 更新mongo client
func SetMongoDBClient(client *mongo.Client) {
	mongoClient.Lock.Lock()
	defer mongoClient.Lock.Unlock()

	mongoClient.Client = client
}

func DisconnectMongoDB() error {
	if err := mongoClient.Client.Disconnect(context.Background()); err != nil {
		logger.GetLogger().Error("disconnect mongo db err:%v", zap.Any("err", err))
		return err
	}
	logger.GetLogger().Info("disconnect mongo db success")
	return nil
}
