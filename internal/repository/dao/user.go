package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("user duplicate email")
	ErrUserNotFound       = gorm.ErrRecordNotFound
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
	// 使用 1062 错误码判断是否是唯一索引冲突
	err := ud.db.WithContext(ctx).Create(&user).Error
	if sqlError, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo = 1062
		if sqlError.Number == uniqueIndexErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (ud *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (ud *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).First(&User{}, id).First(&user).Error
	return user, err
}

func (ud *UserDAO) UpdateById(ctx context.Context, user User) error {
	return ud.db.WithContext(ctx).Model(user).Where("id = ?", user.Id).Updates(user).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey, autoIncrement"`
	Email    string `gorm:"type:varchar(100);uniqueIndex"`
	Password string `gorm:"type:varchar(100)"`

	// 创建时间 毫秒级
	Ctime int64
	// 更新时间 毫秒级
	Utime int64

	// 用户信息
	NickName     string `gorm:"type:varchar(100)"`
	Birthday     int64
	Introduction string `gorm:"type:varchar(100)"`
}
