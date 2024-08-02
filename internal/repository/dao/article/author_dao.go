package article

import (
	"context"
	"gorm.io/gorm"
)

type AuthorDao interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

func NewAuthorDao(db *gorm.DB) AuthorDao {
	panic("implement me")
}
