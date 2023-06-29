package logic

import (
	"algo_assess/global"
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDefaultAlgoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDefaultAlgoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDefaultAlgoLogic {
	return &GetDefaultAlgoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetDefaultAlgo 取默认数据，用于需选择算法的场景下默认无数据的页面
func (l *GetDefaultAlgoLogic) GetDefaultAlgo(in *proto.DefaultReq) (*proto.DefaultRsp, error) {
	l.Logger.Infof("in GetDefaultAlgo ,get req:%+v", in)
	var algoId int32
	var algoTypeName, algoName, provider string
	if in.GetScene() == 1 { // 算法动态，只需返加一个algoId
		out, err := l.svcCtx.SummaryRepo.GetDefaultAlgoId(l.ctx, in.GetStartTime(), in.GetUserId(), in.GetUserType())
		if err != nil {
			l.Logger.Error("GetDefaultAlgoId error:", err)
			return &proto.DefaultRsp{
				Code:     380,
				Msg:      err.Error(),
				AlgoId:   0,
				AlgoName: "",
			}, nil
		}
		algoId = int32(out.AlgoId)
		if out.AlgoType == 1 {
			algoTypeName = global.AlgoTypeNameT0
		} else if out.AlgoType == 2 {
			algoTypeName = global.AlgoTypeNameSplit
		}
		algoName = out.AlgoName
		provider = out.Provider
	} else if in.GetScene() == 4 { // 胜率、信号
		out, err := l.svcCtx.SummaryRepo.GetAdvanceDefaultAlgo(l.ctx, in.GetStartTime(), in.GetEndTime(), in.GetUserId(), in.GetUserType())
		if err != nil {
			l.Logger.Error("GetDefaultAlgoId error:", err)
			return &proto.DefaultRsp{
				Code:     380,
				Msg:      err.Error(),
				AlgoId:   0,
				AlgoName: "",
			}, nil
		}
		algoId = int32(out.AlgoId)
		if out.AlgoType == 1 {
			algoTypeName = global.AlgoTypeNameT0
		} else if out.AlgoType == 2 {
			algoTypeName = global.AlgoTypeNameSplit
		}
		algoName = out.AlgoName
		provider = out.Provider
	} else if in.GetScene() == 5 { // 绩效分析
		out, err := l.svcCtx.SummaryOrigRepo.GetOrigDefaultAlgoId(l.ctx, in.GetStartTime(), in.GetUserId(), in.GetUserType(), in.GetBatchNo())
		if err != nil {
			l.Logger.Error("GetOrigDefaultAlgoId error:", err)
			return &proto.DefaultRsp{
				Code:     380,
				Msg:      err.Error(),
				AlgoId:   0,
				AlgoName: "",
			}, nil
		}
		algoId = int32(out.AlgoId)
		if out.AlgoType == 1 {
			algoTypeName = global.AlgoTypeNameT0
		} else if out.AlgoType == 2 {
			algoTypeName = global.AlgoTypeNameSplit
		}
		algoName = out.AlgoName
		provider = out.Provider
	}

	return &proto.DefaultRsp{
		Code:         200,
		Msg:          "success",
		AlgoId:       algoId,
		AlgoName:     algoName,
		AlgoTypeName: algoTypeName,
		Provider:     provider,
	}, nil
}
