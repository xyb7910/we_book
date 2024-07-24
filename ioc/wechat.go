package ioc

import "we_book/internal/service/oauth2/wechat"

func InitWechatService() wechat.Service {
	// 更改为自己的appId和appSecret
	appId := "wx1234567890abcdef"
	appSecret := "1234567890abcdef"
	return wechat.NewService(appId, appSecret)
}
