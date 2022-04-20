// Package logic
/*
 Author: hawrkchen
 Date: 2022/3/24 15:14
 Desc:
*/
package consumer

import (
	"algo_assess/assess-mq-server/proto"
	"algo_assess/global"
	"algo_assess/market-mq-server/internal/svc"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
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
	// s.Logger.Info("cumsume market info|key:", key, ", value:", val)
	// 先解头
	headLen := unsafe.Sizeof(global.QuoteHead{})
	head, err := ParseQuoteHead(val[:headLen])
	if err != nil {
		s.Logger.Error("parse quote head fail:", err)
		return nil
	}
	s.Logger.Info("get head ", head.QuoteID, head.SecBitType, head.QuoteType)
	if head.QuoteType != uint16(global.SubscribeQuoteTypeLevel2) { // 非level2的数据不要
		return nil
	}
	// 解结构体
	level2, err := parsTagTagQuoteClientLevel2Data(val)
	if err != nil {
		s.Logger.Error("parse quote level2 fail:", err)
	}
	s.Logger.Infof("level2:%+v", level2)

	secId := string(level2.SecID[:])
	s.Logger.Info("get secId:", secId)
	// 落redis

	transAct := BuildMarketTime(level2.OrigTime)
	hKey := fmt.Sprintf("level2:%s:%s", level2.SecID[:], transAct)
	s.Logger.Info("get redis key:", hKey)
	s.svcCtx.RedisClient.Hmset(hKey, BuildLevel2ToRedisHash(level2))

	// 推送到 assess mq server
	req := &proto.MarketDataReq{
		UseculityId:   0,
		SecId:         secId,
		EntrustBidVol: level2.TotalBidVol,
		EntrustAskVol: level2.TotalAskVol,
		OrgiTime:      level2.OrigTime,
		TotalTradeVol: level2.TotalTradeVol,
		LastPrice:     uint64(level2.LastPrice),
	}
	s.Logger.Infof("send req:%+v", req)
	rsp, err := s.svcCtx.AssessMQClient.PullMarketData(s.ctx, req)
	if err != nil {
		s.Logger.Error("call assess mq server fail:", err)
		return nil
	}
	//TODO: 其实回执都不需要
	s.Logger.Info("get rsp code:", rsp.Code)
	time.Sleep(time.Second * 10)

	return nil
}
