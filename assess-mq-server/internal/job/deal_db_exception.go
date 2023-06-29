// Package job
/*
 Author: hawrkchen
 Date: 2022/11/29 11:23
 Desc:
*/
package job

import (
	"algo_assess/global"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
)

func (o *AssessJob) DealDBException() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("in DealDBException recover :", err)
		}
	}()
	o.Logger.Info("DealDBException start...")
	//time.Sleep(time.Second * 60)
	// 取异常队列下面的所有Key
	// 1. 取profile 队列数据
	o.DoProfileWork()
	//2. 取 profit 队列数据
	o.DoProfileSumWork()
	//3. 取timeline 队列数据
	o.DoTimeLineWork()

	o.Logger.Info("DealDBException end...")
}

func (o *AssessJob) DoProfileWork() {
	data, err := o.s.RedisClient.Hgetall(global.FailProfileKey)
	if err != nil {
		o.Logger.Error("hgetall error:", err)
		return
	}
	for k, v := range data {
		var pf global.Profile
		if err := json.Unmarshal([]byte(v), &pf); err != nil {
			o.Logger.Error("Unmarshal error:", err)
			continue
		}
		record, err := o.s.ProfileRepo.GetDataByProfileKey(o.ctx, pf.Date, pf.AccountId, pf.AlgoId, pf.SecId, pf.AlgoOrderId)
		if err != nil {
			o.Logger.Error("GetDataByProfileKey error:", err)
			continue
		}
		if record.AccountId == "" { // DB表无记录，则直接插入
			if err := o.s.ProfileRepo.CreateAlgoProfile(o.ctx, &pf); err != nil {
				o.Logger.Error("CreateAlgoProfile error:", err)
				continue
			}
			// TODO:
			//if err := o.s.SummaryRepo.CreateAlgoSummary(o.ctx, &profile); err != nil {
			//	o.Logger.Error("CreateAlgoSummary error:", err)
			//	continue
			//}
			_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailProfileKey, k)
			if err != nil {
				o.Logger.Error("hdel key error:", err)
			}
			continue
		}
		// 更新时间与交易时间比较
		uTime := cast.ToInt64(record.UpdateTime.Format(global.TimeFormatMinInt))
		o.Logger.Info("get db update Time:", uTime, ",redis trans time:", pf.TransAt)
		if uTime > pf.TransAt { // DB更新时间 大于 交易时间时，表示该交易已更新，不需要处理,直接清除该记录
			o.Logger.Info("update time greater than trans time, continue")
			_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailProfileKey, k)
			if err != nil {
				o.Logger.Error("hdel key error:", err)
			}
			continue
		}
		// 否则更新DB的数据
		if err := o.s.ProfileRepo.UpdateAlgoProfile(o.ctx, &pf); err != nil {
			o.Logger.Error("UpdateAlgoProfile error:", err)
			continue
		}
		//TODO:
		//if err := o.s.SummaryRepo.UpdateAlgoSummary(o.ctx, &profile); err != nil {
		//	o.Logger.Error("UpdateAlgoSummary error:", err)
		//	continue
		//}
		// 删除异常队列中该笔记录
		_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailProfileKey, k)
		if err != nil {
			o.Logger.Error("hdel key error:", err)
		}
	}
}

func (o *AssessJob) DoProfileSumWork() {
	data, err := o.s.RedisClient.Hgetall(global.FailProfileSumKey)
	if err != nil {
		o.Logger.Error("hgetall error:", err)
		return
	}
	for k, v := range data {
		var ps global.ProfileSum
		if err := json.Unmarshal([]byte(v), &ps); err != nil {
			o.Logger.Error("Unmarshal error:", err)
			continue
		}
		record, err := o.s.SummaryRepo.GetDataBySummaryKey(o.ctx, ps.Date, ps.AccountId, ps.AlgoId)
		if err != nil {
			o.Logger.Error("GetDataBySummaryKey error:", err)
			continue
		}

		if record.UserId == "" { // DB表无记录，则直接插入
			if err := o.s.SummaryRepo.CreateAlgoSummary(o.ctx, &ps); err != nil {
				o.Logger.Error("CreateAlgoSummary error:", err)
				continue
			}
			_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailProfileSumKey, k)
			if err != nil {
				o.Logger.Error("hdel key error:", err)
			}
			continue
		}
		// 更新时间与交易时间比较
		uTime := cast.ToInt64(record.UpdateTime.Format(global.TimeFormatMinInt))
		if uTime >= ps.TransAt { // DB更新时间 大于 交易时间时，表示该交易已更新，不需要处理
			o.Logger.Info("update time greater than trans time, continue")
			_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailProfileSumKey, k)
			if err != nil {
				o.Logger.Error("hdel key error:", err)
			}
			continue
		}
		// 否则更新DB的数据
		if err := o.s.SummaryRepo.UpdateAlgoSummary(o.ctx, &ps); err != nil {
			o.Logger.Error("UpdateAlgoProfit error:", err)
			continue
		}
		// 删除异常队列中该笔记录
		_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailProfileSumKey, k)
		if err != nil {
			o.Logger.Error("hdel key error:", err)
		}
	}
}

func (o *AssessJob) DoTimeLineWork() {
	data, err := o.s.RedisClient.Hgetall(global.FailTimeLineKey)
	if err != nil {
		o.Logger.Error("hgetall error:", err)
		return
	}
	for k, v := range data {
		var tl global.ProfileSum
		if err := json.Unmarshal([]byte(v), &tl); err != nil {
			o.Logger.Error("Unmarshal error:", err)
			continue
		}
		record, err := o.s.TimeLineRepo.GetDataByTimeLineKey(o.ctx, tl.AccountId, tl.TransAt, tl.AlgoId)
		if err != nil {
			o.Logger.Error("GetDataByProfileKey error:", err)
			continue
		}
		if record.AccountId == "" { // DB表无记录，则直接插入
			if err := o.s.TimeLineRepo.CreateAlgoTimeLine(o.ctx, &tl); err != nil {
				o.Logger.Error("CreateAlgoTimeLine error:", err)
				continue
			}
			_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailTimeLineKey, k)
			if err != nil {
				o.Logger.Error("hdel key error:", err)
			}
			continue
		}
		// 更新时间与交易时间比较
		if record.TransactTime >= tl.TransAt { // DB更新时间 大于 交易时间时，表示该交易已更新，不需要处理
			o.Logger.Info("update time greater than trans time, continue")
			_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailTimeLineKey, k)
			if err != nil {
				o.Logger.Error("hdel key error:", err)
			}
			continue
		}
		// 否则更新DB的数据
		if err := o.s.TimeLineRepo.UpdateAlgoTimeLine(o.ctx, &tl); err != nil {
			o.Logger.Error("UpdateAlgoTimeLine error:", err)
			continue
		}
		// 删除异常队列中该笔记录
		_, err = o.s.RedisClient.HdelCtx(o.ctx, global.FailTimeLineKey, k)
		if err != nil {
			o.Logger.Error("hdel key error:", err)
		}
	}
}
