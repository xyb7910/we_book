package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"we_book/internal/domain"
)

var (
	luaIncrReadCnt string
)

const (
	fileReadCnt    = "read_cnt"
	fileCollectCnt = "collect_cnt"
	fileLikeCnt    = "like_cnt"
)

type InteractiveCache interface {
	// IncrReadCntIfPresent 如果在缓存中有对应的数据，就 +1
	IncrReadCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	IncrLikeCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context,
		biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	// Get 查询缓存中数据
	// 事实上，这里 liked 和 collected 是不需要缓存的
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
}

type RedisInteractiveCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func (r *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}

func (r *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrReadCnt,
		[]string{r.key(biz, bizId)},
		fileReadCnt, 1).Err()
}

func (r *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrReadCnt,
		[]string{r.key(biz, bizId)},
		fileLikeCnt, 1).Err()
}

func (r *RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrReadCnt,
		[]string{r.key(biz, bizId)},
		fileLikeCnt, -1).Err()
}

func (r *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return r.client.Eval(ctx, luaIncrReadCnt,
		[]string{r.key(biz, bizId)},
		fileCollectCnt, 1).Err()
}

func (r *RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	data, err := r.client.HGetAll(ctx, r.key(biz, bizId)).Result()
	if err != nil {
		return domain.Interactive{}, nil
	}
	if len(data) == 0 {
		return domain.Interactive{}, ErrKeyNotExists
	}

	collectCnt, _ := strconv.ParseInt(data[fileCollectCnt], 10, 64)
	likeCnt, _ := strconv.ParseInt(data[fileLikeCnt], 10, 64)
	readCnt, _ := strconv.ParseInt(data[fileReadCnt], 10, 64)

	return domain.Interactive{
		CollectCnt: collectCnt,
		LikedCnt:   likeCnt,
		ReadCnt:    readCnt,
	}, err
}

func (r *RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
	key := r.key(biz, bizId)
	err := r.client.HMSet(ctx, key,
		fileReadCnt, intr.ReadCnt,
		fileLikeCnt, intr.LikedCnt,
		fileCollectCnt, intr.CollectCnt).Err()
	if err != nil {
		return err
	}
	return r.client.Expire(ctx, key, r.expiration).Err()
}

func NewRedisInteractiveCache(client redis.Cmdable, expiration time.Duration) InteractiveCache {
	return &RedisInteractiveCache{
		client:     client,
		expiration: expiration,
	}
}
