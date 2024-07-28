package startup

import (
	"we_book/internal/pkg/logger"
	"we_book/internal/service/oauth2/wechat"
	"we_book/internal/web"
)

func InitWechatService(l logger.V1) wechat.Service {
	// 更改为自己的appId和appSecret
	appId := "wx1234567890abcdef"
	appSecret := "1234567890abcdef"
	return wechat.NewService(appId, appSecret, l)
}

func NewWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
