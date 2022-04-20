package svc

import (
	"algo_assess/assess-api-server/internal/config"
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/assess-rpc-server/assessservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AssessClient   assessservice.AssessService
	AssessMQClient mqservice.AssessMqService
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		AssessClient:   assessservice.NewAssessService(zrpc.MustNewClient(c.AssessRPC)),
		AssessMQClient: mqservice.NewAssessMqService(zrpc.MustNewClient(c.AssessMQRPC)),
	}
}
