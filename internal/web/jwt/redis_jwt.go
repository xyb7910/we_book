package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	AtKey = []byte("95osj3fUD8fo0mlYdDbncXz4VD4ogvf4")
	RtKey = []byte("95osj6fUD8fo0mlYdDbncXz4VD4igvf7")
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) Handler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (rj *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := rj.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = rj.SetRefreshToken(ctx, uid, ssid)
	return err
}

func (rj *RedisJWTHandler) SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 过期时间我们设置为 7 days
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		//ctx.String(http.StatusInternalServerError, "internal server error, Login failed")
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (rj *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		//ctx.String(http.StatusInternalServerError, "internal server error, Login failed")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (rj *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenStr := ctx.GetHeader("Authorization")
	segs := strings.Split(tokenStr, " ")
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}

func (rj *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	val, err := rj.cmd.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if val == 0 {
			return nil
		}
		return errors.New("session expired")
	default:
		return err
	}
}

func (rj *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")

	claims := ctx.MustGet("claims").(*UserClaims)
	return rj.cmd.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}
