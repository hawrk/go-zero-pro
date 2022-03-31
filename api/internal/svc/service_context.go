package svc

import (
	"algo_assess/api/internal/config"
	"algo_assess/rpc/assess/assess"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config       config.Config
	AssessClient assess.Assess
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		AssessClient: assess.NewAssess(zrpc.MustNewClient(c.AssessRPC)),
	}
}
