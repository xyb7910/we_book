package repository

import (
	"context"
	"we_book/internal/repository/cache"
)

var (
	ErrCodeTooMany     = cache.ErrCodeSendTooMany
	ErrCodeVerifyError = cache.ErrCodeVerifyError
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (cr *CodeRepository) Store(ctx context.Context, biz string,
	phone string, code string) error {
	return cr.cache.Set(ctx, biz, phone, code)
}

func (cr *CodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return cr.cache.Verify(ctx, biz, phone, inputCode)
}
