package service

import (
	"context"
	"we_book/internal/domain"
	"we_book/internal/repository/article"
)

type articleService struct {
	repo       article.ArticleRepository
	readerRepo article.ArticleReaderRepository
	authorRepo article.ArticleAuthorRepository
}

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) (int64, error)
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, set int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{repo: repo}
}

func (asv *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	return asv.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func NewArticleServiceV1(readerRepo article.ArticleReaderRepository, authorRepo article.ArticleAuthorRepository) ArticleService {
	return &articleService{readerRepo: readerRepo, authorRepo: authorRepo}
}

func (asv *articleService) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	return asv.repo.GetPubById(ctx, id)
}

func (asv *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return asv.repo.GetById(ctx, id)
}

func (asv *articleService) List(ctx context.Context, uid int64, set int, limit int) ([]domain.Article, error) {
	return asv.repo.List(ctx, uid, set, limit)
}

func (asv *articleService) Edit(ctx context.Context, article domain.Article) (int64, error) {
	return asv.repo.Create(ctx, article)
}

func (asv *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	article.Status = domain.ArticleStatusUnpublished
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

func (asv *articleService) PublishV1(ctx context.Context, article domain.Article) (int64, error) {
	panic("implement me")
}

func (asv *articleService) Publish(ctx context.Context, article domain.Article) (int64, error) {
	var (
		id  = article.Id
		err error
	)
	if article.Id > 0 {
		err = asv.authorRepo.Update(ctx, article)
	} else {
		id, err = asv.authorRepo.Create(ctx, article)
	}

	if err != nil {
		return 0, err
	}

	article.Id = id
	for i := 0; i < 3; i++ {
		id, err = asv.readerRepo.Save(ctx, article)
		if err == nil {
			break
		}
	}
	return id, nil
}
