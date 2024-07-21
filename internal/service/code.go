package service

import (
	"context"
	"fmt"
	"math/rand"
	"we_book/internal/repository"
	"we_book/internal/service/sms"
)

const codeTplId = "1877556"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func (svc *CodeService) genCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}

func (svc *CodeService) Send(ctx context.Context,
	// 区别业务场景
	biz string,
	phone string) error {
	// 生成一个验证码
	code := svc.genCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, code)
}
