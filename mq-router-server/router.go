package main

import (
	"algo_assess/global"
	"algo_assess/mq-router-server/internal/job"
	"algo_assess/mq-router-server/internal/listen"
	"algo_assess/mq-router-server/ordproto"
	"flag"
	"fmt"

	"algo_assess/mq-router-server/internal/config"
	"algo_assess/mq-router-server/internal/server"
	"algo_assess/mq-router-server/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/router.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewRouterServiceServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ordproto.RegisterRouterServiceServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	// log、prometheus、trace、metricsUrl.
	if err := c.SetUp(); err != nil {
		panic(err)
	}

	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	for _, mq := range listen.Mqs(c, ctx) {
		serviceGroup.Add(mq)
	}
	// 加载rpc 服务
	serviceGroup.Add(s)
	// 注册自身服务节点到Redis，作动态扩容
	RegisterNode(c, ctx)

	// 定时更新mq服务节点
	go job.StartAsyncNodes(c, ctx)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	serviceGroup.Start()
	//s.Start()
}

// RegisterNode 注册自身服务节点--assess-mq-server
// TODO:
func RegisterNode(c config.Config, svcCtx *svc.ServiceContext) {
	_, err := svcCtx.RedisClient.Sadd(global.ConsisHasNode, c.ListenOn)
	if err != nil {
		fmt.Println("Sadd key error:", err)
	}
}
