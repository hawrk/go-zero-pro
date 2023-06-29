package svc

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/mq-router-server/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	RedisClient     *redis.Redis
	AssessMQMClient mqservice.AssessMqService
	AssessMQSClient mqservice.AssessMqService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		RedisClient:     redis.New(c.Redis.Host),
		AssessMQMClient: mqservice.NewAssessMqService(zrpc.MustNewClient(c.AssessMQMRPC)),
		AssessMQSClient: mqservice.NewAssessMqService(zrpc.MustNewClient(c.AssessMQSRPC)),
	}
}
