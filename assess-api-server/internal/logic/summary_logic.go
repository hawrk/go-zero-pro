package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SummaryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SummaryLogic {
	return &SummaryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Summary 算法dashboard 里汇总信息
func (l *SummaryLogic) Summary(req *types.AlgoComsumReq) (resp *types.AlgoComsumRsp, err error) {
	l.Logger.Infof("dashboard Summary, get req:%+v", req)

	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	//end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))

	oReq := &assessservice.OrderSummaryReq{
		Date:     start,
		UserId:   req.UserId,
		UserType: int32(req.UserType),
	}
	oRsp, err := l.svcCtx.AssessClient.GetUserOrderSummary(l.ctx, oReq)
	if err != nil {
		l.Logger.Error("call rpc GetUserOrderSummary error:", err)
		return
	}
	//拼装返回数据
	s := types.DBTradeSide{
		BuyRate:  oRsp.GetBuyRate() * 100,
		SellRate: oRsp.GetSellRate() * 100,
	}
	m := types.DBMarketRateInfo{
		HugeRate:   oRsp.GetFundRate().GetHuge() * 100,
		BigRate:    oRsp.GetFundRate().GetBig() * 100,
		MiddleRate: oRsp.GetFundRate().GetMiddle() * 100,
		SmallRate:  oRsp.GetFundRate().GetSmall() * 100,
	}
	out := &types.AlgoComsumRsp{
		UserCnt:          oRsp.UserNum,
		TotalUserCnt:     oRsp.TotalUserNum,
		AlgoCnt:          int64(oRsp.TradeAlgoCount),
		TotalAlgoCnt:     int64(oRsp.AlgoCount),
		TradeVol:         float64(oRsp.TotalTradeAmount) / 10000,
		OrderCnt:         oRsp.OrderNum,
		Side:             s,
		ProviderCnt:      int64(oRsp.TradeProviderCount),
		TotalProviderCnt: int64(oRsp.ProviderCount),
		MarketRate:       m,
		Progress:         oRsp.Progress * 100,
	}

	return out, nil
}
