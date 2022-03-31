package main

import (
	"flag"
	"fmt"

	"algo_assess/rpc/assess/internal/config"
	"algo_assess/rpc/assess/internal/server"
	"algo_assess/rpc/assess/internal/svc"
	"algo_assess/rpc/assess/proto"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/assess.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	svr := server.NewAssessServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		proto.RegisterAssessServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	// add for zero mq
	//go mq.StartZeroMQ()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
