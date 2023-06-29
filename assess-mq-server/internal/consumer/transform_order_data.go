// Package consumer
/*
 Author: hawrkchen
 Date: 2022/5/11 14:13
 Desc:
*/
package consumer

import (
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"
)

// TransAlgoOrderData 总线母单消息推送结构体转换
func TransAlgoOrderData(data *pb.AlgoOrderPerf) global.MAlgoOrder {
	t := time.UnixMicro(int64(data.GetTransactTime())).Format(global.TimeFormatMinInt)
	transactAt := tools.TimeMoveForward(t)
	var batchNo int64
	if data.BatchNo == 0 { // 总线数据为0，默认置为1
		batchNo = global.DefaultBatchNo
	} else {
		batchNo = data.BatchNo
	}
	// 如果是日内回转的，同一个母单有买卖两个子单，委托数量要*2
	var orderQty int64
	if data.AlgorithmType == global.AlgoTypeT0 {
		orderQty = int64(data.AlgoOrderQty) * 2
	} else {
		orderQty = int64(data.AlgoOrderQty)
	}
	out := global.MAlgoOrder{
		BatchNo:         batchNo,
		BasketId:        int(data.BasketId),
		UserId:          strings.TrimSpace(tools.RMu0000(data.BusUserId)),
		AlgoId:          int(data.Id),          // 母单单号
		AlgorithmId:     int(data.AlgorithmId), // 母单算法ID
		AlgorithmType:   int(data.AlgorithmType),
		UsecId:          int(data.USecurityId),
		SecId:           strings.TrimSpace(tools.RMu0000(data.SecurityId)),
		AlgoOrderQty:    orderQty / 100, // 该母单委托数量，非成交数量
		TransTime:       transactAt,
		StartTime:       int64(data.StartTime),
		SourceFrom:      int(data.SourceFrom),
		EndTime:         int64(data.EndTime),
		UnixTime:        cast.ToString(data.TransactTime)[:10],
		UnixTimeMillSec: int64(data.TransactTime),
		SourcePrx:       GetSourcePrx(data.SourceFrom, data.BatchNo),
	}
	return out
}

func TransChildOrderData(data *pb.ChildOrderPerf) global.ChildOrderData {
	t := time.UnixMicro(int64(data.GetTransactTime())).Format(global.TimeFormatMinInt)
	transactAt := tools.TimeMoveForward(t)
	var batchNo int64
	if data.BatchNo == 0 { // 总线数据为0，默认置为1
		batchNo = global.DefaultBatchNo
	} else {
		batchNo = data.BatchNo
	}
	out := global.ChildOrderData{
		BatchNo:          batchNo,
		OrderId:          int64(data.Id),
		AlgoOrderId:      int64(data.AlgoOrderId),
		AlgorithmType:    int(data.AlgorithmType),
		AlgoId:           int(data.AlgorithmId),
		UsecId:           uint(data.USecurityId),
		UserId:           strings.TrimSpace(tools.RMu0000(data.BusUserId)),
		SecId:            strings.TrimSpace(data.SecurityId),
		TradeSide:        tools.GetOrderTradeSide(data.Side), // 买卖方向    1-买    2-卖
		OrderQty:         int64(data.OrderQty) / 100,
		Price:            int64(data.Price), // 委托价格， 价格先不转
		OrderType:        uint(data.OrderType),
		LastPx:           int64(data.LastPx), // 成交价格
		LastQty:          int64(data.LastQty) / 100,
		ComQty:           int64(data.CumQty) / 100,
		ArrivePrice:      int64(data.ArrivedPrice),  // 到达价格
		TotalFee:         cast.ToInt64(data.Charge), // 总手续费
		ChildOrderStatus: uint(data.ChildOrdStatus),
		TransTime:        transactAt,
		CurDate:          cast.ToInt64(strconv.FormatInt(transactAt, 10)[:8]),
		UnixTime:         cast.ToInt64(cast.ToString(data.TransactTime)[:10]),
		UnixTimeMillSec:  int64(data.TransactTime),
		SourceFrom:       int(data.SourceFrom),
		SourcePrx:        GetSourcePrx(data.SourceFrom, data.BatchNo),
	}
	return out
}

func GetSourcePrx(s int32, batchNo int64) string {
	switch s {
	case global.SourceFromBus:
		return global.OrderSourceNor + ":1"
	case global.SourceFromFix:
		return global.OrderSourceFix + ":1"
	case global.SourceFromOrigin: // 原始订单
		return global.OrderSourceOri + ":" + cast.ToString(batchNo)
	//case global.SourceFromImport:      // 订单导入与原始订单合并
	//	return global.OrderSourceImp + ":" + cast.ToString(batchNo)
	default:
		return global.OrderSourceAbt
	}
}
