package repo_interface

import (
	"context"
	"github.com/Monstergogo/beauty-share/internal/entity"
)

type MongoRepo interface {
	AddShare(ctx context.Context, shareInfo interface{}) error
	GetShareByPage(ctx context.Context, pageIndex, pageSize int64) (int64, []*entity.ShareInfo, error)
}
