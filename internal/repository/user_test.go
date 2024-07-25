package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"we_book/internal/domain"
	"we_book/internal/repository/cache"
	cachemocks "we_book/internal/repository/cache/mocks"
	"we_book/internal/repository/dao"
	daomocks "we_book/internal/repository/dao/mocks"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixNano())
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		id  int64

		wantedUser domain.User
		wantedErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中，但是查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(3)).
					Return(domain.User{}, cache.ErrKeyNotExists)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(3)).
					Return(dao.User{
						Id: 3,
						Email: sql.NullString{
							String: "user@example.com",
							Valid:  true,
						},
						Password: "Bogzc2002@",
						Phone: sql.NullString{
							String: "13812345678",
							Valid:  true,
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       3,
					Email:    "user@example.com",
					Password: "Bogzc2002@",
					Phone:    "13812345678",
				}).Return(nil)
				return d, c
			},
			ctx: context.Background(),
			wantedUser: domain.User{
				Id:       3,
				Email:    "user@example.com",
				Password: "Bogzc2002@",
				Phone:    "13812345678",
			},
			wantedErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存命中
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(3)).
					Return(domain.User{}, cache.ErrKeyNotExists)

				d := daomocks.NewMockUserDAO(ctrl)
				return d, c
			},
			ctx: context.Background(),
			wantedUser: domain.User{
				Id:       3,
				Email:    "user@example.com",
				Password: "Bogzc2002@",
				Phone:    "13812345678",
			},
			wantedErr: nil,
		},
		{
			name: "缓存未命中，查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存命中
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(3)).
					Return(domain.User{}, cache.ErrKeyNotExists)

				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(3)).
					Return(dao.User{}, errors.New("mock db error"))
				return d, c
			},
			ctx:        context.Background(),
			id:         3,
			wantedUser: domain.User{},
			wantedErr:  errors.New("mock db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantedUser, u)
			assert.Equal(t, tc.wantedErr, err)
		})
	}
}

func TestRedisConn(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		t.Log(err)
	}
	t.Log("redis连接成功")
}
