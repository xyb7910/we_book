package sms

import "context"

type Service interface {
	Send(ctx context.Context, biz string, args []string, number ...string) error
}
