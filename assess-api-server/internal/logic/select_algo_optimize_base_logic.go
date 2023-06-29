package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectAlgoOptimizeBaseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectAlgoOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectAlgoOptimizeBaseLogic {
	return &SelectAlgoOptimizeBaseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectAlgoOptimizeBaseLogic) SelectAlgoOptimizeBase(req *types.SelectOptimizeBaseReq) (resp *types.SelectOptimizeBaseRsp, err error) {
	rsp, err := l.svcCtx.AssessClient.SelectOptimizeBase(l.ctx, &assessservice.SelectOptimizeBaseReq{
		ProviderId: req.ProviderId,
		SecId:      req.SecId,
		AlgoId:     req.AlgoId,
		Page:       req.Page,
		Limit:      req.Limit,
	})
	p := &types.SelectOptimizeBaseRsp{
		Code:  int(rsp.GetCode()),
		Msg:   rsp.GetMsg(),
		Total: rsp.GetTotal(),
		Data:  ToOptimizeRsp(rsp.List),
	}
	return p, err
}

func ToOptimizeRsp(data []*assessservice.OptimizeBase) (ret []types.OptimizeBase) {
	for _, v := range data {
		o := types.OptimizeBase{
			Id:           v.GetId(),
			ProviderId:   v.GetProviderId(),
			ProviderName: v.GetProviderName(),
			SecId:        v.GetSecId(),
			SecName:      v.GetSecName(),
			AlgoId:       v.GetAlgoId(),
			AlgoType:     v.GetAlgoType(),
			AlgoName:     v.GetAlgoName(),
			OpenRate:     v.GetOpenRate(),
			IncomeRate:   v.GetIncomeRate(),
			BasisPoint:   v.GetBasisPoint(),
			CreateTime:   v.CreateTime,
			UpdateTime:   v.GetUpdateTime(),
		}
		ret = append(ret, o)
	}
	return ret
}
