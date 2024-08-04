package article

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"gorm.io/gorm"
	"time"
	"we_book/internal/domain"
	"we_book/internal/repository"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao/article"
	dao "we_book/internal/repository/dao/article"
	"we_book/pkg/logger"
)

type CacheArticleRepository struct {
	dao article.ArticleDAO

	userRepo repository.UserRepository
	cache    cache.ArticleCache

	authorDAO article.AuthorDao
	readerDAO article.ArticleDAO

	db *gorm.DB
	l  logger.V1
}

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, article domain.Article) (int64, error)
	SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, set int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, id int64) (domain.Article, error)
}

func NewArticleRepository(dao article.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{dao: dao}
}

func (c CacheArticleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}

func (c *CacheArticleRepository) GetPubById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	usr, err := c.userRepo.FindById(ctx, art.AuthorId)
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.ArticleStatus(art.Status),
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.NickName,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil
}

func (c *CacheArticleRepository) GetById(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(data), nil
}

func (c *CacheArticleRepository) SyncStatus(ctx context.Context, id int64, author int64, status domain.ArticleStatus) error {
	return c.dao.SyncStatus(ctx, id, author, uint8(status))
}

// Sync 数据在同一个表中
func (c *CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := c.dao.Sync(ctx, c.toEntity(art))
	if err == nil {
		c.cache.DelFirstPage(ctx, art.Author.Id)
		er := c.cache.Set(ctx, art)
		if er != nil {
			fmt.Errorf("set cache err: %v", er)
		}
	}
	return id, err
}

// SyncV1 数据存在统一数据库但是在不同的表中---非事务
func (c *CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  int64
		err error
	)
	if art.Id > 0 {
		err = c.readerDAO.UpdateById(ctx, c.toEntity(art))
	} else {
		id, err = c.readerDAO.Insert(ctx, c.toEntity(art))
	}
	if err != nil {
		return 0, err
	}
	err = c.readerDAO.Upsert(ctx, c.toEntity(art))
	return id, err
}

func (c *CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	// 首先开启一个事务
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	defer tx.Rollback()

	// 使用 tx 来构建 dao
	author := article.NewAuthorDao(tx)
	reader := article.NewReaderDAO(tx)

	var (
		id  int64
		err error
	)
	if art.Id > 0 {
		err = author.UpdateById(ctx, c.toEntity(art))
	} else {
		id, err = author.Insert(ctx, c.toEntity(art))
	}
	if err != nil {
		return 0, err
	}
	err = reader.UpsertV2(ctx, article.PublishArticleDAO{Article: c.toEntity(art)})
	tx.Commit()
	return id, err
}

func (c *CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	defer func() {
		err := c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			fmt.Errorf("del first page err: %v", err)
		}
	}()
	return c.dao.Insert(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	})
}

func (c *CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	defer func() {
		err := c.cache.DelFirstPage(ctx, art.Author.Id)
		if err != nil {
			fmt.Errorf("del first page err: %v", err)
		}
	}()
	return c.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	})
}

func (c *CacheArticleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	// 优先返回第一页数据
	if offset == 0 && limit <= 100 {
		data, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				c.preCache(ctx, data)
			}()
			return data, nil
		}
	}
	res, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[dao.Article, domain.Article](res, func(idx int, item dao.Article) domain.Article {
		return c.toDomain(item)
	})
	// 注意回写入缓存
	go func() {
		err = c.cache.SetFirstPage(ctx, uid, data)
		fmt.Errorf("set first page err: %v", err)
		c.preCache(ctx, data)
	}()
	return data, nil
}

func (c *CacheArticleRepository) preCache(ctx context.Context, data []domain.Article) {
	const MAX_CACHE_SIZE = 1024 * 1024
	if len(data) > 0 && len(data[0].Content) < MAX_CACHE_SIZE {
		err := c.cache.Set(ctx, data[0])
		if err != nil {
			fmt.Errorf("set cache err: %v", err)
		}
	}
}

func (c *CacheArticleRepository) toDomain(item article.Article) domain.Article {
	return domain.Article{
		Id:      item.Id,
		Title:   item.Title,
		Content: item.Content,
		Status:  domain.ArticleStatus(item.Status),
		Author: domain.Author{
			Id: item.AuthorId,
		},
		Ctime: time.UnixMilli(item.Ctime),
		Utime: time.UnixMilli(item.Utime),
	}
}
