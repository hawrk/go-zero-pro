// Package job
/*
 Author: hawrkchen
 Date: 2022/4/1 15:27
 Desc:  处理定时任务崩溃时，DB未处理的订单
*/
package job

import (
	"algo_assess/assess-mq-server/internal/consumer"
	"algo_assess/assess-mq-server/internal/dao"
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
				AlgorithmType:    int(detail.AlgorithmType),
				AlgoId:           int(detail.AlgorithmId),
				UserId:           detail.UserId,
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

	s := consumer.NewAlgoPlatformOrderTrade(o.ctx, o.s)
	// 合并到本地缓存中
	for _, data := range datas {
		var virtualUser, normalUser string
		var adminUser []string
		virtualUser = "0"
		dao.GAccountMap.RWMutex.RLock()
		if _, exist := dao.GAccountMap.Account[data.UserId]; exist {
			if dao.GAccountMap.Account[data.UserId].UserType == 1 { // 普通账户，还需要找到其管理员账户
				normalUser = data.UserId
				adminUser = dao.GAccountMap.Account[data.UserId].ParUserId
			} else if dao.GAccountMap.Account[data.UserId].UserType == 3 { // 管理员账户，直接取其user_id
				adminUser = []string{data.UserId}
			} else if dao.GAccountMap.Account[data.UserId].UserType == 2 { // 算法厂商用户
				virtualUser = data.UserId
			}
		} else {
			o.Logger.Error(" userId not found:", data.UserId)
			// 本地缓存找不到账户信息时，只能计算所有用户的汇总绩效
		}
		dao.GAccountMap.RWMutex.RUnlock()
		if normalUser != "" {
			algoKey := fmt.Sprintf("%d:%s:%d:%s", data.TransTime, normalUser, data.AlgoId, data.SecId)
			s.Dispatch(normalUser, algoKey, &data)
		}
		for _, v := range adminUser {
			algoKey := fmt.Sprintf("%d:%s:%d:%s", data.TransTime, v, data.AlgoId, data.SecId)
			s.Dispatch(v, algoKey, &data)
		}
		if virtualUser != "" {
			algoKey := fmt.Sprintf("%d:%s:%d:%s", data.TransTime, virtualUser, data.AlgoId, data.SecId)
			s.Dispatch(virtualUser, algoKey, &data)
		}
	}

}
