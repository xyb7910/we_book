package web

import (
	"github.com/gin-gonic/gin"
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
}

func (at *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Req
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
	id, err := at.svc.Edit(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
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
