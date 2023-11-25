package repo_interface

import (
	"github.com/Monstergogo/beauty-share/init/db"
	"github.com/Monstergogo/beauty-share/internal/repo"
	"go.mongodb.org/mongo-driver/mongo"
)

func MongoDBClientProvider() *mongo.Client {
	return db.GetMongoDB()
}

func MongoRepoProvider(dbClient *mongo.Client) MongoRepo {
	return &repo.MongoRepoImpl{
		DB: dbClient,
	}
}
