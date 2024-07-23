package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"we_book/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	key string
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}

// Send 安全发送方法， 其中biz 为业务标识， 必须是线下申请的业务标识， 否则会被拦截
func (s *SMSService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	var tc Claims
	// 如果可以解析成功，说明就是对应的业务标识， 否则就是非法的业务标识
	token, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}
