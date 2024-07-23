package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"we_book/internal/pkg/ratelimit"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
	limiter  ratelimit.Limiter
}

func NewService(client *sms.Client, appId, signName *string, limiter ratelimit.Limiter) *Service {
	return &Service{
		appId:    appId,
		signName: signName,
		client:   client,
		limiter:  limiter,
	}
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers []string) error {
	req := sms.NewSendSmsRequest()
	req.SenderId = s.appId
	req.SignName = s.signName
	req.TemplateId = &tplId
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("send sms failed, code: %s, message: %s", *(status.Code), *(status.Message))
		}
	}
	return nil
}
