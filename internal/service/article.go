package service

import (
	"context"
	"we_book/internal/domain"
	"we_book/internal/repository"
)

type articleService struct {
	repo repository.ArticleRepository
}

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) (int64, error)
	Save(ctx context.Context, article domain.Article) (int64, error)
	create(ctx context.Context, article domain.Article) (int64, error)
	update(ctx context.Context, article domain.Article) error
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (asv *articleService) Edit(ctx context.Context, article domain.Article) (int64, error) {
	return asv.repo.Create(ctx, article)
}

func (asv *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	// 如果 article 的 ID 为 0 则调用 create 方法
	if article.Id > 0 {
		err := asv.update(ctx, article)
		return article.Id, err
	}
	return asv.create(ctx, article)
}

func (asv *articleService) create(ctx context.Context, article domain.Article) (int64, error) {
	return asv.repo.Create(ctx, article)
}

func (asv *articleService) update(ctx context.Context, article domain.Article) error {
	return asv.repo.Update(ctx, article)
}

func (asv *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	panic("implement me")
}
