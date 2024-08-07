package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
	"net/url"
	"we_book/internal/domain"
	logger2 "we_book/pkg/logger"
)

var redirectUrl = url.PathEscape("https://test.com/oauth2/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
	logger    logger2.V1
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenID  string `json:"openid"`
	Scope   string `json:"scope"`
	UnionID string `json:"unionid"`
}

func NewService(appId string, appSecret string, l logger2.V1) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
		logger:    l,
	}
}

func (ws *service) AuthURL(ctx context.Context) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"
	state := uuid.New()
	return fmt.Sprintf(urlPattern, ws.appId, redirectUrl, state), nil
}

func (ws *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, ws.appId, ws.appSecret, code)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := ws.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	decoder := json.NewDecoder(resp.Body)

	var result Result
	err = decoder.Decode(&result)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if result.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("wechat auth error: %s", result.ErrMsg)
	}
	return domain.WechatInfo{
		OpenId:  result.OpenID,
		UnionId: result.UnionID,
	}, nil
}
