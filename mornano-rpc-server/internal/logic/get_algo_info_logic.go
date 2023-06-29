package logic

import (
	"algo_assess/pkg/tools"
	"context"

	"algo_assess/mornano-rpc-server/internal/svc"
	"algo_assess/mornano-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoInfoLogic {
	return &GetAlgoInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAlgoInfo 查询算法基础信息
func (l *GetAlgoInfoLogic) GetAlgoInfo(in *proto.AlgoInfoReq) (*proto.AlgoInfoRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in GetAlgoInfo: get req:", in)
	ret, err := l.svcCtx.AlgoInfoRepo.GetBusAlgoBase(l.ctx)
	if err != nil {
		l.Logger.Error("GetBusAlgoBase error:", err)
		return &proto.AlgoInfoRsp{
			Code:  350,
			Msg:   err.Error(),
			Infos: nil,
		}, nil
	}
	var list []*proto.AlgoInfo
	for _, v := range ret {
		a := &proto.AlgoInfo{
			AlgoId:       int32(v.Id),
			AlgoName:     tools.RMu0000(v.AlgoName),
			AlgoType:     int32(v.AlgorithmType),
			AlgoTypeName: GetAlgoTypeName(v.AlgorithmType),
			Provider:     tools.RMu0000(v.ProviderName),
		}
		list = append(list, a)
	}

	return &proto.AlgoInfoRsp{
		Code:  200,
		Msg:   "success",
		Infos: list,
	}, nil
}

func GetAlgoTypeName(u uint) string {
	if u == 1 {
		return "日内回转"
	} else if u == 2 {
		return "智能委托"
	}
	return "未知算法类型"
}
