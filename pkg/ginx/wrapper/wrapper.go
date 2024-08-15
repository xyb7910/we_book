package wrapper

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	logger2 "we_book/pkg/logger"
)

var logger logger2.V1

type Result struct {
	// 业务代码错误
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func WrapBody[T any](l logger2.V1, fn func(ctx *gin.Context, req T) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}

		res, err := fn(ctx, req)
		if err != nil {
			l.Error("处理业务出错",
				logger.String("path", ctx.Request.URL.Path))
			logger.String("method", ctx.Request.Method)
			logger.String("route", ctx.FullPath())
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapToken[C jwt.Claims](fn func(ctx *gin.Context, claims C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("claims")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, Result{
				Code: 5, Msg: "token 无效", Data: nil,
			})
		}
		c, ok := val.(C)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, Result{
				Code: 5, Msg: "token 错误", Data: nil,
			})
		}

		res, err := fn(ctx, c)
		if err != nil {
			logger.Error("处理业务出错",
				logger.String("path", ctx.Request.URL.Path))
			logger.String("method", ctx.Request.Method)
			logger.String("route", ctx.FullPath())
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WarpBodyANDToken[T any, C jwt.Claims](fn func(ctx *gin.Context, req T, claims C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req T
		if err := ctx.Bind(&req); err != nil {
			return
		}

		val, ok := ctx.Get("claims")
		if !ok {
			ctx.JSON(http.StatusUnauthorized, Result{
				Code: 5, Msg: "token 无效", Data: nil,
			})
		}
		c, ok := val.(C)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, Result{
				Code: 5, Msg: "token 错误", Data: nil,
			})
		}
		res, err := fn(ctx, req, c)
		if err != nil {
			logger.Error("处理业务出错",
				logger.String("path", ctx.Request.URL.Path))
			logger.String("method", ctx.Request.Method)
			logger.String("route", ctx.FullPath())
		}
		ctx.JSON(http.StatusOK, res)
	}
}
