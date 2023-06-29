package main

import (
	"algo_assess/market-mq-server/internal/listen"
	"algo_assess/market-mq-server/job"
	"flag"
	"fmt"

	"algo_assess/market-mq-server/internal/config"
	"algo_assess/market-mq-server/internal/server"
	"algo_assess/market-mq-server/internal/svc"
	"algo_assess/market-mq-server/proto"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/market.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	svcCtx := svc.NewServiceContext(c)
	svr := server.NewMarketServiceServer(svcCtx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		proto.RegisterMarketServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// 加载MQ
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("mq recover:", err)
			}
		}()
		// log、prometheus、trace、metricsUrl.
		if err := c.SetUp(); err != nil {
			panic(err)
		}

		serviceGroup := service.NewServiceGroup()
		defer serviceGroup.Stop()

		for _, mq := range listen.Mqs(c, svcCtx) {
			serviceGroup.Add(mq)
		}
		serviceGroup.Start()
	}()

	done := make(chan struct{})
	go job.StartMarketJob(c, done)
	defer func() {
		done <- struct{}{}
	}()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
