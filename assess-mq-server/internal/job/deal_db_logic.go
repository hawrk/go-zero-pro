// Package job
/*
 Author: hawrkchen
 Date: 2022/4/1 15:27
 Desc:  处理定时任务崩溃时，DB未处理的订单
*/
package job

import (
	"algo_assess/assess-mq-server/internal/logic"
	"algo_assess/global"
	"fmt"
	"github.com/spf13/cast"
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
	var datas []global.ChildOrderData
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
			data := global.ChildOrderData{
				OrderId:          detail.ChildOrderId,
				AlgoOrderId:      int64(detail.AlgoOrderId),
				AlgorithmType:    detail.AlgorithmType,
				AlgoId:           detail.AlgorithmId,
				UsecId:           detail.UsecurityId,
				SecId:            detail.SecurityId,
				OrderQty:         detail.OrderQty,
				Price:            detail.Price,
				OrderType:        detail.OrderType,
				LastPx:           detail.LastPx,
				LastQty:          detail.LastQty,
				ComQty:           detail.ComQty,
				ArrivePrice:      detail.ArrivedPrice,
				ChildOrderStatus: detail.OrdStatus,
				TransTime:        detail.TransactTime,
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
		algoKey := fmt.Sprintf("%d:%d:%s", data.TransTime, data.AlgoId, data.SecId)
		v := global.GlobalAssess.CalAlgo[algoKey]
		out, err := logic.RealTimeCal(o.s, v, &data)
		if err != nil {
			o.Logger.Error("error cal assess:", err)
			continue
		}
		global.GlobalAssess.CalAlgo[algoKey] = out
	}
	global.GlobalAssess.RWMutex.Unlock()

}
