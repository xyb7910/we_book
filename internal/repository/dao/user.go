package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("user duplicate email or phone")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindByWechatOpenId(ctx context.Context, OpenId string) (User, error)
	UpdateById(ctx context.Context, user User) error
}

type GORMUserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{db: db}
}

func (ud *GORMUserDAO) Insert(ctx context.Context, user User) error {
	// 创建时间 毫秒级
	now := time.Now().UnixMilli()
	user.Ctime = now
	user.Utime = now
	// 使用 1062 错误码判断是否是唯一索引冲突
	err := ud.db.WithContext(ctx).Create(&user).Error
	if sqlError, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrNo = 1062
		if sqlError.Number == uniqueIndexErrNo {
			return ErrUserDuplicate
		}
	}
	return err
}

func (ud *GORMUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

func (ud *GORMUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).First(&User{}, id).First(&user).Error
	return user, err
}

func (ud *GORMUserDAO) UpdateById(ctx context.Context, user User) error {
	return ud.db.WithContext(ctx).Model(user).Where("id = ?", user.Id).Updates(user).Error
}

func (ud *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	return user, err
}

func (ud *GORMUserDAO) FindByWechatOpenId(ctx context.Context, OpenId string) (User, error) {
	var user User
	err := ud.db.WithContext(ctx).Where("wechat_open_id = ?", OpenId).First(&user).Error
	return user, err
}

type User struct {
	Id       int64          `gorm:"primaryKey, autoIncrement"`
	Email    sql.NullString `gorm:"type:varchar(100);uniqueIndex"`
	Password string         `gorm:"type:varchar(100)"`

	// 创建时间 毫秒级
	Ctime int64
	// 更新时间 毫秒级
	Utime int64

	// 用户信息
	NickName      string `gorm:"type:varchar(100)"`
	Birthday      int64
	Introduction  string         `gorm:"type:varchar(100)"`
	WeChatOpenId  sql.NullString `gorm:"type:varchar(100);uniqueIndex"`
	WeChatUnionId sql.NullString `gorm:"type:varchar(100);uniqueIndex"`

	Phone sql.NullString `gorm:"type:varchar(100);uniqueIndex"`
}
