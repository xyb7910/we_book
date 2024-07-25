package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
}

// UserClaims jwt token 携带的信息
type UserClaims struct {
	jwt.RegisteredClaims

	// 声明自己的字段
	Uid       int64 `json:"uid"`
	Ssid      string
	UserAgent string
}
