package repository

import (
	"context"
	"we_book/internal/domain"
	"we_book/internal/repository/dao"
)

type CacheArticleRepository struct {
	dao dao.ArticleDAO
}

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{dao: dao}
}

func (c CacheArticleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}
