package main

import (
	"algo_assess/assess-api-server/internal/config"
	"algo_assess/assess-api-server/internal/handler"
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/pkg/zookeeper"
	"flag"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/assess-api.yaml", "the config file")
var zk *zookeeper.Zookeeper

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	//server := rest.MustNewServer(c.RestConf, rest.WithCors(),rest.WithUnauthorizedCallback(middleware.NewInterceptorMiddleware().Handle()))
	//s := rest.MustNewServer(c.RestConf, rest.WithCors())
	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(
		func(header http.Header) {
			header.Set("Access-Control-Allow-Origin", "*")                              // 这是允许访问所有域
			header.Add("Access-Control-Allow-Headers", "x-requested-with,content-type") // 这是允许访问所有域
			header.Add("Access-Control-Expose-Headers", "Content-Disposition")
		}, nil))
	//server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)
	// 注册zk
	if c.Zookeeper.EnableRegister {
		RegisterZookeeper(c)
		defer zk.Conn.Close()
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	server.Start()
}

func RegisterZookeeper(c config.Config) {
	zk = zookeeper.NewZookeeper(c.Zookeeper.Host, c.Zookeeper.TimeOut)
	if err := zk.Connect(); err != nil {
		fmt.Println("connect zookeeper fail:")
		return
	}
	//serverHost := fmt.Sprintf("%s:%d", c.Host, c.Port)
	var serverHost string
	if c.Domain.UseDomain {
		serverHost = c.Domain.DomainName
	} else {
		serverHost = fmt.Sprintf("%s:%d", c.Host, c.Port)
	}
	path := "/assess-api-server"
	if err := zk.RegisterServer(path, serverHost); err != nil {
		fmt.Println("register server fail：", err)
		return
	}
	fmt.Println("register zookeeper success:", serverHost)
}
