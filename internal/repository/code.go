package repository

import (
	"context"
	"we_book/internal/repository/cache"
)

var (
	ErrCodeTooMany     = cache.ErrCodeSendTooMany
	ErrCodeVerifyError = cache.ErrCodeVerifyError
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: c,
	}
}

func (cr *CacheCodeRepository) Store(ctx context.Context, biz string,
	phone string, code string) error {
	return cr.cache.Set(ctx, biz, phone, code)
}

func (cr *CacheCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return cr.cache.Verify(ctx, biz, phone, inputCode)
}
