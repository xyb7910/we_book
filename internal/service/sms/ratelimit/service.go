package ratelimit

import (
	"context"
	"fmt"
	"we_book/internal/pkg/ratelimit"
	"we_book/internal/service/sms"
)

var ErrLimited = fmt.Errorf("ratelimit: limited")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent:send")
	if err != nil {
		return fmt.Errorf("ratelimit: %w", err)
	}
	if limited {
		return ErrLimited
	}
	err = s.svc.Send(ctx, tpl, args, number...)
	return err
}
