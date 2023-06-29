package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	pt "google.golang.org/protobuf/proto"
	"strings"
)

type OrigOrderAnalyseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrigOrderAnalyseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrigOrderAnalyseLogic {
	return &OrigOrderAnalyseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// OrigOrderAnalyse 原始订单分析
func (l *OrigOrderAnalyseLogic) OrigOrderAnalyse(in *proto.OrigAnalyseReq) (*proto.OrigAnalyseResp, error) {
	l.Logger.Infof("OrigAnalyseReq:%v", in)
	//var algoOrder []*models.TbAlgoOrder
	//var childOrder []*models.TbAlgoOrderDetail
	//var err error
	batchNo := tools.GeneralID()
	var timeList []int64
	if in.GetOrderType() == 1 { // 母单
		// 根据日期和ID找母单
		for _, v := range in.GetOrders() {
			out, err := l.svcCtx.AlgoOrderRepo.QueryAlgoOrderByDate(v.GetDate(), v.GetId())
			if err != nil {
				l.Logger.Error("QueryAlgoOrderByDate error:", err)
				continue
			}
			if out.Id <= 0 { // 判断是否有记录
				continue
			}
			timeList = append(timeList, out.TransTime)

			l.SendAlgoOrderKafka(batchNo, out)

			// 根据母单继续找子单
			childs, err := l.svcCtx.OrderDetailRepo.QueryChildOrderByDate(out.Date, out.AlgoId)
			if err != nil {
				l.Logger.Error("QueryChildOrerByDate error:", err)
				continue
			}
			for _, cv := range childs {
				l.SendChildOrderKafka(batchNo, *cv)
			}
		}
	} else if in.GetOrderType() == 2 { // 子单
		var algos []string             // 存放母单非重复ID
		m := make(map[string]struct{}) // 辅助 判断algos列表是否重复
		var perfChilds []models.TbAlgoOrderDetail
		for _, v := range in.GetOrders() {
			child, err := l.svcCtx.OrderDetailRepo.QueryChildOrderByChildOrder(v.GetDate(), v.GetId())
			if err != nil {
				l.Logger.Error("QueryChildOrderByChildOrder: ", err)
				continue
			}
			if child.Id <= 0 { // 判断是否有记录
				continue
			}
			// 以 date:algoOrderId 作为key,判断母单号是否已存在
			a := fmt.Sprintf("%d:%d", child.Date, child.AlgoOrderId)
			if _, exist := m[a]; !exist {
				algos = append(algos, a)
				m[a] = struct{}{}
			}
			perfChilds = append(perfChilds, child) // 保存好该子单信息列表
		}
		l.Logger.Infof("get algos:%+v", algos)
		// 找到归属母单后，把该母单的详情查回来
		for _, v := range algos {
			arr := strings.Split(v, ":")
			if len(arr) < 2 {
				continue
			}
			out, err := l.svcCtx.AlgoOrderRepo.QueryAlgoOrderByDate(cast.ToInt32(arr[0]), cast.ToInt32(arr[1]))
			if err != nil {
				l.Logger.Error("QueryAlgoOrderByDate error:", err)
				continue
			}
			if out.Id <= 0 { // 判断是否有记录
				continue
			}
			timeList = append(timeList, out.TransTime)

			l.SendAlgoOrderKafka(batchNo, out)

		}
		// 发完母单后，再发子单
		for _, cc := range perfChilds {
			l.SendChildOrderKafka(batchNo, cc)
		}

	} else {
		l.Logger.Error("unknown orderType")
		return &proto.OrigAnalyseResp{
			Code: 305,
			Msg:  "unknown orderType",
		}, nil
	}
	start, end := tools.GetDoubleTimePoint(timeList)
	l.Logger.Info("OrigOrderAnalyse general batchNo:", batchNo, ",start:", start, ",end:", end)
	return &proto.OrigAnalyseResp{
		Code:      200,
		Msg:       "success",
		BatchNo:   batchNo,
		StartTime: start,
		EndTime:   end,
	}, nil
}

// SendAlgoOrderKafka 发送母单Kafka
func (l *OrigOrderAnalyseLogic) SendAlgoOrderKafka(batchNo int64, algo models.TbAlgoOrder) {
	b := &proto.AlgoOrderPerf{
		Id:            uint32(algo.AlgoId), // 母单ID
		BasketId:      uint32(algo.BasketId),
		AlgorithmType: uint32(algo.AlgorithmType),
		AlgorithmId:   uint32(algo.AlgorithmId), // 算法id
		USecurityId:   uint32(algo.UsecId),
		SecurityId:    algo.SecId,
		AlgoOrderQty:  uint64(algo.AlgoOrderQty) * 100,
		TransactTime:  uint64(algo.UnixTime),
		StartTime:     uint64(algo.StartTime),
		EndTime:       uint64(algo.EndTime),
		BusUserId:     algo.UserId,
		BatchNo:       batchNo,
		BatchName:     "",
		SourceFrom:    2,
	}
	// 发送Kafka
	bytes, err := pt.Marshal(b)
	if err != nil {
		l.Logger.Error("Marshal msg error:", err)
		return
	}
	msg := &sarama.ProducerMessage{
		Topic:     l.svcCtx.Config.Kafka.AlgoTopic, // 统一放入正常队列
		Value:     sarama.ByteEncoder(bytes),
		Partition: 0, // 指定0分区，需要在初始化时指定sarama.NewManualPartitioner
	}
	_, _, err = l.svcCtx.SyncProducer.SendMessage(msg)
	if err != nil {
		l.Logger.Errorf("send message(%s) err=%s \n", msg.Value, err)
		return
	}
}

// SendChildOrderKafka 发送子单Kafka
func (l *OrigOrderAnalyseLogic) SendChildOrderKafka(batchNo int64, child models.TbAlgoOrderDetail) {
	b := &proto.ChildOrderPerf{
		Id:        uint32(child.ChildOrderId),
		BusUserId: child.UserId,
		//BusUuserId:     uint32(accountInfos[a.UserId]),
		AlgoOrderId:    uint32(child.AlgoOrderId),
		AlgorithmType:  uint32(child.AlgorithmType),
		AlgorithmId:    uint32(child.AlgorithmId),
		USecurityId:    uint32(child.UsecurityId),
		SecurityId:     child.SecurityId,
		Side:           uint32(child.TradeSide),
		OrderQty:       uint64(child.OrderQty) * 100,
		Price:          uint64(child.Price),
		OrderType:      uint32(child.OrderType),
		CumQty:         uint64(child.ComQty) * 100,
		LastPx:         uint64(child.LastPx),
		LastQty:        uint64(child.LastQty) * 100,
		Charge:         child.TotalFee,
		ArrivedPrice:   uint64(child.ArrivedPrice),
		ChildOrdStatus: uint32(child.OrdStatus),
		TransactTime:   uint64(child.TransactTime),
		Version:        0,
		BatchNo:        batchNo,
		BatchName:      "",
		SourceFrom:     2,
	}
	// 发送Kafka
	bytes, err := pt.Marshal(b)
	if err != nil {
		l.Logger.Error("Marshal msg error:", err)
		return
	}
	msg := &sarama.ProducerMessage{
		Topic:     l.svcCtx.Config.Kafka.ChildTopic, // 统一放入正常队列
		Value:     sarama.ByteEncoder(bytes),
		Partition: 0, // 指定0分区，需要在初始化时指定sarama.NewManualPartitioner
	}
	_, _, err = l.svcCtx.SyncProducer.SendMessage(msg)
	if err != nil {
		l.Logger.Errorf("send message(%s) err=%s \n", msg.Value, err)
		return
	}
}
