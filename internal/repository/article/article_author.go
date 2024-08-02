package article

import (
	"context"
	"we_book/internal/domain"
)

type ArticleAuthorRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}
