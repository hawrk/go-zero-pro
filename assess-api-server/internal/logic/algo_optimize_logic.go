package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"algo_assess/assess-rpc-server/assessservice"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type AlgoOptimizeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoOptimizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoOptimizeLogic {
	return &AlgoOptimizeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlgoOptimizeLogic) AlgoOptimize(req *types.OptimizeReq) (resp *types.OptimizeRsp, err error) {
	optimizeReq := &assessservice.OptimizeReq{
		SecurityId: req.SecurityId,
		AlgoIds:    req.AlgoIds,
	}
	rsp, err := l.svcCtx.AssessClient.GetOptimize(l.ctx, optimizeReq)
	if err != nil {
		return nil, err
	}
	p := &types.OptimizeRsp{
		Code:  int(rsp.GetCode()),
		Msg:   rsp.GetMsg(),
		Total: rsp.GetTotal(),
		Data:  GetOptimize(rsp.List),
	}
	return p, nil
}

func GetOptimize(data []*assessservice.Optimize) (ret []types.OptimizeInfo) {
	for _, v := range data {
		optimize := types.OptimizeInfo{
			Id:            v.Id,
			Provider_id:   int(v.ProviderId),
			Provider_name: v.ProviderName,
			SecId:         v.SecId,
			SecName:       v.SecName,
			AlgoId:        int(v.AlgoId),
			AlgoName:      v.AlgoName,
		}
		ret = append(ret, optimize)
	}
	return ret
}
