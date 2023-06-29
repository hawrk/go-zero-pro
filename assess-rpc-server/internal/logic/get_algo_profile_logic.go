package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoProfileLogic {
	return &GetAlgoProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAlgoProfile 获取算法画像信息
func (l *GetAlgoProfileLogic) GetAlgoProfile(in *proto.AlgoProfileReq) (*proto.AlgoProfileRsp, error) {
	l.Logger.Infof("get req:%+v", in)
	var infos []*proto.ProfileInfo
	var total int64
	// 根据 算法类型名称反查一下算法类型ID
	var algoType int32
	if in.GetAlgoTypeName() != "" {
		r, err := l.svcCtx.AlgoInfoRepo.GetAlgoTypeIdByAlgoTypeName(l.ctx, in.GetAlgoTypeName())
		if err != nil {
			l.Logger.Error("error get GetAlgoTypeIdByAlgoTypeName:", err)
			return &proto.AlgoProfileRsp{}, nil
		}
		algoType = r
	}
	var count int64
	var result []*models.TbAlgoProfile
	var err error
	if in.SourceFrom == global.SourceFromImport || in.SourceFrom == global.SourceFromOrigin {
		result, count, err = l.svcCtx.ProfileOrigRepo.GetProfiles(l.ctx, in, algoType)
	} else {
		result, count, err = l.svcCtx.ProfileRepo.GetProfiles(l.ctx, in, algoType)
	}
	if err != nil {
		l.Logger.Error("error get GetProfileForEconomy:", err)
		return &proto.AlgoProfileRsp{}, nil
	}
	for _, v := range result {
		//if v.AccountType == 4 { // 不返回超管虚拟用户信息
		//	continue
		//}
		d := proto.ProfileInfo{
			BatchNo:     v.BatchNo,
			AccoutId:    v.AccountId,
			AccountName: v.AccountName,
			AlgoId:      int32(v.AlgoId),
			AlgoType:    int32(v.AlgoType),
			AlgoName:    v.AlgoName,
			TradeVol:    v.TradeCost,
			Profit:      v.ProfitAmount,
			ProfitRate:  tools.MulPercent(v.ProfitRate),
			TotalFee:    v.TotalTradeFee,
			CrossFee:    v.CrossFee,
			// add guangda
			AvgTradePrice:  v.AvgTradePrice,
			AvgArrivePrice: v.AvgArrivePrice,
			Pwp:            v.Pwp,
			AlgoDuration:   v.AlgoDuration,
			Twap:           v.Twap,
			TwapDev:        v.TwapDev,
			Vwap:           v.Vwap,
			VwapDev:        v.VwapDev,

			CancelRate:    tools.MulPercent(v.CancelRate),
			MinSplitOrder: int64(v.MiniSplitOrder),
			Progress:      v.ProgressRate,
			MinJointRate:  tools.MulPercent(v.MiniJointRate),
			WithdrawRate:  v.WithdrawRate,
			VwapStdDev:    v.StandardDeviation,
			PfRateStdDev:  v.PfRateStdDev, // 收益率标准差
			CreateTime:    tools.TimeMini2String(v.OrderTime),
			Provider:      v.Provider,
			SecId:         v.SecId,
			SecName:       v.SecName,
			AlgoOrderId:   v.AlgoOrderId,
			Industry:      v.Industry,
			FundType:      int32(v.FundType),
			Liquidity:     int32(v.Liquidity),
			DealEffi:      v.DealEffi,
			AlgoOrderFit:  v.AlgoOrderFit,
			//PriceFit:          v.PriceFit,
			TradeVolFit:       v.TradeVolFit,
			TradeVolFitStdDev: v.DealFitStdDev,
			TimeFitStdDev:     v.TimeFitStdDev,
			AssessProfitRate:  v.Factor,
		}
		infos = append(infos, &d)
	}
	total = count

	// ret
	l.Logger.Info("get ret total:", total, ",ret len:", len(infos))
	rsp := &proto.AlgoProfileRsp{
		Code:  0,
		Msg:   "success",
		Total: total,
		Info:  infos,
	}
	return rsp, nil
}
