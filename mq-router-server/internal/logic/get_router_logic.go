package logic

import (
	"algo_assess/mq-router-server/ordproto"
	"context"

	"algo_assess/mq-router-server/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetRouterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRouterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRouterLogic {
	return &GetRouterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  proto
func (l *GetRouterLogic) GetRouter(in *ordproto.RouterReq) (*ordproto.RouterRsp, error) {
	// todo: add your logic here and delete this line

	return &ordproto.RouterRsp{}, nil
}
