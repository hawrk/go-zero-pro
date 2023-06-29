package logic

import (
	"algo_assess/assess-mq-server/internal/consumer"
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendAlgoOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendAlgoOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendAlgoOrderLogic {
	return &SendAlgoOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  SendAlgoOrder api测试接口用
func (l *SendAlgoOrderLogic) SendAlgoOrder(in *proto.ApiAlgoOrderReq) (*proto.ApiAlgoOrderRsp, error) {
	// todo: add your logic here and delete this line
	s := consumer.NewAlgoOrderTrade(l.ctx, l.svcCtx)
	_ = s.Consume("key", string(in.GetValue()))

	return &proto.ApiAlgoOrderRsp{
		Code: 200,
		Msg:  "success",
	}, nil
}
