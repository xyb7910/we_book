package service

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	"math"
	"time"
	"we_book/internal/domain"
	"we_book/internal/repository"
)

type RankingService interface {
	TopN(ctx context.Context) error
	topN(ctx context.Context) ([]domain.Article, error)
}

type BatchRankingService struct {
	artSvc    ArticleService
	interSvc  InteractiveService
	repo      repository.RankingRepository
	batchSize int
	n         int
	scoreFunc func(t time.Time, likeCnt int64) float64
}

func (b *BatchRankingService) TopN(ctx context.Context) error {
	res, err := b.topN(ctx)
	if err != nil {
		return err
	}
	//log.Println(res)
	return b.repo.ReplaceTopN(ctx, res)
}

func (b *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	type Score struct {
		art   domain.Article
		score float64
	}
	now := time.Now()
	offset := 0
	topN := queue.NewConcurrentPriorityQueue[Score](b.n,
		func(a, b Score) int {
			if a.score == b.score {
				return 0
			}
			if a.score > b.score {
				return 1
			}
			return -1
		})

	for {
		arts, err := b.artSvc.ListPub(ctx, now, offset, b.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts,
			func(idx int, src domain.Article) int64 {
				return src.Id
			})

		inter, err := b.interSvc.GetByIds(ctx, "article", ids)
		if err != nil {
			return nil, err
		}

		for _, art := range arts {
			intr, ok := inter[art.Id]
			if !ok {
				continue
			}
			score := b.scoreFunc(art.Ctime, intr.LikedCnt)

			err = topN.Enqueue(Score{
				art:   art,
				score: score,
			})

			if errors.Is(err, queue.ErrOutOfCapacity) {
				val, _ := topN.Dequeue()
				if val.score < score {
					_ = topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				} else {
					_ = topN.Enqueue(val)
				}
			}
		}

		if len(arts) < b.batchSize || now.Sub(arts[len(arts)-1].Ctime).Hours() > 7*24 {
			break
		}
		offset = offset + len(arts)
	}

	res := make([]domain.Article, b.n)
	for i := b.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			break
		}
		res[i] = val.art
	}
	return res, nil
}

func NewBatchRankingService(artSvc ArticleService, interSvc InteractiveService) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		interSvc:  interSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			sec := time.Since(t).Seconds()
			return float64(likeCnt-1) / math.Pow(float64(sec+2), 1.5)
		},
	}
}
