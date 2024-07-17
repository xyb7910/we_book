package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (ud *UserDAO) Insert(ctx context.Context, user User) error {
	// 创建时间 毫秒级
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	return ud.db.WithContext(ctx).Create(&user).Error
}

type User struct {
	Id       int    `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"type:varchar(100);uniqueIndex"`
	Password string `gorm:"type:varchar(100)"`

	// 创建时间 毫秒级
	Ctime int64
	// 更新时间 毫秒级
	Utime int64
}
