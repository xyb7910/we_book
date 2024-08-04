package web

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"we_book/internal/domain"
	"we_book/internal/service"
	ijwt "we_book/internal/web/jwt"
	"we_book/pkg/ginx/wrapper"
	logger2 "we_book/pkg/logger"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger2.V1
}

func NewArticleHandler(svc service.ArticleService, l logger2.V1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (at *ArticleHandler) RegisterRouters(server *gin.Engine) {
	article := server.Group("/articles")

	article.POST("/edit", at.Edit)
	article.POST("/publish", at.Publish)
	article.POST("/withdraw", at.Withdraw)
	article.POST("/list",
		wrapper.WarpBodyANDToken[ListReq, ijwt.UserClaims](at.List))
	article.GET("/detail/:id",
		wrapper.WrapToken[ijwt.UserClaims](at.Detail))

	pub := article.Group("/pub")
	pub.GET("/:id", at.PubDetail)
}

func (at *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "id error",
		})
		return
	}
	art, err := at.svc.GetPubById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "system error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "success",
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Content: art.Content,
			Status:  uint8(art.Status),
			Author:  art.Author.Name,
			Ctime:   art.Ctime.Format("2006-01-02 15:04:05"),
			Utime:   art.Utime.Format("2006-01-02 15:04:05"),
		},
	})
}

func (at *ArticleHandler) Detail(ctx *gin.Context, claims ijwt.UserClaims) (Result, error) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return Result{
			Code: 4,
			Msg:  "id error",
		}, err
	}
	art, err := at.svc.GetById(ctx, id)
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "system error",
		}, err
	}
	if art.Author.Id != claims.Uid {
		// 此时用户为非法攻击用户
		// 应该报警
		return Result{
			Code: 4,
			Msg:  "not your article",
		}, err
	}
	return Result{
		Code: 2,
		Msg:  "success",
		Data: ArticleVO{
			Id:      art.Id,
			Title:   art.Title,
			Content: art.Content,
			Status:  uint8(art.Status),
			Ctime:   art.Ctime.Format("2006-01-02 15:04:05"),
			Utime:   art.Utime.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

func (at *ArticleHandler) List(ctx *gin.Context, req ListReq, uc ijwt.UserClaims) (Result, error) {
	res, err := at.svc.List(ctx, uc.Uid, req.OffSet, req.Limit)
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "system error",
		}, err
	}
	return Result{
		Code: 2,
		Msg:  "success",
		Data: slice.Map[domain.Article, ArticleVO](res,
			func(idx int, src domain.Article) ArticleVO {
				return ArticleVO{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract(),
					Ctime:    src.Ctime.Format("2006-01-02 15:04:05"),
					Utime:    src.Utime.Format("2006-01-02 15:04:05"),
				}
			}),
	}, nil
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
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "Bind error",
		})
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

func (at *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
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
	err := at.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "withdraw article error",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "success",
	})
}
