package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
	pt "google.golang.org/protobuf/proto"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportChildOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportChildOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportChildOrderLogic {
	return &ImportChildOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ImportChildOrder 订单导入：子单信息导入
func (l *ImportChildOrderLogic) ImportChildOrder(in *proto.ChildOrderPerfs) (*proto.PushDataRsp, error) {
	l.Logger.Info("in ImportChildOrder, get req len:", len(in.GetParts()))
	var reqIds []int64 // 请求子单列表
	for _, v := range in.GetParts() {
		reqIds = append(reqIds, int64(v.Id))
	}
	childIds, err := l.svcCtx.OrderDetailRepo.BatchQueryChildOrder(l.ctx, reqIds)
	if err != nil {
		l.Logger.Error("BatchQueryChildOrder error:", err)
		return &proto.PushDataRsp{
			Code: 10000,
			Msg:  "子单信息查询失败",
		}, nil
	}

	// db 子单列表转成map
	m := make(map[int64]struct{}, len(childIds))
	for _, v := range childIds {
		m[v] = struct{}{}
	}

	for _, part := range in.GetParts() {
		if _, exist := m[int64(part.GetId())]; !exist {
			t := time.UnixMicro(int64(part.GetTransactTime())).Format(global.TimeFormatMinInt)
			transactAt := tools.TimeMoveForward(t)
			child := &global.ChildOrderData{
				BatchNo:          global.DefaultBatchNo, // 与总线保持一致
				OrderId:          int64(part.Id),
				AlgoOrderId:      int64(part.AlgoOrderId),
				AlgorithmType:    int(part.AlgorithmType),
				AlgoId:           int(part.AlgorithmId),
				UsecId:           uint(part.USecurityId),
				UserId:           strings.TrimSpace(tools.RMu0000(part.BusUserId)),
				SecId:            strings.TrimSpace(part.SecurityId),
				TradeSide:        tools.GetOrderTradeSide(part.Side), // 买卖方向    1-买    2-卖
				OrderQty:         int64(part.OrderQty) / 100,
				Price:            int64(part.Price), // 委托价格， 价格先不转
				OrderType:        uint(part.OrderType),
				LastPx:           int64(part.LastPx), // 成交价格
				LastQty:          int64(part.LastQty) / 100,
				ComQty:           int64(part.CumQty) / 100,
				ArrivePrice:      int64(part.ArrivedPrice),  // 到达价格
				TotalFee:         cast.ToInt64(part.Charge), // 总手续费
				ChildOrderStatus: uint(part.ChildOrdStatus),
				TransTime:        transactAt,
				UnixTimeMillSec:  int64(part.TransactTime),
				SourceFrom:       int(part.SourceFrom), // 数据修复标识
			}

			if err := l.svcCtx.OrderDetailRepo.CreateOrderDetail(l.ctx, child); err != nil {
				l.Logger.Error("CreateOrderDetail error:", err)
				continue
			}
		}

		// 发送kafka
		bytes, ptErr := pt.Marshal(part)
		if ptErr != nil {
			l.Logger.Error("Marshal msg error:", ptErr)
			continue
			//return &proto.PushDataRsp{
			//	Code: 10000,
			//	Msg:  "proto格式错误",
			//}, nil
		}
		msg := &sarama.ProducerMessage{
			Topic:     l.svcCtx.Config.Kafka.ChildTopic, // 放修复的子单队列
			Value:     sarama.ByteEncoder(bytes),
			Partition: 0, // 指定0分区，需要在初始化时指定sarama.NewManualPartitioner
		}
		_, _, kafkaErr := l.svcCtx.SyncProducer.SendMessage(msg)
		if kafkaErr != nil {
			l.Logger.Errorf("send message(%s) err=%s \n", msg.Value, kafkaErr)
			continue
			//return &proto.PushDataRsp{
			//	Code: 10000,
			//	Msg:  "kafka推送错误",
			//}, nil
		}
	}

	return &proto.PushDataRsp{
		Code: 200,
		Msg:  "success",
	}, nil
}
