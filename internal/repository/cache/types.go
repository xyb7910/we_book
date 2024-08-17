package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)
}

type LocalCache struct {
}

type RedisCache struct {
}

type DoubleCache struct {
	local Cache
	redis Cache
}

func (d *DoubleCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (d *DoubleCache) Get(ctx context.Context, key string) (any, error) {
	//TODO implement me
	panic("implement me")
}
