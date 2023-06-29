// Package consumer
/*
 Author: hawrkchen
 Date: 2022/5/16 15:51
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
	"strings"
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
	// 二进制val 长度在350左右，字符串val 长度在250到300左右
	if key == "MarketInfo3" {
		return s.ConsumeStringMsg(val)
	} else {
		return s.ConsumeBinaryMsg(val)
	}
}

func (s *AlgoPlatformSHMarketInfo) ConsumeBinaryMsg(val string) error {
	//time.Sleep(time.Second * 1)
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
	//s.Logger.Infof("get level2:%+v", level2)
	//secId := strings.TrimSpace(tools.Bytes2String(level2.SecID[:]))
	//s.Logger.Info("get secId:", secId)
	// 上交所时间格式为 HHMMSS  拿 9:29的数据需要填充9:30--hawrk 20221216 统一成HHMMSSsss精确到毫秒
	if level2.OrigTime < 92900000 ||
		(level2.OrigTime > 113000000 && level2.OrigTime < 130000000) ||
		level2.OrigTime > 150100000 {
		//s.Logger.Info("un arrive time , continue:", level2.OrigTime)
		return nil
	}

	//s.Logger.Infof("get level2:%+v", level2)
	// 过滤债券和期权   49 -> ascii为1    50 -> ascii 为2
	//secId := strings.TrimSpace(tools.Bytes2String(level2.SecID[:]))
	if level2.SecID[0] == 48 || level2.SecID[0] == 49 || level2.SecID[0] == 50 { // 沪市的话一般保留6开头的就可以了
		//s.Logger.Info("unnessary secID.... ", secId)
		return nil
	}
	// TODO: secid  600346
	//secId := strings.TrimSpace(tools.Bytes2String(level2.SecID[:]))
	//if secId != "600436" {
	//	//s.Logger.Info("not my secid :", secId)
	//	return nil
	//}

	//s.Logger.Infof("get sh level2:%+v", level2)
	shQuote := TransShQuoteData(&level2)
	// 为避免网络延迟，数据推送出现时间乱序问题，校验一下时间顺序
	if shQuote.OrigTime < global.ShOriginTime {
		s.Logger.Error("Error: current OriginTime less than pre OriginTime:", shQuote.OrigTime, global.ShOriginTime)
		return nil
	}
	global.ShOriginTime = shQuote.OrigTime
	//s.Logger.Info("get originTime:", global.ShOriginTime, ", currentTime :", shQuote.OrigTime)

	//s.Logger.Infof("get shQuote:%+v", shQuote)

	// 实时计算市场vwap
	netTotalVol := CalculateMarketVwap(&shQuote)
	// 数据存Redis
	// 格式hash    key->   sz000038:20220520   field -> 0930   value -> QuoteLevel2Data结构化Json
	t := cast.ToString(shQuote.OrigTime)
	hashKey := fmt.Sprintf("%s:%s", shQuote.SecID, t[:8])
	hashValue, _ := json.Marshal(shQuote)
	//s.Logger.Info("get haskKey:", hashKey, ", hashField:", hashField, ", hashValue:", hashValue)
	m := map[string]string{
		t: tools.Bytes2String(hashValue),
	}
	// 落库用，数据本身不参与实时计算
	if err := s.svcCtx.RedisClient.Hmset(hashKey, m); err != nil {
		s.Logger.Error("set redis hash fail:", err)
	}
	// 写入redis,一期计算市场指标时用到
	if err := BuildLevel2Data(s.svcCtx.RedisClient, netTotalVol, &shQuote); err != nil {
		s.Logger.Error("build level2 data fail:", err)
	}
	//time.Sleep(time.Second*1)
	//s.Logger.Info("done...")
	return nil
}

func (s *AlgoPlatformSHMarketInfo) ConsumeStringMsg(val string) error {
	shQuote := s.TransStr2Struct(val)
	if shQuote == nil {
		s.Logger.Error("parse string msg error")
		return nil
	}
	// 实时计算市场vwap
	netTotalVol := CalculateMarketVwap(shQuote)
	// 数据存Redis
	// 格式hash    key->   sz000038:20220520   field -> 0930   value -> QuoteLevel2Data结构化Json
	t := cast.ToString(shQuote.OrigTime)
	hashKey := fmt.Sprintf("%s:%s", shQuote.SecID, t[:8])
	hashValue, _ := json.Marshal(shQuote)
	//s.Logger.Info("get haskKey:", hashKey, ", hashField:", hashField, ", hashValue:", hashValue)
	m := map[string]string{
		t: tools.Bytes2String(hashValue),
	}
	if err := s.svcCtx.RedisClient.Hmset(hashKey, m); err != nil {
		s.Logger.Error("set redis hash fail:", err)
	}
	if err := BuildLevel2Data(s.svcCtx.RedisClient, netTotalVol, shQuote); err != nil {
		s.Logger.Error("build level2 data fail:", err)
	}
	return nil
}

func (s *AlgoPlatformSHMarketInfo) TransStr2Struct(val string) *global.QuoteLevel2Data {
	sli := strings.Split(val, "|")
	if len(sli) < 10 {
		return nil
	}
	shQuote := &global.QuoteLevel2Data{
		SecID:         sli[0],
		OrigTime:      cast.ToInt64(sli[1]),
		LastPrice:     cast.ToInt64(sli[2]),
		TotalTradeVol: cast.ToInt64(sli[3]),
		AskPrice:      sli[4],
		AskVol:        sli[5],
		BidPrice:      sli[6],
		BidVol:        sli[7],
		TotalBidVol:   cast.ToInt64(sli[8]),
		TotalAskVol:   cast.ToInt64(sli[9]),
	}

	return shQuote
}
