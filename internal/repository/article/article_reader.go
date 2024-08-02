package article

import (
	"context"
	"we_book/internal/domain"
)

type ArticleReaderRepository interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
}
