// Package mqueue
/*
 Author: hawrkchen
 Date: 2022/3/24 9:59
 Desc:
*/
package main

import (
	"algo_assess/mqueue/internal/config"
	"algo_assess/mqueue/internal/job"
	"algo_assess/mqueue/internal/listen"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var configFile = flag.String("f", "etc/mq.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)

	// log、prometheus、trace、metricsUrl.
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	// 初始化全局计算缓存变量

	// 禁止打印状态日志
	logx.DisableStat()
	for _, mq := range listen.Mqs(c) {
		serviceGroup.Add(mq)
	}

	done := make(chan struct{})
	go job.StartOrderJob(c, done)
	defer func() {
		done <- struct{}{}
	}()


	serviceGroup.Start()
}
