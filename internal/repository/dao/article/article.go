package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	Upsert(ctx context.Context, article Article) error
	SyncStatus(ctx context.Context, id int64, author int64, u uint8) error
	Transaction(ctx context.Context, bizFunc func(txDAO ArticleDAO) error) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{db: db}
}

type Article struct {
	Id       int64  `gorm:"primaryKey,  autoIncrement"`
	Title    string `gorm:"type:varchar(255)"`
	Content  string `gorm:"type:text"`
	AuthorId int64  `gorm:"index"`
	Status   uint8
	Ctime    int64
	Utime    int64
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
