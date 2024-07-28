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
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (as articleService) Edit(ctx context.Context, article domain.Article) (int64, error) {
	return as.repo.Create(ctx, article)
}
