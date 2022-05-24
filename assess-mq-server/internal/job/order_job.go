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
	spec := "0 */1 * * * ?"
	_, err := job.AddFunc(spec, func() {
		//order.GetAssessWithCalcute()
		order.GetAssessWithCache()
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
