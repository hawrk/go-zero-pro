package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOptimizeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOptimizeBaseLogic {
	return &DeleteOptimizeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除一键优选基础数据
func (l *DeleteOptimizeBaseLogic) DeleteOptimizeBase(in *proto.DeleteOptimizeBaseReq) (*proto.DeleteOptimizeBaseRsp, error) {
	//optimizeBase, _ := l.svcCtx.OptimizeBaseRepo.SelectOptimizeBaseById(in.GetId())
	err := l.svcCtx.OptimizeBaseRepo.DeleteAlgoOptimizeBase(in.GetId())
	if err != nil {
		return &proto.DeleteOptimizeBaseRsp{
			Code: 10000,
			Msg:  "删除失败",
		}, err
	}
	//optimize, _ := l.svcCtx.OptimizeRepo.GetAlgoOptimize(optimizeBase.SecId, optimizeBase.AlgoId)
	//if optimize != nil {
	//	infos, _ := l.svcCtx.OptimizeBaseRepo.SelectOptimizeBaseBySecId(optimizeBase.SecId)
	//	var algoOptimize *models.TbAlgoOptimize
	//	if len(infos) == 0 {
	//		algoOptimize = &models.TbAlgoOptimize{
	//			Id:        optimize.Id,
	//			AlgoId:    0,
	//			AlgoType:  0,
	//			AlgoName:  "",
	//			Score:     0.0,
	//			UpateTime: time.Now(),
	//		}
	//		fmt.Printf("algoOptimize%v\n", algoOptimize)
	//		_ = l.svcCtx.OptimizeRepo.InitAlgoOptimize(algoOptimize)
	//	} else {
	//		for _, info := range infos {
	//			algoOptimizeTemp := &models.TbAlgoOptimize{
	//				Id:        optimize.Id,
	//				SecId:     info.SecId,
	//				SecName:   info.SecName,
	//				AlgoId:    info.AlgoId,
	//				AlgoType:  info.AlgoType,
	//				AlgoName:  info.AlgoName,
	//				Score:     util.CalOptimizeScore(info.OpenRate, info.IncomeRate, info.BasisPoint),
	//				UpateTime: time.Now(),
	//			}
	//			if algoOptimize != nil {
	//				if algoOptimizeTemp.Score > algoOptimize.Score {
	//					algoOptimize = algoOptimizeTemp
	//				}
	//			} else {
	//				algoOptimize = algoOptimizeTemp
	//			}
	//		}
	//		_ = l.svcCtx.OptimizeRepo.UpdateAlgoOptimize(algoOptimize)
	//	}
	//}
	return &proto.DeleteOptimizeBaseRsp{
		Code: 0,
		Msg:  "删除成功",
	}, nil
}
