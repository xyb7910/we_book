package memory

import (
	"context"
	"fmt"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// Send 模拟发送短信
func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	fmt.Println(args)
	return nil
}
