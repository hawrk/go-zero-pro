package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	pb "algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type PushSzQuoteLevelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushSzQuoteLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushSzQuoteLevelLogic {
	return &PushSzQuoteLevelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  深市行情信息
func (l *PushSzQuoteLevelLogic) PushSzQuoteLevel(in *pb.ReqPushSzLevel) (*pb.PushDataRsp, error) {
	l.Logger.Info("into PushSzQuoteLevel ...")
	//TODO:
	// 先直接入库
	for _, v := range in.GetQuote() {
		if err := l.svcCtx.MarketLevelRepo.CreateMarketLevel(l.ctx, &global.QuoteLevel2Data{
			SecID:         v.SeculityId,
			OrigTime:      v.OrgiTime,
			LastPrice:     v.LastPrice,
			TotalTradeVol: v.TotalTradeVol,
			AskPrice:      v.AskPrice,
			AskVol:        v.AskVol,
			BidPrice:      v.BidPrice,
			BidVol:        v.BidVol,
			TotalBidVol:   v.TotalBidVol,
			TotalAskVol:   v.TotalAskVol,
			Vwap:          0,
		}); err != nil {
			l.Logger.Error("CreateMarketLevel error:", err)
		}
	}

	/*
		// 这里行情修复不需要写表，在行情计算处理完成后会写入
		topic := l.svcCtx.Config.Kafka.SZMarketTopic
		for _, v := range in.GetQuote() {
			b, err := proto.Marshal(v)
			if err != nil {
				l.Logger.Error("Marshal error:", err)
				continue
			}
			msg := &sarama.ProducerMessage{
				Topic:     topic,
				Value:     sarama.ByteEncoder(b),
				Partition: 0,
			}
			_, _, kafkaErr := l.svcCtx.SyncProducer.SendMessage(msg)
			if kafkaErr != nil {
				l.Logger.Error("send message(%s) err=%s", msg.Value, kafkaErr)
				return &pb.PushDataRsp{
					Code: 10000,
					Msg:  "kafka推送错误",
				}, nil
			}
		}

	*/

	return &pb.PushDataRsp{
		Code: 200,
		Msg:  "success",
	}, nil
}
