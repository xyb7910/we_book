package service

import (
	"context"
	"time"
	events "we_book/events/article"
	"we_book/internal/domain"
	"we_book/internal/repository/article"
	"we_book/pkg/logger"
)

type articleService struct {
	repo       article.ArticleRepository
	readerRepo article.ArticleReaderRepository
	authorRepo article.ArticleAuthorRepository
	l          logger.V1
	producer   events.Producer
	ch         chan readInfo
}

type readInfo struct {
	Uid int64
	Aid int64
}

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article) (int64, error)
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) (int64, error)
	Withdraw(ctx context.Context, article domain.Article) error
	List(ctx context.Context, uid int64, set int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, aid, uid int64) (domain.Article, error)
}

func NewArticleService(repo article.ArticleRepository, l logger.V1, producer events.Producer) ArticleService {
	return &articleService{
		repo:     repo,
		l:        l,
		producer: producer,
		//ch: make(chan readInfo, 10),
	}
}

func NewArticleServiceV2(repo article.ArticleRepository,
	l logger.V1, producer events.Producer) ArticleService {
	ch := make(chan readInfo, 10)
	go func() {
		for {
			Uids := make([]int64, 0, 10)
			Aids := make([]int64, 0, 10)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			for i := 0; i < 10; i++ {
				select {
				case <-ctx.Done():
					break
				case Info, ok := <-ch:
					if !ok {
						cancel()
						return
					}
					Uids = append(Uids, Info.Uid)
					Aids = append(Aids, Info.Aid)
				}
			}
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			producer.ProducerReadEventV1(ctx, events.ReadEventV1{
				Uids: Uids,
				Aids: Aids,
			})
			cancel()
		}
	}()
	return &articleService{
		repo:     repo,
		producer: producer,
		l:        l,
		ch:       ch,
	}
}

func (asv *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	return asv.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func NewArticleServiceV1(readerRepo article.ArticleReaderRepository, authorRepo article.ArticleAuthorRepository) ArticleService {
	return &articleService{readerRepo: readerRepo, authorRepo: authorRepo}
}

func (asv *articleService) GetPubById(ctx context.Context, aid, uid int64) (domain.Article, error) {
	art, err := asv.repo.GetPubById(ctx, aid)
	if err == nil {
		go func() {
			er := asv.producer.ProducerReadEvent(
				ctx,
				events.ReadEvent{
					Uid: uid,
					Aid: aid,
				})
			if er == nil {
				asv.l.Error("producer read event error")
			}
		}()

		go func() {
			asv.ch <- readInfo{
				Uid: uid,
				Aid: aid,
			}
		}()
	}
	return art, err
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
