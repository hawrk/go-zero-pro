package main

import (
	"algo_assess/assess-mq-server/internal/dao"
	"algo_assess/assess-mq-server/internal/job"
	"algo_assess/assess-mq-server/internal/listen"
	"flag"
	"fmt"

	"algo_assess/assess-mq-server/internal/config"
	"algo_assess/assess-mq-server/internal/server"
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/mqassess.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	svr := server.NewAssessMqServiceServer(svcCtx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		proto.RegisterAssessMqServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	//defer s.Stop()

	// 加载DB基础数据
	if err := dao.Run(c, svcCtx); err != nil {
		fmt.Println("load data error:", err)
		// 加载不到基础数据，直接panic掉，防止数据计算错误
		panic(err)
	}
	//dao.GScoreConf.GetProgressScore(30.0)

	// test   facade
	//dbTask := dao.NewLoadTask()
	//dbTask.StartLoadTask(svcCtx)

	// 加载二期内存计算数据
	if err := dao.Reload(c, svcCtx); err != nil {
		fmt.Println("reload db data error:", err)
		//
	}

	if err := dao.SetHeartBeat(svcCtx); err != nil {
		fmt.Println("set heart beat error:", err)
	}

	// 加载MQ
	// log、prometheus、trace、metricsUrl.
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range listen.Mqs(c, svcCtx) {
		serviceGroup.Add(mq)
	}
	// 加载rpc 服务
	serviceGroup.Add(s)

	done := make(chan struct{})
	go job.StartOrderJob(c, done)
	defer func() {
		done <- struct{}{}
	}()

	go job.StartAsyncJob(c, svcCtx)

	fmt.Printf("Starting assess mq rpc server at %s...\n", c.ListenOn)
	serviceGroup.Start()

	//s.Start()
}
