package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/models"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetOptimizeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOptimizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOptimizeLogic {
	return &GetOptimizeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOptimize 一键优选
func (l *GetOptimizeLogic) GetOptimize(in *proto.OptimizeReq) (*proto.OptimizeRsp, error) {
	//time.Sleep(time.Second * 20)
	l.Logger.Infof("get req", in)
	var result []*models.TbAlgoOptimizeBase
	var err error
	var count int64
	//count, result, err = l.svcCtx.OptimizeRepo.GetAlgoOptimizeBySecurityId(in.GetSecurityId())
	count, result, err = l.svcCtx.OptimizeBaseRepo.GetOptimize(in.GetSecurityId(), in.GetAlgoIds())
	l.Logger.Infof("get result:%+v", result)
	if err != nil {
		return &proto.OptimizeRsp{
			Code:  1000,
			Msg:   "查询失败",
			Total: 0,
			List:  nil,
		}, err
	}
	return &proto.OptimizeRsp{
		Code:  0,
		Msg:   "查询成功",
		Total: count,
		List:  GetOptimize(result),
	}, nil
}

func GetOptimize(result []*models.TbAlgoOptimizeBase) (ret []*proto.Optimize) {
	for _, optimizeInfo := range result {
		optimize := &proto.Optimize{
			Id:           optimizeInfo.Id,
			ProviderId:   int32(optimizeInfo.ProviderId),
			ProviderName: optimizeInfo.ProviderName,
			SecId:        optimizeInfo.SecId,
			SecName:      optimizeInfo.SecName,
			AlgoId:       int32(optimizeInfo.AlgoId),
			AlgoName:     optimizeInfo.AlgoName,
		}
		ret = append(ret, optimize)
	}
	return ret
}
