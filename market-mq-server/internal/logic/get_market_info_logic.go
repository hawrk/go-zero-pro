package logic

import (
	"context"

	"algo_assess/market-mq-server/internal/svc"
	"algo_assess/market-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMarketInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketInfoLogic {
	return &GetMarketInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取行情价格数量信息
func (l *GetMarketInfoLogic) GetMarketInfo(in *proto.MarkReq) (*proto.MarkRsp, error) {
	// todo: add your logic here and delete this line
	l.Infof("get req:%+v", in)

	return &proto.MarkRsp{}, nil
}
