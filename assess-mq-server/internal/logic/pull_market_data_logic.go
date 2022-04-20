package logic

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"
	"algo_assess/global"
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

type PullMarketDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPullMarketDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PullMarketDataLogic {
	return &PullMarketDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  推送行情数据
func (l *PullMarketDataLogic) PullMarketData(in *proto.MarketDataReq) (*proto.MarketDataRsp, error) {
	l.Logger.Infof("get req:%+v", in)
	// 收到行情推送数据，然后保存到本地缓存就可以了
	curKey := fmt.Sprintf("%s:%s", in.GetSecId(), BuildCurMarketTime(in.GetOrgiTime()))
	lastKey := fmt.Sprintf("%s:%s", in.GetSecId(), BuildLastMarketTime(in.GetOrgiTime()))
	l.Logger.Info("get curKey:", curKey, ", get lastKey:", lastKey)

	qty := in.GetEntrustAskVol() + in.GetEntrustBidVol() // 当前委托数量
	global.QuoteMap[curKey] = qty                        // 当前时间委托数量
	entrustVol := qty - global.QuoteMap[lastKey]         // 当前时间委托数量 - 前一分钟委托数量

	global.GlobalMarketLevel2.RWMutex.Lock()
	global.GlobalMarketLevel2.EntrustVol[curKey] = entrustVol
	global.GlobalMarketLevel2.TradeVol[curKey] = in.GetTotalTradeVol() - global.GlobalMarketLevel2.TradeVol[lastKey]
	global.GlobalMarketLevel2.LastPrice[curKey] = int64(in.GetLastPrice())

	l.Logger.Infof("get entrust map:%+v", global.GlobalMarketLevel2.EntrustVol)
	l.Logger.Infof("get total trade map:%+v", global.GlobalMarketLevel2.TradeVol)
	l.Logger.Infof("get last price  map:%+v", global.GlobalMarketLevel2.LastPrice)
	global.GlobalMarketLevel2.RWMutex.Unlock()

	return &proto.MarketDataRsp{}, nil
}
