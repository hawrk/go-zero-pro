// Package job
/*
 Author: hawrkchen
 Date: 2022/5/24 18:59
 Desc:
*/
package job

import (
	"algo_assess/market-mq-server/internal/config"
	"fmt"
	"github.com/robfig/cron/v3"
)

func newCronCnf() *cron.Cron {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(parser), cron.WithChain())
}

func StartMarketJob(c config.Config, d chan struct{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("in StartOrderJob recover :", err)
		}
	}()
	// job
	m := NewMarketJob(c)
	job := newCronCnf()
	spec := "0 10 15 * * ?" //
	_, err := job.AddFunc(spec, func() {
		//m.SayHello()
		m.SyncRedisToDb()
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
