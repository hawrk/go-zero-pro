package logic

import (
	"algo_assess/assess-mq-server/internal/consumer"
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendChildOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendChildOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendChildOrderLogic {
	return &SendChildOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SendChildOrder api测试接口子单
func (l *SendChildOrderLogic) SendChildOrder(in *proto.ApiAlgoOrderReq) (*proto.ApiAlgoOrderRsp, error) {
	s := consumer.NewAlgoPlatformOrderTrade(l.ctx, l.svcCtx)
	_ = s.Consume("key", string(in.GetValue()))

	return &proto.ApiAlgoOrderRsp{
		Code: 200,
		Msg:  "success",
	}, nil
}
