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
	//l.Logger.Infof("get req:%+v", in)
	// 收到行情推送数据，然后保存到本地缓存就可以了  --去掉前缀
	curKey := fmt.Sprintf("%s:%d", in.GetSecId()[3:], in.GetOrgiTime())
	//lastKey := fmt.Sprintf("%s:%s", in.GetSecId(), BuildLastMarketTime(in.GetOrgiTime()))
	//l.Logger.Info("get curKey:", curKey, ", get lastKey:", lastKey)

	entrustVol := in.GetEntrustAskVol() + in.GetEntrustBidVol()
	//l.Logger.Info("get CurEntrustVol :", entrustVol)
	global.GlobalMarketLevel2.RWMutex.Lock()
	global.GlobalMarketLevel2.EntrustVol[curKey] = entrustVol
	global.GlobalMarketLevel2.TradeVol[curKey] += in.GetNetTotalTradeVol() // 成交总增量, 上游已经计算好了
	global.GlobalMarketLevel2.LastPrice[curKey] = in.GetLastPrice()
	// l.Logger.Infof("get entrust map:%+v", global.GlobalMarketLevel2.EntrustVol)
	// l.Logger.Infof("get total trade map:%+v", global.GlobalMarketLevel2.TradeVol)
	// l.Logger.Infof("get last price  map:%+v", global.GlobalMarketLevel2.LastPrice)
	global.GlobalMarketLevel2.RWMutex.Unlock()

	return &proto.MarketDataRsp{}, nil
}
