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
	"github.com/spf13/cast"
	"strings"
	"time"
)

// TransAlgoOrderData 总线母单消息推送结构体转换
func TransAlgoOrderData(data *pb.AlgoOrderPerf) global.MAlgoOrder {
	transactAt := time.UnixMicro(int64(data.GetTransactTime())).Format(global.TimeFormatMinInt)
	out := global.MAlgoOrder{
		AlgoId:       int(data.Id),
		AlgorithmId:  int(data.AlgorithmId),
		UsecId:       int(data.USecurityId),
		SecId:        strings.TrimSpace(data.SecurityId),
		AlgoOrderQty: int64(data.AlgoOrderQty) / 100,
		TransTime:    cast.ToInt64(transactAt),
	}
	return out
}

// TransChildOrderData 总线子单消息推送结构体转换
func TransChildOrderData(data *pb.ChildOrderPerf) global.ChildOrderData {
	t := time.UnixMicro(int64(data.GetTransactTime())).Format(global.TimeFormatMinInt)
	transactAt := cast.ToInt64(t)
	out := global.ChildOrderData{
		OrderId:          int64(data.Id),
		AlgoOrderId:      int64(data.AlgoOrderId),
		AlgorithmType:    uint(data.AlgorithmType),
		AlgoId:           uint(data.AlgorithmId),
		UsecId:           uint(data.USecurityId),
		SecId:            strings.TrimSpace(data.SecurityId),
		OrderQty:         int64(data.OrderQty) / 100,
		Price:            int64(data.Price) / 100, // 总线过来的价格是以元为单位再乘以10000，所以这里要转成分只需要除以100就可以了
		OrderType:        uint(data.OrderType),
		LastPx:           int64(data.LastPx) / 100,
		LastQty:          int64(data.LastQty) / 100,
		ComQty:           int64(data.CumQty) / 100,
		ArrivePrice:      int64(data.ArrivedPrice) / 100,
		ChildOrderStatus: uint(data.ChildOrdStatus),
		TransTime:        transactAt,
	}
	return out
}
