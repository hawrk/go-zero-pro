package logic

import (
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoOrderSummaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoOrderSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoOrderSummaryLogic {
	return &GetAlgoOrderSummaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAlgoOrderSummary 根据算法类型取订单汇总信息 （dashboard 算法列表）
func (l *GetAlgoOrderSummaryLogic) GetAlgoOrderSummary(in *proto.AlgoOrderReq) (*proto.AlgoOrderRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in GetAlgoOrderSummary get req:", in)
	result, count, err := l.svcCtx.SummaryRepo.GetAlgoOrderSummary(l.ctx, in.GetDate(), in.GetUserId(), in.GetUserType(),
		in.GetAlgoType(), int(in.GetPage()), int(in.GetLimit()))
	if err != nil {
		l.Logger.Error("query algo order summary error:", err)
		return &proto.AlgoOrderRsp{}, nil
	}
	l.Logger.Info("get result:", result, "count:", count)
	var list []*proto.AlgoTradeInfo
	for _, v := range result {
		tradeSideSum := v.BuyAmount + v.SellAmount
		var buyRate, sellRate float64
		if tradeSideSum > 0 {
			buyRate = float64(v.BuyAmount) / float64(tradeSideSum)
			sellRate = float64(v.SellAmount) / float64(tradeSideSum)
		}
		i := &proto.AlgoTradeInfo{
			Provider:         v.Provider,
			UserNum:          v.UserNum,
			TotalTradeAmount: v.TradeAmount,
			ProfitRate:       v.ProfitRate,
			OrderNum:         v.OrderNum,
			BuyRate:          buyRate,
			SellRate:         sellRate,
		}
		list = append(list, i)
	}
	rsp := &proto.AlgoOrderRsp{
		Code:  200,
		Msg:   "success",
		Total: count,
		Info:  list,
	}
	l.Logger.Info("get rsp:", rsp)

	return rsp, nil
}
