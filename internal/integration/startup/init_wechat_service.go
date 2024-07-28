package startup

import (
	"we_book/internal/pkg/logger"
	"we_book/internal/service/oauth2/wechat"
)

func InitPhoneToWechatService(v1 logger.V1) wechat.Service {
	return wechat.NewService("", "", v1)
}
