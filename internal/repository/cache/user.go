package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"we_book/internal/domain"
)

var ErrKeyNotExists = redis.Nil

type UserCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: 10 * time.Minute,
	}
}

func (uc *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

func (uc *RedisUserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := uc.key(id)
	// 如果缓存中有，直接返回
	val, err := uc.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(val, &user)
	return user, err
}

func (uc *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	key := uc.key(user.Id)
	return uc.client.Set(ctx, key, val, uc.expiration).Err()
}
