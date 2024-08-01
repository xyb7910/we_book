package article

import (
	"context"
	"gorm.io/gorm"
	"we_book/internal/domain"
	"we_book/internal/repository/dao/article"
)

type CacheArticleRepository struct {
	dao article.ArticleDAO

	authorDAO article.AuthorDao
	readerDAO article.ArticleDAO

	db *gorm.DB
}

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Sync(ctx context.Context, article domain.Article) (int64, error)
}

func NewArticleRepository(dao article.ArticleDAO, author article.AuthorDao, reader article.ArticleDAO) ArticleRepository {
	return &CacheArticleRepository{dao: dao, authorDAO: author, readerDAO: reader}
}

func (c CacheArticleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}

// Sync 数据在同一个表中
func (c CacheArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Sync(ctx, c.toEntity(art))
}

// SyncV1 数据存在统一数据库但是在不同的表中---非事务
func (c CacheArticleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
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

func (c CacheArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
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

func (c CacheArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

func (c CacheArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
