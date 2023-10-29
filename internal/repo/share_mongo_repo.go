package repo

import (
	"context"
	myDB "github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/internal/entity"
	"github.com/Monstergogo/beauty-share/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoRepoImpl struct {
}

func (s *MongoRepoImpl) GetShareByPage(ctx context.Context, pageIndex, pageSize int64) (int64, []*entity.ShareInfo, error) {
	var (
		res   []*entity.ShareInfo
		total int64
	)
	db := myDB.GetMongoDB()

	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	findOptions.SetSkip((pageIndex - 1) * pageSize)
	findOptions.SetSort(bson.M{"created_at": -1})

	collection := db.Database(util.MongoShareDBName).Collection(util.MongoShareCollectName)
	cur, err := collection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Printf("get share info by page err: %e", err)
		return total, res, err
	}
	defer cur.Close(ctx)
	res = make([]*entity.ShareInfo, 0)
	for cur.Next(ctx) {
		var temp *entity.ShareInfo
		if err = cur.Decode(&temp); err != nil {
			log.Fatalln(err)
			return total, res, err
		}
		res = append(res, temp)
	}
	total, err = collection.CountDocuments(ctx, &bson.M{})
	if err != nil {
		log.Fatalln(err)
		return total, res, err
	}
	return total, res, err
}

func (s *MongoRepoImpl) AddShare(ctx context.Context, shareInfo interface{}) error {
	db := myDB.GetMongoDB()
	err := db.Ping(ctx, nil)
	if err != nil {
		log.Printf("mongo ping err:%v", err)
		return err
	}
	collection := db.Database(util.MongoShareDBName).Collection(util.MongoShareCollectName)
	res, err := collection.InsertOne(ctx, shareInfo)
	if err != nil {
		log.Printf("mongo insert one err:%v", err)
		return err
	}
	log.Println(res.InsertedID)
	return err
}
