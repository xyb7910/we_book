package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"we_book/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		if err != nil {
			return err
		}
		log.Println(err)
	}
	return errors.New("all sms service failed")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < length+idx; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			// 输出日志
		}
	}
	return errors.New("all sms service failed")
}
