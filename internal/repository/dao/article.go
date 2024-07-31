package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type GORMArticleDAO struct {
	db *gorm.DB
}

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

func NewGORMArticleDAO(db *gorm.DB) ArticleDAO {
	return &GORMArticleDAO{db: db}
}

type Article struct {
	Id       int64  `gorm:"primaryKey,  autoIncrement"`
	Title    string `gorm:"type:varchar(255)"`
	Content  string `gorm:"type:text"`
	AuthorId int64  `gorm:"index"`
	Ctime    int64
	Utime    int64
}

func (g GORMArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now
	err := g.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (g GORMArticleDAO) UpdateById(ctx context.Context, article Article) error {
	now := time.Now().UnixMilli()
	article.Utime = now

	res := g.db.WithContext(ctx).Model(&article).
		Where("id = ? and author_id = ?", article.Id, article.AuthorId).
		Updates(map[string]any{
			"title":   article.Title,
			"content": article.Content,
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
