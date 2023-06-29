package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddOptimizeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddOptimizeBaseLogic {
	return &AddOptimizeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddOptimizeBase 新增一键优选基础数据
func (l *AddOptimizeBaseLogic) AddOptimizeBase(in *proto.AddOptimizeBaseReq) (*proto.AddOptimizeBaseRsp, error) {
	if in.GetAlgoType() != 1 {
		return &proto.AddOptimizeBaseRsp{
			Code: 10000,
			Msg:  "算法类型错误，请检查算法类型后重试",
		}, nil
	}
	count, _ := l.svcCtx.OptimizeBaseRepo.CountOptimize(in.GetSecId(), int(in.GetAlgoId()))
	if count != 0 {
		return &proto.AddOptimizeBaseRsp{
			Code: 10000,
			Msg:  "证券id或算法id重复，请检查数据后重试",
		}, nil
	}
	err := l.svcCtx.OptimizeBaseRepo.AddAlgoOptimizeBase(in)
	if err != nil {
		return &proto.AddOptimizeBaseRsp{
			Code: 10000,
			Msg:  "新增失败",
		}, err
	}
	//secId := in.GetSecId()
	//info, _ := l.svcCtx.OptimizeRepo.GetSingleAlgoOptimize(secId)
	//score := util.CalOptimizeScore(in.GetOpenRate(), in.GetIncomeRate(), in.GetBasisPoint())
	//if info != nil {
	//	if score > info.Score {
	//		algoOptimize := &models.TbAlgoOptimize{
	//			Id:           info.Id,
	//			ProviderId:   int(in.GetProviderId()),
	//			ProviderName: in.GetProviderName(),
	//			SecName:      in.GetSecName(),
	//			AlgoId:       int(in.GetAlgoId()),
	//			AlgoType:     int(in.GetAlgoType()),
	//			AlgoName:     in.GetAlgoName(),
	//			Score:        score,
	//			UpateTime:    time.Now(),
	//		}
	//		fmt.Printf("algoOptimize%v\n", algoOptimize)
	//		algoOptimizeErr := l.svcCtx.OptimizeRepo.UpdateAlgoOptimize(algoOptimize)
	//		if algoOptimizeErr != nil {
	//			return &proto.AddOptimizeBaseRsp{
	//				Code: 10000,
	//				Msg:  "优选数据修改失败，请检查数据后重试",
	//			}, algoOptimizeErr
	//		}
	//	}
	//} else {
	//	algoOptimize := &models.TbAlgoOptimize{
	//		SecId:        secId,
	//		ProviderId:   int(in.GetProviderId()),
	//		ProviderName: in.GetProviderName(),
	//		SecName:      in.GetSecName(),
	//		AlgoId:       int(in.GetAlgoId()),
	//		AlgoType:     int(in.GetAlgoType()),
	//		AlgoName:     in.GetAlgoName(),
	//		Score:        score,
	//		CreateTime:   time.Now(),
	//		UpateTime:    time.Now(),
	//	}
	//	fmt.Printf("algoOptimize%v\n", algoOptimize)
	//	algoOptimizeErr := l.svcCtx.OptimizeRepo.AddAlgoOptimize(algoOptimize)
	//	if algoOptimizeErr != nil {
	//		return &proto.AddOptimizeBaseRsp{
	//			Code: 10000,
	//			Msg:  "优选数据新增失败，请检查数据后重试",
	//		}, algoOptimizeErr
	//	}
	//}
	return &proto.AddOptimizeBaseRsp{
		Code: 0,
		Msg:  "新增成功",
	}, nil
}
