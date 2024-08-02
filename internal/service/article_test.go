package service

import (
	"context"
	"errors"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"testing"
	"we_book/internal/domain"
	"we_book/internal/repository/article"
	articlerepomock "we_book/internal/repository/article/mocks"
)

func Test_articleService_Publish(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository)

		art domain.Article

		wantedErr error
		wantedId  int64
	}{
		{
			name: "新建发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository, article.ArticleReaderRepository) {
				author := articlerepomock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Create(gomock.Any(), domain.Article{
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				reader := articlerepomock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					Id:      1,
					Title:   "title",
					Content: "content",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return author, reader
			},
			art: domain.Article{
				Title:   "title",
				Content: "content",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantedId: 1,
		},

		{
			name: "修改并发表成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository,
				article.ArticleReaderRepository) {
				author := articlerepomock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				reader := articlerepomock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用了制作库 ID
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)
				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantedId: 2,
		},

		{
			name: "保存到制作库失败",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository,
				article.ArticleReaderRepository) {
				author := articlerepomock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(errors.New("mock db error"))
				reader := articlerepomock.NewMockArticleReaderRepository(ctrl)
				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantedId:  0,
			wantedErr: errors.New("mock db error"),
		},
		{
			name: "保存到制作库成到线上库成功",
			mock: func(ctrl *gomock.Controller) (article.ArticleAuthorRepository,
				article.ArticleReaderRepository) {
				author := articlerepomock.NewMockArticleAuthorRepository(ctrl)
				author.EXPECT().Update(gomock.Any(), domain.Article{
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(nil)
				reader := articlerepomock.NewMockArticleReaderRepository(ctrl)
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用了制作库 ID
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(0), errors.New("mock db error"))
				// 重试成功
				reader.EXPECT().Save(gomock.Any(), domain.Article{
					// 确保使用了制作库 ID
					Id:      2,
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(2), nil)
				return author, reader
			},
			art: domain.Article{
				Id:      2,
				Title:   "我的标题",
				Content: "我的内容",
				Author: domain.Author{
					Id: 123,
				},
			},
			wantedId:  2,
			wantedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			author, reader := tc.mock(ctrl)
			svc := NewArticleServiceV1(reader, author)
			id, err := svc.Publish(context.Background(), tc.art)
			assert.Equal(t, tc.wantedErr, err)
			assert.Equal(t, tc.wantedId, id)
		})
	}
}
