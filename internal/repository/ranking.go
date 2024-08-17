package repository

import (
	"context"
	"we_book/internal/domain"
	"we_book/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}

type CachedRankingRepository struct {
	redis *cache.RankingRedisCache
	local *cache.RankingLocalCache
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	// 首先从本地缓存中拿去
	_ = c.local.Set(ctx, arts)
	return c.redis.Set(ctx, arts)
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	// 首先从本地缓存中拿去
	data, err := c.local.Get(ctx)
	if err != nil {
		return nil, nil
	}
	data, err = c.redis.Get(ctx)
	if err != nil {
		return c.local.ForceGet(ctx)
	} else {
		_ = c.local.Set(ctx, data)
	}
	return data, nil
}

func NewRankingRepository(redis *cache.RankingRedisCache, local *cache.RankingLocalCache) RankingRepository {
	return &CachedRankingRepository{
		redis: redis,
		local: local,
	}
}
