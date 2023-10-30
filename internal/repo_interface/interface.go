package repo_interface

import (
	"context"
	"github.com/Monstergogo/beauty-share/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoRepo interface {
	AddShare(ctx context.Context, shareInfo entity.ShareInfo) error
	GetShareByPage(ctx context.Context, lastID primitive.ObjectID, pageSize int64) (int64, []*entity.ShareInfo, error)
}
