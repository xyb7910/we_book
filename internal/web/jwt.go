package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtHandler struct{}

func (j *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		//ctx.String(http.StatusInternalServerError, "internal server error, Login failed")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

// UserClaims jwt token 携带的信息
type UserClaims struct {
	jwt.RegisteredClaims

	// 声明自己的字段
	Uid       int64 `json:"uid"`
	UserAgent string
}
