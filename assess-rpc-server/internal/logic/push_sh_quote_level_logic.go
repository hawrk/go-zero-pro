package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	pb "algo_assess/assess-rpc-server/proto"
	"context"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type PushShQuoteLevelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushShQuoteLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushShQuoteLevelLogic {
	return &PushShQuoteLevelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  沪市行情信息推送
func (l *PushShQuoteLevelLogic) PushShQuoteLevel(in *pb.ReqPushShLevel) (*pb.PushDataRsp, error) {
	l.Logger.Info("into PushShQuoteLevel ...")
	topic := l.svcCtx.Config.Kafka.SHMarketTopic
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
	return &pb.PushDataRsp{
		Code: 200,
		Msg:  "kafka推送成功",
	}, nil
}
