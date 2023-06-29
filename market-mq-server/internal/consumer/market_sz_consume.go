// Package logic
/*
 Author: hawrkchen
 Date: 2022/3/24 15:14
 Desc:
*/
package consumer

import (
	"algo_assess/global"
	"algo_assess/market-mq-server/internal/svc"
	"algo_assess/pkg/tools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"unsafe"
)

type AlgoPlatformMarketInfo struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoPlatformMarketInfo(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoPlatformMarketInfo {
	return &AlgoPlatformMarketInfo{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *AlgoPlatformMarketInfo) Consume(key string, val string) error {
	//s.Logger.Info("cumsume market info|key:", key, ", value:", val)
	//time.Sleep(time.Second * 1)
	// 先解头
	headLen := unsafe.Sizeof(global.QuoteHead{})
	head, err := ParseQuoteHead(val[:headLen])
	if err != nil {
		s.Logger.Error("parse quote head fail:", err)
		return nil
	}
	if head.QuoteType != uint16(global.SubscribeQuoteTypeLevel2) { // 非level2的数据不要
		s.Logger.Info("not level2,continue, type:", head.QuoteType)
		return nil
	}
	// 解结构体
	level2, err := parsTagTagQuoteClientLevel2Data(val)
	if err != nil {
		s.Logger.Error("parse quote level2 fail:", err)
	}

	//s.Logger.Infof("level2:%+v", level2)
	//secId := strings.TrimSpace(tools.Bytes2String(level2.SecID[:]))
	//s.Logger.Info("get secId:", secId)
	//time.Sleep(time.Millisecond * 100)
	// 过滤 休市的数据   深交所交易时间格式 HHMMSSsss,精确到毫秒
	if level2.OrigTime < 92800000 ||
		(level2.OrigTime > 113000000 && level2.OrigTime < 130000000) ||
		level2.OrigTime > 150100000 {
		//s.Logger.Info("unarrive time , contine:", level2.OrigTime)
		return nil
	}
	// TODO: test

	//if secId != "000038" {
	//	//s.Logger.Info("unnessary secID: ", secId)
	//	return nil
	//}
	// 过滤债券和期权   49 -> ascii为1    50 -> ascii 为2
	if level2.SecID[0] == 49 || level2.SecID[0] == 50 {
		//s.Logger.Info("unnecessary secID.... ")
		return nil
	}
	//s.Logger.Infof("get level2:%+v", level2)
	// trans quote data
	szQuote := TransSzQuoteData(&level2)
	// 为避免网络延迟，数据推送出现时间乱序问题，校验一下时间顺序
	if szQuote.OrigTime < global.SzOriginTime {
		s.Logger.Error("Error: current OriginTime less than pre OriginTime:", szQuote.OrigTime, global.ShOriginTime)
		return nil
	}
	global.SzOriginTime = szQuote.OrigTime
	//s.Logger.Infof("get sz szQuote:%+v", szQuote)

	// 实时计算市场vwap
	netTotalVol := CalculateMarketVwap(&szQuote)
	//s.Logger.Info("get netTotalVol:", netTotalVol)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		wg.Done()
		// 落mysql
		// 优化： 本地缓存数据， 避免DB先查一次再插入或更新操作
		//quoteKey := szQuote.SecID + ":" + cast.ToString(szQuote.OrigTime)
		////s.Logger.Info("get sz quate key:", quoteKey)
		//val, _ := global.QuateKeyMap.Load(quoteKey)
		//_, ok := val.(int)
		//if !ok {
		//	global.QuateKeyMap.Store(quoteKey, 1)
		//	if err := s.svcCtx.MarketLevelRepo.CreateMarketLevel(s.ctx, &szQuote); err != nil {
		//		s.Logger.Error("insert db market level fail:", err)
		//	}
		//} else {
		//	if err := s.svcCtx.MarketLevelRepo.UpdateMarketLevel(s.ctx, &szQuote); err != nil {
		//		s.Logger.Error("update db market level fail:", err)
		//	}
		//}
	}()

	go func() {
		wg.Done()
		// 数据存Redis
		// 格式hash    key->   sz000038:20220520   field -> 0930   value -> QuoteLevel2Data结构化Json
		t := cast.ToString(szQuote.OrigTime)
		hashKey := fmt.Sprintf("%s:%s", szQuote.SecID, t[:8])
		hashValue, _ := json.Marshal(szQuote)
		//s.Logger.Info("get haskKey:", hashKey, ", hashField:", hashField, ", hashValue:", hashValue)
		m := map[string]string{
			t: tools.Bytes2String(hashValue),
		}
		// 保存DB的信息，缓存落地信息
		if err := s.svcCtx.RedisClient.Hmset(hashKey, m); err != nil {
			s.Logger.Error("set redis hash fail:", err)
		}
		if err := BuildLevel2Data(s.svcCtx.RedisClient, netTotalVol, &szQuote); err != nil {
			s.Logger.Error("build level2 data fail:", err)
		}

		// 反序列
		//var out global.QuoteLevel2Data
		//if err := json.Unmarshal(hashValue, &out); err != nil {
		//	s.Logger.Error("unnarshal error:", err)
		//}
		//s.Logger.Infof("get unmarshal data:%+v", out)

	}()

	go func() {
		wg.Done()
		// 推送到 assess mq server
		//req := &proto.MarketDataReq{
		//	UseculityId:      0,
		//	SecId:            szQuote.SecID,
		//	EntrustBidVol:    szQuote.TotalBidVol,
		//	EntrustAskVol:    szQuote.TotalAskVol,
		//	OrgiTime:         szQuote.OrigTime,
		//	TotalTradeVol:    szQuote.TotalTradeVol,
		//	LastPrice:        szQuote.LastPrice,
		//	NetTotalTradeVol: netTotalVol,
		//}
		////s.Logger.Infof("send req:%+v", req)
		//_, err = s.svcCtx.AssessMQClient.PullMarketData(s.ctx, req)
		//if err != nil {
		//	s.Logger.Error("call assess mq server fail:", err)
		//}
	}()
	wg.Wait()
	//time.Sleep(time.Second * 10)
	//s.Logger.Info("done....")
	return nil
}
