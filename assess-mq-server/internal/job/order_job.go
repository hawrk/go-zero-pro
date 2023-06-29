// Package job
/*
 Author: hawrkchen
 Date: 2022/3/29 15:21
 Desc:
*/
package job

import (
	"algo_assess/assess-mq-server/internal/config"
	"fmt"
	"github.com/robfig/cron/v3"
)

func newCronCnf() *cron.Cron {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(parser), cron.WithChain())
}

func StartOrderJob(c config.Config, d chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("in StartOrderJob recover :", err)
		}
	}()
	order := NewOrderJob(c)

	//go func() {
	//	order.DealDBAssess()
	//}()

	job := newCronCnf()
	spec := "0 */1 * * * ?"          // 每一分钟落地一期绩效数据
	spec2 := "0 0 20 * * ?"          // 每天晚上20点清除缓存数据
	exceptionSpec := "0 */5 * * * ?" // 每五分钟扫描异常队列中的数据
	heatBeatSpec := "*/30 * * * ?"   // 每30秒同步一次心跳数据
	// 一期实时计算数据落DB
	_, err := job.AddFunc(spec, func() {
		//order.GetAssessWithCalcute()
		order.GetAssessWithCache()
	})
	// 清除二期本地缓存数据
	_, err = job.AddFunc(spec2, func() {
		order.ClearLocalCache()
	})
	// 二期 异常处理落表失败,定时Job
	_, err = job.AddFunc(exceptionSpec, func() {
		order.DealDBException()
	})

	_, err = job.AddFunc(heatBeatSpec, func() {
		order.KeepHeartBeat()
	})
	if err != nil {
		fmt.Println("add cron error:", err)
		return
	}
	job.Start()
	defer job.Stop()

	for {
		select {
		case <-d:
			return
		}
	}
}
