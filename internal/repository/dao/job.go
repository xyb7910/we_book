package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	UpdateNextUtime(ctx context.Context, id int64, next time.Time) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (G *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := G.db.WithContext(ctx)
	for {
		now := time.Now()
		var j Job
		err := db.WithContext(ctx).Where("status = ? AND next_time <= ?", jobStatusWaiting, now).Error
		if err != nil {
			return Job{}, err
		}
		// 乐观锁， CAS 代替 for update
		res := db.Where("id = ? AND version = ?", j.Id, j.Version).Model(&Job{}).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"utime":   now,
				"version": j.Version + 1,
			})
		if res.Error != nil {
			return Job{}, err
		}
		if res.RowsAffected == 0 {
			continue
		}
		return j, nil
	}
}

func (G *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return G.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error

}

func (G *GORMJobDAO) Release(ctx context.Context, id int64) error {
	return G.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (G *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return G.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusPaused,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (G *GORMJobDAO) UpdateNextUtime(ctx context.Context, id int64, next time.Time) error {
	return G.db.WithContext(ctx).Where("id = ?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

//func NewGORMJobDAO(db *gorm.DB) JobDAO {
//	return &*GORMJobDAO{db: db}
//}

type Job struct {
	Id       int64 `gorm:"primary_key,autoIncrement"`
	Cfg      string
	Executor string
	Name     string `gorm:"unique"`
	Status   int
	NextTime int64 `gorm:"index"`
	Cron     string
	Version  int
	Ctime    int64
	Utime    int64
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusPaused
)
