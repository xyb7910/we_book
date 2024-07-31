package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"we_book/internal/domain"
	"we_book/internal/pkg/logger"
	"we_book/internal/service"
	ijwt "we_book/internal/web/jwt"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.V1
}

func NewArticleHandler(svc service.ArticleService, l logger.V1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (at *ArticleHandler) RegisterRouters(server *gin.Engine) {
	article := server.Group("/articles")

	article.POST("/edit", at.Edit)
	article.POST("/publish", at.Publish)
}

func (at *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 获取作者的id
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(200, Result{
			Code: 5,
			Msg:  "token error",
		})
		at.l.Error("not find user session")
		return
	}
	// 校验部分
	id, err := at.svc.Edit(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(200, Result{
			Code: 5,
			Msg:  "system error",
		})
		at.l.Error("save article error")
		return
	}
	ctx.JSON(200, Result{
		Msg:  "success",
		Data: id,
	})
}

func (at *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "system error",
		})
		at.l.Error("not find user session")
		return
	}
	println(claims.Uid)
	id, err := at.svc.Publish(ctx, req.toDomain(claims.Uid))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "publish article error",
		})
		at.l.Error("save article error")
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "success",
		Data: id,
	})
}

type ArticleReq struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (req ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
