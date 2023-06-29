package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	pt "google.golang.org/protobuf/proto"
	"strings"
	"time"
)

type PushChildOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPushChildOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PushChildOrderLogic {
	return &PushChildOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// PushChildOrder 子单绩效修复
func (l *PushChildOrderLogic) PushChildOrder(in *proto.ChildOrderPerfs) (*proto.PushDataRsp, error) {
	parts := in.GetParts()
	// 先校验一下母单是否已上传
	/*
		var algoIds []int
		for _, v := range parts {
			if !tools.InSlice(algoIds, int(v.AlgoOrderId)) {
				algoIds = append(algoIds, int(v.AlgoOrderId))
			}
		}
		ids, err := l.svcCtx.AlgoOrderRepo.BatchQueryAlgoOrder(l.ctx, algoIds)
		if err != nil {
			l.Logger.Error("BatchQueryAlgoOrder error:", err)
			return &proto.PushDataRsp{
				Code: 10000,
				Msg:  "检查母单信息失败",
			}, nil
		}
		if len(ids) != len(algoIds) { // 检查母单信息是否满足
			l.Logger.Error("algo order not match...")
			return &proto.PushDataRsp{
				Code: 10000,
				Msg:  "母单信息缺失，请先导入母单数据",
			}, nil
		}
	*/
	var date string    // 格式 "202305056"
	var reqIds []int64 // 请求子单列表
	for _, v := range parts {
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

	for _, part := range parts {
		t := time.UnixMicro(int64(part.GetTransactTime())).Format(global.TimeFormatMinInt)
		transactAt := tools.TimeMoveForward(t)
		date = cast.ToString(transactAt)[:8]
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
		if _, exist := m[int64(part.GetId())]; !exist { // 子单号不存在，补一份到DB
			if err := l.svcCtx.OrderDetailRepo.CreateOrderDetail(l.ctx, child); err != nil {
				l.Logger.Error("CreateOrderDetail error:", err)
				continue
			}
		} else { // 子单数据有可会改变，这里需要更新
			if err := l.svcCtx.OrderDetailRepo.UpdateOrderDetail(l.ctx, child); err != nil {
				l.Logger.Error("UpdateOrderDetail error:", err)
				continue
			}
		}
	}
	// 关联的算法ID
	var algorithmIds []int32
	mId := make(map[int32]struct{}) // 辅助 childOrderIds 去重
	for _, v := range parts {
		if _, exist := mId[int32(v.AlgorithmId)]; !exist {
			algorithmIds = append(algorithmIds, int32(v.AlgorithmId))
			mId[int32(v.AlgorithmId)] = struct{}{}
		}
	}
	// 关联母单信息
	l.PushAlgoOrder2Kafka(algorithmIds, date)

	// 关联子单信息
	l.PushChildOrder2Kafka(algorithmIds, date)

	return &proto.PushDataRsp{
		Code: 200,
		Msg:  "处理成功",
	}, nil
}

// PushAlgoOrder2Kafka 根据算法ID 找到当天的所有母单信息，然后推到Kafka
func (l *PushChildOrderLogic) PushAlgoOrder2Kafka(algothmIds []int32, date string) {
	//把所有该算法ID的母单号找出来
	out, err := l.svcCtx.AlgoOrderRepo.GetAlgoByAlgorithms(l.ctx, algothmIds, date)
	if err != nil {
		l.Logger.Error("get algoOrder List error:", err)
		return
	}

	for _, v := range out {
		o := &proto.AlgoOrderPerf{
			Id:            uint32(v.AlgoId),
			BasketId:      uint32(v.BasketId),
			AlgorithmType: uint32(v.AlgorithmType),
			AlgorithmId:   uint32(v.AlgorithmId),
			USecurityId:   uint32(v.UsecId),
			SecurityId:    v.SecId,
			AlgoOrderQty:  uint64(v.AlgoOrderQty) * 100,
			TransactTime:  uint64(v.UnixTime),
			StartTime:     uint64(v.StartTime),
			EndTime:       uint64(v.EndTime),
			BusUserId:     v.UserId,
			BatchNo:       v.BatchNo,
			BatchName:     v.BatchName,
			SourceFrom:    1,
		}
		// 发送Kafka
		bytes, ptErr := pt.Marshal(o)
		if ptErr != nil {
			l.Logger.Error("Marshal msg error:", ptErr)
			continue
		}
		l.Logger.Infof("push ProtoMessage:%+v", o)
		msg := &sarama.ProducerMessage{
			Topic:     l.svcCtx.Config.Kafka.AlgoTopic, // 统一放入正常队列
			Value:     sarama.ByteEncoder(bytes),
			Partition: 0, // 指定0分区，需要在初始化时指定sarama.NewManualPartitioner
		}
		_, _, kafkaErr := l.svcCtx.SyncProducer.SendMessage(msg)
		if kafkaErr != nil {
			l.Logger.Error("send message error:", msg.Value, kafkaErr)
			continue
		}
	}
}

func (l *PushChildOrderLogic) PushChildOrder2Kafka(algothmIds []int32, date string) {
	//把所有该算法ID的子单号找出来
	out, err := l.svcCtx.OrderDetailRepo.GetChildByAlgorithms(l.ctx, algothmIds, date)
	if err != nil {
		l.Logger.Error("get GetChildOrder List error:", err)
		return
	}
	for _, v := range out {
		o := &proto.ChildOrderPerf{
			Id:        uint32(v.ChildOrderId),
			BusUserId: v.UserId,
			//BusUuserId:     v.UsecurityId,
			AlgoOrderId:    uint32(v.AlgoOrderId),
			AlgorithmType:  uint32(v.AlgorithmType),
			AlgorithmId:    uint32(v.AlgorithmId),
			USecurityId:    uint32(v.UsecurityId),
			SecurityId:     v.SecurityId,
			Side:           uint32(v.TradeSide),
			OrderQty:       uint64(v.OrderQty) * 100,
			Price:          uint64(v.Price),
			OrderType:      uint32(v.OrderType),
			CumQty:         uint64(v.ComQty) * 100,
			LastPx:         uint64(v.LastPx),
			LastQty:        uint64(v.LastQty) * 100,
			Charge:         v.TotalFee,
			ArrivedPrice:   uint64(v.ArrivedPrice),
			ChildOrdStatus: uint32(v.OrdStatus),
			TransactTime:   uint64(v.TransactTime),
			Version:        0,
			BatchNo:        v.BatchNo,
			BatchName:      v.BatchName,
			SourceFrom:     1,
		}
		// 发送kafka
		bytes, ptErr := pt.Marshal(o)
		if ptErr != nil {
			l.Logger.Error("Marshal msg error:", ptErr)
			continue

		}
		l.Logger.Infof("push ProtoMessage:%+v", o)
		msg := &sarama.ProducerMessage{
			Topic:     l.svcCtx.Config.Kafka.ChildTopic, // 放修复的子单队列
			Value:     sarama.ByteEncoder(bytes),
			Partition: 0, // 指定0分区，需要在初始化时指定sarama.NewManualPartitioner
		}
		_, _, kafkaErr := l.svcCtx.SyncProducer.SendMessage(msg)
		if kafkaErr != nil {
			l.Logger.Errorf("send message(%s) err=%s \n", msg.Value, kafkaErr)
			continue
		}
	}
}
