// Package job
/*
 Author: hawrkchen
 Date: 2023/6/27 10:04
 Desc:
*/
package job

import (
	"algo_assess/global"
	"algo_assess/mq-router-server/internal/config"
	"algo_assess/mq-router-server/internal/svc"
	"algo_assess/pkg/tools"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logx"
)

var ConsisHash *tools.HashRing
var NodeNum int

func newCronCnf() *cron.Cron {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(parser), cron.WithChain())
}

func StartAsyncNodes(c config.Config, svcContext *svc.ServiceContext) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("in StartOrderJob recover :", err)
		}
	}()

	ConsisHash = tools.NewHashRing(20, nil)
	// 服务启动，先从Redis中取列表初始化节点
	str, err := svcContext.RedisClient.Smembers(global.ConsisHasNode)
	fmt.Println("get node list:", str)
	ConsisHash.AddNodes(str...)
	NodeNum = len(str)

	job := newCronCnf()
	spec := "0 */1 * * * ?" // 每一分钟从Redis更新节点

	_, err = job.AddFunc(spec, func() {
		AsyncNodes(c, svcContext)
	})
	if err != nil {
		fmt.Println("add cron error:", err)
		return
	}

	job.Start()
	defer job.Stop()

	select {}
}

// AsyncNodes 同步Redis节点信息
func AsyncNodes(c config.Config, svcContext *svc.ServiceContext) {
	str, err := svcContext.RedisClient.Smembers(global.ConsisHasNode)
	if err != nil {
		logx.Error("Smembers error:", err)
		return
	}
	// redis取到的节点个数与 当前保存的个数一致时，则跳过处理
	//logx.Info("redis nodes num:", len(str), ",curr nodes num:", NodeNum)
	if len(str) == NodeNum {
		logx.Info("same nodes, continue...")
		return
	}
	logx.Info("nodes change, reset ....")
	ConsisHash.Reset(str...)
	// 同步当前节点个数
	NodeNum = len(str)
}
