package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"time"
	"we_book/internal/domain"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (ur *UserRepository) Create(ctx context.Context, u domain.User) error {
	return ur.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (ur *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := ur.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

// 因为 dao 层的 user 与 domain 层的 user 字段名不一致，所以需要转换
func (ur *UserRepository) toDomainUser(u dao.User) domain.User {
	return domain.User{
		Id:           u.Id,
		Email:        u.Email,
		Password:     u.Password,
		NickName:     u.NickName,
		Birthday:     time.UnixMilli(u.Birthday),
		Introduction: u.Introduction,
	}
}

func (ur *UserRepository) toDaoUser(u domain.User) dao.User {
	return dao.User{
		Id:           u.Id,
		Email:        u.Email,
		Password:     u.Password,
		NickName:     u.NickName,
		Birthday:     u.Birthday.UnixMilli(),
		Introduction: u.Introduction,
	}
}

func (ur *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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
	err = ur.cache.Set(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (ur *UserRepository) Edit(ctx *gin.Context, info domain.User) error {
	return ur.dao.UpdateById(ctx, ur.toDaoUser(info))
}
