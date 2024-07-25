package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	ijwt "we_book/internal/web/jwt"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要进行校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 使用 JWT 进行校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 用户代理不一致
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 实现刷新 token
		//now := time.Now()
		//// 每十秒刷新一次
		//if claims.ExpiresAt.Sub(now) < time.Second*10 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
		//	if err != nil {
		//		log.Println("token 刷新失败")
		//	}
		//	ctx.Header("x-jwt-token", tokenStr)
		//}
		ctx.Set("claims", claims)
	}
}
