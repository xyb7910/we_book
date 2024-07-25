package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisJWTHandler struct {
	cmd redis.Cmdable
}

func NewRedisJWTHandler(cmd redis.Cmdable) *RedisJWTHandler {
	return &RedisJWTHandler{
		cmd: cmd,
	}
}

func (rj *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New()
	err := rj.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
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
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD4ogvf0"))
	if err != nil {
		//ctx.String(http.StatusInternalServerError, "internal server error, Login failed")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
