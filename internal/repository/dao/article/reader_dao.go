package article

import (
	"context"
	"gorm.io/gorm"
)

type ReaderDAO interface {
	Upsert(ctx context.Context, article Article) error
	UpsertV2(ctx context.Context, dao PublishArticleDAO) error
}

type PublishArticleDAO struct {
	Article
}

func NewReaderDAO(db *gorm.DB) ReaderDAO {
	panic("implement me")
}
