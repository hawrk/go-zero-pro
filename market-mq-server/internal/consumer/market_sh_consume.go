// Package consumer
/*
 Author: hawrkchen
 Date: 2022/5/16 15:51
 Desc:
*/
package consumer

import (
	"algo_assess/assess-mq-server/proto"
	"algo_assess/global"
	"algo_assess/market-mq-server/internal/svc"
	"algo_assess/pkg/tools"
	"context"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type AlgoPlatformSHMarketInfo struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoPlatformSHMarketInfo(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoPlatformSHMarketInfo {
	return &AlgoPlatformSHMarketInfo{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *AlgoPlatformSHMarketInfo) Consume(key string, val string) error {
	//s.Logger.Info("cumsume market info|key:", key, ", value:", val)
	headLen := unsafe.Sizeof(global.QuoteHead{})
	head, err := ParseQuoteHead(val[:headLen])
	if err != nil {
		s.Logger.Error("parse quote head fail:", err)
		return nil
	}
	//s.Logger.Info("get head:", head)
	if head.QuoteType != uint16(global.SubscribeQuoteTypeLevel2) { // 非level2的数据不要
		s.Logger.Info("not level2,continue, type:", head.QuoteType)
		return nil
	}
	// 解结构体
	level2, err := parsTagTagQuoteClientLevel2Data(val)
	if err != nil {
		s.Logger.Error("parse quote level2 fail:", err)
	}
	// 上交所时间格式为 HHMMSS
	if level2.OrigTime < 93000 ||
		(level2.OrigTime > 113000 && level2.OrigTime < 130000) ||
		level2.OrigTime > 150100 {
		//s.Logger.Info("unarrive time , contine:", level2.OrigTime)
		return nil
	}
	// 过滤债券和期权   49 -> ascii为1    50 -> ascii 为2
	//secId := strings.TrimSpace(tools.Bytes2String(level2.SecID[:]))
	if level2.SecID[0] == 48 || level2.SecID[0] == 49 || level2.SecID[0] == 50 { // 沪市的话一般保留6开头的就可以了
		//s.Logger.Info("unnessary secID.... ", secId)
		return nil
	}
	// TODO: secid  600346
	//secId := strings.TrimSpace(tools.Bytes2String(level2.SecID[:]))
	//if secId != "600346" {
	//	//s.Logger.Info("not my secid :", secId)
	//	return nil
	//}
	shQuote := s.TransShQuoteData(&level2)
	//s.Logger.Infof("get shQuote:%+v", shQuote)
	// 实时计算市场vwap
	netTotalVol := CalculateMarketVwap(&shQuote)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		wg.Done()
		quoteKey := shQuote.SecID + ":" + cast.ToString(shQuote.OrigTime)
		v, _ := global.QuateKeyMap.Load(quoteKey)
		_, ok := v.(int)
		if !ok {
			global.QuateKeyMap.Store(quoteKey, 1)
			if err := s.svcCtx.ShMarketLevelRepo.CreateShMarketLevel(s.ctx, &shQuote); err != nil {
				s.Logger.Error("insert db market level fail:", err)
			}
		} else {
			if err := s.svcCtx.ShMarketLevelRepo.UpdateShMarketLevel(s.ctx, &shQuote); err != nil {
				s.Logger.Error("update db market level fail:", err)
			}
		}
	}()

	go func() {
		wg.Done()
		// 推送到 assess mq server
		req := &proto.MarketDataReq{
			UseculityId:      0,
			SecId:            shQuote.SecID,
			EntrustBidVol:    shQuote.TotalBidVol,
			EntrustAskVol:    shQuote.TotalAskVol,
			OrgiTime:         shQuote.OrigTime,
			TotalTradeVol:    shQuote.TotalTradeVol,
			LastPrice:        shQuote.LastPrice,
			NetTotalTradeVol: netTotalVol,
		}
		//s.Logger.Infof("send req:%+v", req)
		_, err = s.svcCtx.AssessMQClient.PullMarketData(s.ctx, req)
		if err != nil {
			s.Logger.Error("call assess mq server fail:", err)
		}
	}()
	wg.Wait()
	time.Sleep(time.Millisecond * 1)
	return nil
}

// TransShQuoteData 转换上交所行情数据结构
func (s *AlgoPlatformSHMarketInfo) TransShQuoteData(data *global.TagQuoteClientLevel2Data) global.QuoteLevel2Data {
	out := global.QuoteLevel2Data{
		SecID:         "sh:" + strings.TrimSpace(tools.Bytes2String(data.SecID[:])),
		OrigTime:      tools.GetDayByEnv(s.svcCtx.Config.Deployment.Env, data.OrigTime),
		LastPrice:     int64(data.LastPrice) / 10,
		TotalTradeVol: int64(data.TotalTradeVol) / 1000,
		AskPrice:      tools.Byte32ArrayToStringForSh(data.AskPrice),
		AskVol:        tools.Byte64ArrayToStringForSh(data.AskVol),
		BidPrice:      tools.Byte32ArrayToStringForSh(data.BidPrice),
		BidVol:        tools.Byte64ArrayToStringForSh(data.BidVol),
		TotalAskVol:   int64(data.TotalAskVol) / 1000,
		TotalBidVol:   int64(data.TotalBidVol) / 1000,
	}
	return out
}
