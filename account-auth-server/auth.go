package main

import (
	"account-auth/account-auth-server/internal/config"
	"account-auth/account-auth-server/internal/handler"
	"account-auth/account-auth-server/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

var configFile = flag.String("f", "etc/auth-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	//server := rest.MustNewServer(c.RestConf)
	//server := rest.MustNewServer(c.RestConf, rest.WithCors())

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(
		func(header http.Header) {
			header.Set("Access-Control-Allow-Origin", "*")                                                   // 这是允许访问所有域
			header.Add("Access-Control-Allow-Headers", "x-requested-with,x-token,Access-Token,Content-Type") // 这是允许访问所有域
			header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			//header.Add("Access-Control-Expose-Headers", "Content-Disposition")
		}, nil))
	//server := rest.MustNewServer(c.RestConf,
	//	rest.WithNotAllowedHandler(middleware.NewCorsMiddleware().Handler()))
	//server.Use(middleware.NewCorsMiddleware().Handle)

	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
