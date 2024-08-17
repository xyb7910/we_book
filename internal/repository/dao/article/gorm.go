package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

func (g *GORMArticleDAO) ListPub(ctx context.Context, start time.Time, offset int, limit int) ([]Article, error) {
	var res []Article
	milli := start.UnixMilli()
	err := g.db.WithContext(ctx).Where("utime<?", milli).Order("utime desc").Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{db: db}
}

func (g *GORMArticleDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	var pub PublishedArticle
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&pub).Error
	return pub, err
}

func (g *GORMArticleDAO) GetById(ctx context.Context, id int64) (Article, error) {
	var article Article
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&article).First(&article).Error
	return article, err
}

func (g *GORMArticleDAO) GetByAuthor(ctx context.Context, uid int64, offset int, limit int) ([]Article, error) {
	var articles []Article
	err := g.db.WithContext(ctx).Where("author_id = ?", uid).
		Offset(offset).
		Limit(limit).
		Order("ctime desc").
		Find(&articles).Error
	return articles, err
}

func (g *GORMArticleDAO) Transaction(ctx context.Context, bizFunc func(txDAO ArticleDAO) error) error {
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDAO(tx)
		return bizFunc(txDAO)
	})
}

func (g *GORMArticleDAO) SyncStatus(ctx context.Context, id int64, author int64, u uint8) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? AND author_id = ?", id, author).
			Updates(map[string]any{
				"status": u,
				"utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return fmt.Errorf("update article status failed")
		}
		return tx.Model(&Article{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"status": u,
				"utime":  now,
			}).Error
	})
}

func (g *GORMArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {
	var id = article.Id
	var err error
	err = g.db.Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDAO(tx)
		if article.Id > 0 {
			err = txDAO.UpdateById(ctx, article)
		} else {
			id, err = txDAO.Insert(ctx, article)
		}
		if err != nil {
			return err
		}
		return txDAO.Upsert(ctx, article)
	})
	return id, err
}

func (g *GORMArticleDAO) Upsert(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now

	return g.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(
			map[string]interface{}{
				"title":   article.Title,
				"status":  article.Status,
				"content": article.Content,
				"utime":   article.Utime,
			}),
	}).Create(&article).Error
}

func (g *GORMArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := g.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (g *GORMArticleDAO) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now

	res := g.db.WithContext(ctx).Model(&article).
		Where("id = ? and author_id = ?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
			"status":  article.Status,
			"utime":   article.Utime,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("update article failed")
	}
	return res.Error
}
