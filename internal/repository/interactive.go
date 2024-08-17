package repository

import (
	"context"
	"we_book/internal/domain"
	"we_book/internal/repository/cache"
	"we_book/internal/repository/dao"
	"we_book/pkg/logger"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context,
		biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, ids []int64, bizId []string) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	AddRecord(ctx context.Context, aid int64, uid int64) error
}

type CacheReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     logger.V1
}

func (c *CacheReadCntRepository) AddRecord(ctx context.Context, aid int64, uid int64) error {
	// TODO:implement me
	panic("implement me")
}

func (c *CacheReadCntRepository) BatchIncrReadCnt(ctx context.Context, ids []int64, bizId []string) error {
	err := c.dao.BatchIncrReadCnt(ctx, ids, bizId)
	if err != nil {
		return err
	}
	// 这里要批处理 redis
	//return c.cache.BatchIncrReadCntIfPresent(ctx, ids, bizId)
	return nil
}

func (c *CacheReadCntRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CacheReadCntRepository) IncrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CacheReadCntRepository) DecrLike(ctx context.Context, biz string, bizId, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CacheReadCntRepository) AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		BizId: bizId,
		Biz:   biz,
		Cid:   cid,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (c *CacheReadCntRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizId)
	if err == nil {
		return intr, nil
	}
	intro, err := c.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interactive{}, err
	}
	res := c.toDomain(intro)
	go func() {
		err = c.cache.Set(ctx, biz, bizId, res)
		if err != nil {
			c.l.Error("set cache failed",
				logger.Int64("biz_id", bizId),
				logger.String("biz", biz),
				logger.Error(err))
		}
	}()
	return res, nil
}

func (c *CacheReadCntRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CacheReadCntRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CacheReadCntRepository) toDomain(intro dao.Interactive) domain.Interactive {
	return domain.Interactive{
		LikedCnt:   intro.LikeCnt,
		ReadCnt:    intro.ReadCnt,
		CollectCnt: intro.CollectCnt,
	}
}

func NewCacheInteractiveRepository(cache cache.InteractiveCache, dao dao.InteractiveDAO, l logger.V1) InteractiveRepository {
	return &CacheReadCntRepository{
		cache: cache,
		dao:   dao,
		l:     l,
	}
}
