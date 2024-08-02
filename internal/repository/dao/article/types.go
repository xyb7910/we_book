package article

import "context"

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, article Article) error
	SyncStatus(ctx context.Context, id int64, author int64, u uint8) error
	Transaction(ctx context.Context, bizFunc func(txDAO ArticleDAO) error) error
}
