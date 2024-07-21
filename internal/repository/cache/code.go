package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany = errors.New("send code too many")
	ErrCodeVerifyError = errors.New("code verify error")
	ErrUnknownError    = errors.New("unknown error")
)

// 编译器在编译的时候，把 set_code 中的代码放进  luaSetCode 这个变量中
//
//go:embed lua/set_code.lua
var luaSetCode string

// 把 verify_code 中的代码放进 luaVerifyCode 这个变量中
//
//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (cc *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (cc *CodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := cc.client.Eval(ctx, luaSetCode, []string{cc.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	case -2:
		return ErrUnknownError
	default:
		return errors.New("系统错误")
	}
}

func (cc *CodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := cc.client.Eval(ctx, luaVerifyCode, []string{cc.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeSendTooMany
	case -2:
		return false, nil
	}
	return false, ErrUnknownError
}
