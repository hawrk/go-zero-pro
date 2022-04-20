package svc

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/market-mq-server/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	RedisClient    *redis.Redis
	AssessMQClient mqservice.AssessMqService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		RedisClient:    redis.New(c.Redis.Host),
		AssessMQClient: mqservice.NewAssessMqService(zrpc.MustNewClient(c.AssessMQRPC)),
	}
}
