package logic

import (
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullShMarketDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPullShMarketDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullShMarketDataLogic {
	return &PullShMarketDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  推送上交所行情数据
func (l *PullShMarketDataLogic) PullShMarketData(in *proto.MarketDataReq) (*proto.MarketDataReq, error) {
	// todo: add your logic here and delete this line
	// 接口废弃
	return &proto.MarketDataReq{}, nil
}
