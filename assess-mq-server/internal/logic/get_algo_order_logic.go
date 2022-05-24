package logic

import (
	"algo_assess/global"
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoOrderLogic {
	return &GetAlgoOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  -----------下面的接口为调试功能，不对外开放
func (l *GetAlgoOrderLogic) GetAlgoOrder(in *proto.AlgoOrderReq) (*proto.AlgoOrderRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("get algo order Id :", in.GetAlgoOrderId())

	global.GlobalAlgoOrder.RWMutex.RLock()
	out := &proto.AlgoOrderRsp{
		Qty: 0,
	}
	global.GlobalAlgoOrder.RWMutex.RUnlock()

	return out, nil
}
