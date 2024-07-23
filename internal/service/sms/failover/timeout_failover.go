package failover

import (
	"context"
	"sync/atomic"
	"we_book/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	// 服务商列表
	svcs []sms.Service
	idx  int32
	// 连续超时的个数
	cnt int32
	// 阈值
	threshold int32
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		// 超过阈值，切换服务商
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 切换成功
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		// 状态被打断
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		return err
	}
}
