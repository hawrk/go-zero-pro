package svc

import (
	"algo_assess/assess-api-server/internal/config"
	"algo_assess/assess-api-server/internal/middleware"
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/assess-rpc-server/assessservice"
	mkservice "algo_assess/market-mq-server/marketservice"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AssessClient   assessservice.AssessService
	AssessMQClient mqservice.AssessMqService
	MarketMQClient mkservice.MarketService
	Interceptor    rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AssessClient:   assessservice.NewAssessService(zrpc.MustNewClient(c.AssessRPC)),
		AssessMQClient: mqservice.NewAssessMqService(zrpc.MustNewClient(c.AssessMQRPC)),
		MarketMQClient: mkservice.NewMarketService(zrpc.MustNewClient(c.MarketMQRPC)),
		Interceptor:    middleware.NewInterceptorMiddleware().Handle,
	}
}
