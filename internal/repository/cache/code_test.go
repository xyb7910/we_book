package cache

import (
	"context"
	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/mock/gomock"
	"testing"
	redismocks "we_book/internal/repository/cache/reidsmocks"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string

		mock func(ctrl *gomock.Controller) redis.Cmdable

		// input
		ctx   context.Context
		biz   string
		phone string
		code  string
		// want error
		wantedErr error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_cod:login:1522"},
					[]string{"123456"}).Return(res)
				return cmd
			},
			ctx:       context.Background(),
			biz:       "login",
			phone:     "1522",
			code:      "123456",
			wantedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			rc := NewRedisCodeCache(tc.mock(ctrl))
			err := rc.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, err, tc.wantedErr)
		})
	}
}
