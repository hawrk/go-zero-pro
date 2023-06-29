package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAlgoOptimizeBaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAlgoOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAlgoOptimizeBaseLogic {
	return &AddAlgoOptimizeBaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAlgoOptimizeBaseLogic) AddAlgoOptimizeBase(req *types.AddOptimizeBaseReq) (resp *types.OptimizeBaseRsp, err error) {
	rsp, err := l.svcCtx.AssessClient.AddOptimizeBase(l.ctx, &assessservice.AddOptimizeBaseReq{
		ProviderId:   req.ProviderId,
		ProviderName: req.ProviderName,
		SecId:        req.SecId,
		SecName:      req.SecName,
		AlgoId:       req.AlgoId,
		AlgoType:     req.AlgoType,
		AlgoName:     req.AlgoName,
		OpenRate:     req.OpenRate,
		IncomeRate:   req.IncomeRate,
		BasisPoint:   req.BasisPoint,
	})
	p := &types.OptimizeBaseRsp{
		Code: int(rsp.GetCode()),
		Msg:  rsp.GetMsg(),
	}
	return p, err
}
