package repo

import (
	"context"
	myDB "github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/internal/entity"
	"github.com/Monstergogo/beauty-share/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoRepoImpl struct {
}

func (s *MongoRepoImpl) GetShareByPage(ctx context.Context, lastId primitive.ObjectID, pageSize int64) (int64, []*entity.ShareInfo, error) {
	var (
		res   []*entity.ShareInfo
		total int64
	)
	db := myDB.GetMongoDB()

	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	findOptions.SetSort(bson.M{"created_at": -1})

	collection := db.Database(util.MongoShareDBName).Collection(util.MongoShareCollectName)
	filter := bson.D{{"_id", bson.M{"$gt": lastId}}}
	cur, err := collection.Find(ctx, filter, findOptions)
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
	total, err = collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatalln(err)
		return total, res, err
	}
	return total, res, err
}

func (s *MongoRepoImpl) AddShare(ctx context.Context, shareInfo entity.ShareInfo) error {
	db := myDB.GetMongoDB()
	collection := db.Database(util.MongoShareDBName).Collection(util.MongoShareCollectName)
	res, err := collection.InsertOne(ctx, shareInfo)
	if err != nil {
		log.Printf("mongo insert one err:%v", err)
		return err
	}
	log.Println(res.InsertedID)
	return err
}
