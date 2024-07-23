package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"time"
	"we_book/internal/domain"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Edit(ctx *gin.Context, info domain.User) error
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (ur *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, ur.toDaoUser(u))
}

func (ur *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return ur.toDomainUser(u), nil
}

// 因为 dao 层的 user 与 domain 层的 user 字段名不一致，所以需要转换
func (ur *CacheUserRepository) toDomainUser(u dao.User) domain.User {
	return domain.User{
		Id:           u.Id,
		Email:        u.Email.String,
		Password:     u.Password,
		NickName:     u.NickName,
		Birthday:     time.UnixMilli(u.Birthday),
		Introduction: u.Introduction,
	}
}

func (ur *CacheUserRepository) toDaoUser(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != ""},
		Password:     u.Password,
		NickName:     u.NickName,
		Birthday:     u.Birthday.UnixMilli(),
		Introduction: u.Introduction,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != ""},
	}
}

func (ur *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 优先从缓存中获取
	user, err := ur.cache.Get(ctx, id)
	// 如果缓存中有， 直接返回
	if err == nil {
		return user, nil
	}
	// 如果缓存中没有， 从数据库中获取
	//if err == cache.ErrCacheMiss {
	//	// 打印日志等等
	//}

	userEntity, err := ur.dao.FindById(ctx, id)
	u := ur.toDomainUser(userEntity)
	_ = ur.cache.Set(ctx, u)
	// 忽略错误，可能是 redis 服务挂了
	return u, nil
}

func (ur *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := ur.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return ur.toDomainUser(user), nil
}

func (ur *CacheUserRepository) Edit(ctx *gin.Context, info domain.User) error {
	return ur.dao.UpdateById(ctx, ur.toDaoUser(info))
}
