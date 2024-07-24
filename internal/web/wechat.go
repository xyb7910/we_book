package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"we_book/internal/service"
	"we_book/internal/service/oauth2/wechat"
)

type OAuth2WeChatHandler struct {
	svc     wechat.Service
	userSvc service.UserService
	jwtHandler
}

func NewOAuth2WeChatHandler(svc wechat.Service, userSvc service.UserService) *OAuth2WeChatHandler {
	return &OAuth2WeChatHandler{
		svc:     svc,
		userSvc: userSvc,
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

func (h *OAuth2WeChatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	info, err := h.svc.VerifyCode(ctx, code, state)
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
	err = h.setJWTToken(ctx, u.Id)
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
