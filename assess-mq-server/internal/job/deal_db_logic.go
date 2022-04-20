// Package job
/*
 Author: hawrkchen
 Date: 2022/4/1 15:27
 Desc:  处理定时任务崩溃时，DB未处理的订单
*/
package job

import (
	"algo_assess/assess-mq-server/internal/logic"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

func (o *AssessJob) DealDBAssess() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("in DealDBAssess recover :", err)
		}
	}()
	o.Logger.Info("in DealDBAssess:deal incomplete orders")
	// check redis 是否有未处理的订单
	arr, err := o.s.RedisClient.Smembers(global.AssessTimeSetKey)
	if err != nil {
		o.Logger.Error("read redis set key error,err:", err)
		return
	}
	if len(arr) == 0 {
		o.Logger.Info(" no task to do, finish...")
		return
	}
	var datas []pb.ChildOrderPerf
	for _, v := range arr {
		transactAt := cast.ToInt64(v)
		result, err := o.s.OrderDetailRepo.QueryOrderDetail(transactAt)
		if err != nil {
			o.Logger.Error("query order detail error:", err)
			continue
		}
		for _, detail := range result {
			if detail.ProcStatus == 1 { // 已处理的不再处理
				continue
			}
			data := pb.ChildOrderPerf{
				Id:             uint32(detail.ChildOrderId),
				AlgoOrderId:    uint32(detail.AlgoOrderId),
				AlgorithmType:  uint32(detail.AlgorithmType),
				AlgorithmId:    uint32(detail.AlgorithmId),
				USecurityId:    uint32(detail.UsecurityId),
				SecurityId:     detail.SecurityId,
				OrderQty:       uint64(detail.OrderQty),
				Price:          uint64(detail.Price),
				OrderType:      uint32(detail.OrderType),
				CumQty:         uint64(detail.ComQty),
				LastPx:         uint64(detail.LastPx),
				LastQty:        uint64(detail.LastQty),
				ArrivedPrice:   uint64(detail.ArrivedPrice),
				ChildOrdStatus: uint32(detail.OrdStatus),
				TransactTime:   uint64(detail.TransactTime),
			}
			datas = append(datas, data)
		}
	}
	if len(datas) == 0 {
		o.Logger.Info("no task to process,finish")
		return
	}
	// 合并到本地缓存中
	global.GlobalAssess.RWMutex.Lock()
	for _, data := range datas {
		transactAt := time.UnixMicro(int64(data.GetTransactTime())).Format(global.TimeFormatMinInt)
		//transact := cast.ToInt64(transactAt)
		algoId := fmt.Sprintf("%s:%d:%d", transactAt, data.AlgorithmId, data.USecurityId)
		v := global.GlobalAssess.CalAlgo[algoId]
		out, err := logic.RealTimeCal(transactAt, v, &data)
		if err != nil {
			o.Logger.Error("error cal assess:", err)
			continue
		}
		global.GlobalAssess.CalAlgo[algoId] = out
	}
	global.GlobalAssess.RWMutex.Unlock()

}
