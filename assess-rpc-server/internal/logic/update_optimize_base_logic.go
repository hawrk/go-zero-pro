package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOptimizeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOptimizeBaseLogic {
	return &UpdateOptimizeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOptimizeBase 修改一键优选基础数据
func (l *UpdateOptimizeBaseLogic) UpdateOptimizeBase(in *proto.UpdateOptimizeBaseReq) (*proto.UpdateOptimizeBaseRsp, error) {
	if in.GetAlgoType() != 1 {
		return &proto.UpdateOptimizeBaseRsp{
			Code: 10000,
			Msg:  "算法类型错误，请检查算法类型后重试",
		}, nil
	}
	err := l.svcCtx.OptimizeBaseRepo.UpdateAlgoOptimizeBase(in)
	if err != nil {
		return &proto.UpdateOptimizeBaseRsp{
			Code: 10000,
			Msg:  "修改失败",
		}, err
	}
	//singleAlgoOptimize, _ := l.svcCtx.OptimizeRepo.GetSingleAlgoOptimize(in.GetSecId())
	//score := util.CalOptimizeScore(in.GetOpenRate(), in.GetIncomeRate(), in.GetBasisPoint())
	//if singleAlgoOptimize.AlgoId == int(in.GetAlgoId()) {
	//	if score >= singleAlgoOptimize.Score {
	//		algoOptimize := &models.TbAlgoOptimize{
	//			Id:           singleAlgoOptimize.Id,
	//			ProviderId:   int(in.GetProviderId()),
	//			ProviderName: in.GetProviderName(),
	//			SecName:      in.GetSecName(),
	//			AlgoName:     in.GetAlgoName(),
	//			Score:        score,
	//			UpateTime:    time.Now(),
	//		}
	//		algoOptimizeErr := l.svcCtx.OptimizeRepo.UpdateAlgoOptimize(algoOptimize)
	//		if algoOptimizeErr != nil {
	//			return &proto.UpdateOptimizeBaseRsp{
	//				Code: 10000,
	//				Msg:  "优选数据修改失败，请检查数据后重试",
	//			}, algoOptimizeErr
	//		}
	//	} else {
	//		infos, err := l.svcCtx.OptimizeBaseRepo.SelectOptimizeBaseBySecId(in.GetSecId())
	//		if err != nil {
	//			return &proto.UpdateOptimizeBaseRsp{
	//				Code: 10000,
	//				Msg:  "优选数据修改失败，请检查数据后重试",
	//			}, err
	//		}
	//		var algoOptimize *models.TbAlgoOptimize
	//		for _, info := range infos {
	//			fmt.Printf("info:%v\n", info)
	//			algoOptimizeTemp := &models.TbAlgoOptimize{
	//				Id:           singleAlgoOptimize.Id,
	//				ProviderId:   info.ProviderId,
	//				ProviderName: info.ProviderName,
	//				SecId:        info.SecId,
	//				SecName:      info.SecName,
	//				AlgoId:       info.AlgoId,
	//				AlgoType:     info.AlgoType,
	//				AlgoName:     info.AlgoName,
	//				Score:        util.CalOptimizeScore(info.OpenRate, info.IncomeRate, info.BasisPoint),
	//				UpateTime:    time.Now(),
	//			}
	//			if algoOptimize != nil {
	//				if algoOptimizeTemp.Score > algoOptimize.Score {
	//					algoOptimize = algoOptimizeTemp
	//				}
	//			} else {
	//				algoOptimize = algoOptimizeTemp
	//			}
	//		}
	//		algoOptimizeErr := l.svcCtx.OptimizeRepo.UpdateAlgoOptimize(algoOptimize)
	//		if algoOptimizeErr != nil {
	//			return &proto.UpdateOptimizeBaseRsp{
	//				Code: 10000,
	//				Msg:  "优选数据修改失败，请检查数据后重试",
	//			}, algoOptimizeErr
	//		}
	//	}
	//} else {
	//	if score > singleAlgoOptimize.Score {
	//		algoOptimize := &models.TbAlgoOptimize{
	//			Id:           singleAlgoOptimize.Id,
	//			ProviderId:   int(in.GetProviderId()),
	//			ProviderName: in.GetProviderName(),
	//			SecId:        in.GetSecId(),
	//			SecName:      in.GetSecName(),
	//			AlgoId:       int(in.GetAlgoId()),
	//			AlgoType:     int(in.GetAlgoType()),
	//			AlgoName:     in.GetAlgoName(),
	//			Score:        score,
	//			UpateTime:    time.Now(),
	//		}
	//		algoOptimizeErr := l.svcCtx.OptimizeRepo.UpdateAlgoOptimize(algoOptimize)
	//		if algoOptimizeErr != nil {
	//			return &proto.UpdateOptimizeBaseRsp{
	//				Code: 10000,
	//				Msg:  "优选数据修改失败，请检查数据后重试",
	//			}, algoOptimizeErr
	//		}
	//	}
	//}

	return &proto.UpdateOptimizeBaseRsp{
		Code: 0,
		Msg:  "修改成功",
	}, nil
}
