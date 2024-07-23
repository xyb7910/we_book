package retryable

import (
	"context"
	"errors"
	"we_book/internal/service/sms"
)

type Service struct {
	svc sms.Service
	// 重试
	retryMax int
}

func (s Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tpl, args, numbers...)
	cnt := 1
	for err != nil && cnt < s.retryMax {
		err = s.svc.Send(ctx, tpl, args, numbers...)
		if err != nil {
			return err
		}
		cnt++
	}
	return errors.New("send sms failed twice")
}
