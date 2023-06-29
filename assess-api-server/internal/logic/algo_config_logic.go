package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlgoConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoConfigLogic {
	return &AlgoConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlgoConfigLogic) AlgoConfig(req *types.AlgoConfigReq) (resp *types.AlgoConfigRsp, err error) {
	l.Logger.Infof("in AlgoConfig get Req:%+v", *req)
	aReq := &mqservice.AlgoConfigReq{
		ProfileType: req.ProfileType,
		AlgoConfig:  req.ConfigJson,
	}
	rsp, err := l.svcCtx.AssessMQClient.AlgoConfig(l.ctx, aReq)
	if err != nil {
		l.Logger.Error("rpc AlgoConfig error:", err)
		return &types.AlgoConfigRsp{
			Code:   360,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}

	return &types.AlgoConfigRsp{
		Code:   int(rsp.GetCode()),
		Msg:    rsp.GetMsg(),
		Result: rsp.GetResult(),
	}, nil
}
