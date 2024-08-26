package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"we_book/api/proto/gen/intr"
	"we_book/interactive/domain"
	"we_book/interactive/service"
)

type InteractiveServiceServer struct {
	intrv1.UnsafeInteractiveServiceServer
	asv service.InteractiveService
}

func (i *InteractiveServiceServer) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	err := i.asv.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

func (i *InteractiveServiceServer) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	err := i.asv.Like(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.LikeResponse{}, err
}

func (i *InteractiveServiceServer) CancelLike(ctx context.Context, request *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	if request.Uid <= 0 {
		return nil, status.Error(codes.InvalidArgument, "uid must be greater than 0")
	}
	err := i.asv.CancelLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.CancelLikeResponse{}, err
}

func (i *InteractiveServiceServer) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := i.asv.Collect(ctx, request.GetBiz(), request.GetBizId(), request.GetCid(), request.GetUid())
	return &intrv1.CollectResponse{}, err
}
func (i *InteractiveServiceServer) Get(ctx context.Context, request *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	res, err := i.asv.Get(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{Intr: i.toDTO(res)}, nil
}

func (i *InteractiveServiceServer) toDTO(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		LikeCnt:    intr.LikedCnt,
		Liked:      intr.Liked,
		ReadCnt:    intr.ReadCnt,
	}
}

func (i *InteractiveServiceServer) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	res, err := i.asv.GetByIds(ctx, request.GetBiz(), request.GetIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*intrv1.Interactive, len(res))
	for k, v := range res {
		m[k] = i.toDTO(v)
	}
	return &intrv1.GetByIdsResponse{
		Intrs: m,
	}, nil
}
