package dao

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type InteractiveDAO interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error)
	DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (Interactive, error)
	InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error
	GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error)
	BatchIncrReadCnt(ctx context.Context, ids []int64, bizs []string) error
}

type GORMInteractiveDAO struct {
	db *gorm.DB
}

func (G *GORMInteractiveDAO) BatchIncrReadCnt(ctx context.Context, ids []int64, bizs []string) error {
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMInteractiveDAO(tx)
		for i := range bizs {
			err := txDAO.IncrReadCnt(ctx, bizs[i], ids[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (G *GORMInteractiveDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()
	return G.db.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"read_cnt": gorm.Expr("read_cnt + 1"),
			"utime":    now,
		}),
	}).Create(&Interactive{
		BizId:   bizId,
		Biz:     biz,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

func (G *GORMInteractiveDAO) InsertLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()
	// 同时记录点赞和点赞数
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"utime":  now,
				"status": 1,
			}),
		}).Create(&UserLikeBiz{
			Biz:    biz,
			BizId:  bizId,
			Uid:    uid,
			Status: 1,
			Ctime:  now,
			Utime:  now,
		}).Error
		if err != nil {
			return err
		}

		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"like_cnt": gorm.Expr("like_cnt + 1"),
				"utime":    now,
			}),
		}).Error
	})
}

func (G *GORMInteractiveDAO) DeleteLikeInfo(ctx context.Context, biz string, bizId, uid int64) error {
	now := time.Now().UnixMilli()
	// 同时记录点赞和点赞数
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&UserLikeBiz{}).
			Where("biz = ? and biz_id = ? and uid = ?", biz, bizId, uid).
			Updates(map[string]any{
				"utime":  now,
				"status": 0,
			}).Error
		if err != nil {
			return err
		}

		return tx.Model(&Interactive{}).
			Where("biz = ? and biz_id = ?", biz, bizId).
			Updates(map[string]any{
				"utime":    now,
				"like_cnt": gorm.Expr("like_cnt - 1"),
			}).Error
	})
}

func (G *GORMInteractiveDAO) Get(ctx context.Context, biz string, bizId int64) (Interactive, error) {
	var res Interactive
	err := G.db.WithContext(ctx).Where("biz = ? and biz_id = ?", biz, bizId).First(&res).Error
	return res, err
}

func (G *GORMInteractiveDAO) InsertCollectionBiz(ctx context.Context, cb UserCollectionBiz) error {
	now := time.Now().UnixMilli()
	cb.Ctime = now
	cb.Utime = now
	return G.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 首先新建一个收藏夹
		err := G.db.WithContext(ctx).Create(&cb).Error
		if err != nil {
			return err
		}
		// 更新收藏的数量
		return tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]any{
				"collect_cnt": gorm.Expr("collect_cnt + 1"),
				"utime":       now,
			}),
		}).Create(&Interactive{
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
			Biz:        cb.Biz,
			BizId:      cb.BizId,
		}).Error
	})
}

func (G *GORMInteractiveDAO) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (UserLikeBiz, error) {
	var res UserLikeBiz
	err := G.db.WithContext(ctx).Where("biz = ? and biz_id = ? and uid = ? and status = ?", biz, bizId, uid, 1).First(&res).Error
	return res, err
}

func (G *GORMInteractiveDAO) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (UserCollectionBiz, error) {
	var res UserCollectionBiz
	err := G.db.WithContext(ctx).Where("biz = ? and biz_id = ? and uid = ?", biz, bizId, uid).First(&res).Error
	return res, err
}

func NewGORMInteractiveDAO(db *gorm.DB) InteractiveDAO {
	return &GORMInteractiveDAO{db: db}
}

type UserLikeBiz struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	BizId  int64  `gorm:"uniqueIndex:idx_biz_id"`
	Biz    string `gorm:"uniqueIndex:idx_biz_type;type:varchar(128)"`
	Uid    int64  `gorm:"uniqueIndex:idx_uid"`
	Ctime  int64
	Utime  int64
	Status int64
}

type UserCollectionBiz struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Cid   int64  `gorm:"uniqueIndex:idx_cid"`
	BizId int64  `gorm:"uniqueIndex:idx_biz_id"`
	Biz   string `gorm:"uniqueIndex:idx_biz_type;type:varchar(128)"`
	Uid   int64  `gorm:"uniqueIndex:idx_uid"`
	Ctime int64
	Utime int64
}

type Collection struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Name  string `gorm:"type:varchar(1024)"`
	Uid   int64  `gorm:"uniqueIndex:idx_uid"`
	Ctime int64
	Utime int64
}

type CollectionItem struct {
	Cid   int64
	Cname string
	BizId int64
	Biz   string
}

type Interactive struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	BizId      int64  `gorm:"uniqueIndex:idx_biz_id"`
	Biz        string `gorm:"uniqueIndex:idx_biz_type;type:varchar(128)"`
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Ctime      int64
	Utime      int64
}
