package dao

import (
	"gorm.io/gorm"
	"we_book/internal/repository/dao/article"
)

func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.Article{}, &article.PublishArticleDAO{})
}
