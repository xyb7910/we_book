package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"we_book/internal/service"
	"we_book/internal/service/oauth2/wechat"
	ijwt "we_book/internal/web/jwt"
)

type OAuth2WeChatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	ijwt.Handler
	stateKey []byte
	cfg      WechatHandlerConfig
}

type WechatHandlerConfig struct {
	Secure bool
}

func NewOAuth2WeChatHandler(svc wechat.Service, userSvc service.UserService,
	jwtHdl ijwt.Handler, cfg WechatHandlerConfig) *OAuth2WeChatHandler {
	return &OAuth2WeChatHandler{
		svc:      svc,
		userSvc:  userSvc,
		Handler:  jwtHdl,
		stateKey: []byte("95osj3fUD7foxmlYdDbncXz4VD2igvf1"),
		cfg:      cfg,
	}
}

func (h *OAuth2WeChatHandler) RegisterRoutes(r *gin.Engine) {
	we := r.Group("/oauth2/wechat")
	we.GET("/authurl", h.AuthURL)
	we.Any("/callback", h.Callback)
}

func (h *OAuth2WeChatHandler) AuthURL(ctx *gin.Context) {
	url, err := h.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "get auth url failed",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  url,
	})
}

func (h *OAuth2WeChatHandler) SetStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return err
	}
	ctx.SetCookie("jwt-state", tokenStr, 600, "/oauth2/wechat/callback", "",
		h.cfg.Secure, true)
	return nil
}

func (h *OAuth2WeChatHandler) VerifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("get cookie failed: %w", err)
	}

	var sc StateClaims
	token, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("verify state failed: %w", err)
	}
	if state != sc.State {
		return errors.New("state not match")
	}
	return nil
}

//

func (h *OAuth2WeChatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	err := h.VerifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "verify state failed",
		})
	}
	info, err := h.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "verify code failed",
		})
	}
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "system error",
		})
	}
	err = h.SetLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "set token failed",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "success",
	})
}

type StateClaims struct {
	State string
	jwt.RegisteredClaims
}
