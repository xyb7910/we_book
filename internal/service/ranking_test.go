package service

import (
	"context"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"we_book/internal/domain"
	svcmocks "we_book/internal/service/mocks"
)

func TestBatchRankingService_TopN(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name      string
		mock      func(ctrl *gomock.Controller) (ArticleService, InteractiveService)
		wantedErr error
		wantedRes []domain.Article
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) (ArticleService, InteractiveService) {
				artSvc := svcmocks.NewMockArticleService(ctrl)
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 0, 3).
					Return([]domain.Article{
						{Id: 1, Ctime: now, Utime: now, Title: "Java"},
						{Id: 2, Ctime: now, Utime: now, Title: "Go"},
						{Id: 3, Ctime: now, Utime: now, Title: "Python"},
					}, nil)
				artSvc.EXPECT().ListPub(gomock.Any(), gomock.Any(), 3, 3).
					Return([]domain.Article{}, nil)

				interSvc := svcmocks.NewMockInteractiveService(ctrl)
				interSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{1, 2, 3}).
					Return(map[int64]domain.Interactive{
						1: {BizId: 1, LikedCnt: 1},
						2: {BizId: 2, LikedCnt: 2},
						3: {BizId: 3, LikedCnt: 3},
					}, nil)
				interSvc.EXPECT().GetByIds(gomock.Any(), "article", []int64{}).
					Return(map[int64]domain.Interactive{}, nil)
				return artSvc, interSvc
			},
			wantedRes: []domain.Article{
				{Id: 3, Ctime: now, Utime: now, Title: "Python"},
				{Id: 2, Ctime: now, Utime: now, Title: "Go"},
				{Id: 1, Ctime: now, Utime: now, Title: "Java"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			artSvc, interSvc := tc.mock(ctrl)
			svc := NewBatchRankingService(artSvc, interSvc)
			svc.batchSize = 3
			svc.n = 3
			svc.scoreFunc = func(t time.Time, likeCnt int64) float64 {
				return float64(likeCnt)
			}
			arts, err := svc.topN(context.Background())
			assert.Equal(t, tc.wantedErr, err)
			assert.Equal(t, tc.wantedRes, arts)
		})
	}
}

//[{3   {0 } 0 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834}
// {2   {0 } 0 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834}
// {1   {0 } 0 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834}]
//does not equal
//[{0   {0 } 0 0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC}
//{2   {0 } 0 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834}
//{1   {0 } 0 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834 2024-08-16 23:01:37.208097 +0800 CST m=+0.008096834}]
