package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
	"we_book/internal/domain"
)

type RankingRedisCache struct {
	client redis.Cmdable
	key    string
}

func (r *RankingRedisCache) Set(ctx context.Context, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = ""
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, val, time.Minute*10).Err()
}

func (r *RankingRedisCache) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(data, &res)
	return res, nil
}

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

func NewRankingRedisCache(client redis.Cmdable, key string) RankingCache {
	return &RankingRedisCache{
		client: client,
		key:    "ranking",
	}
}
