package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProfileLogic {
	return &ProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProfileLogic) Profile(req *types.ProfileReq) (resp *types.ProfileRsp, err error) {
	l.Logger.Infof("get req:%+v", req)
	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatMinInt))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatMinInt))
	var algoId int32
	if req.AlgoName != "" {
		// 先反查一下算法ID
		alReq := &assessservice.ChooseAlgoReq{
			ChooseType: 4,
			AlgoName:   req.AlgoName,
		}

		alRsp, err := l.svcCtx.AssessClient.ChooseAlgoInfo(l.ctx, alReq)
		if err != nil {
			l.Logger.Error("call rpc ChooseAlgoInfo error:", err)
			return nil, err
		}
		algoId = alRsp.GetAlgoId()
	}

	profileReq := &assessservice.AlgoProfileReq{
		Provider:     req.Provider,
		AlgoTypeName: req.AlgoTypeName,
		AlgoId:       algoId,
		AlgoName:     req.AlgoName,
		UserId:       req.UserId,
		StartTime:    start,
		EndTime:      end,
		Page:         req.Page,
		Limit:        req.Limit,
		ProfileType:  req.ProfileType,
		UserType:     int32(req.UserType),
		SourceFrom:   req.SourceFrom, // 区分请求源
		BatchNo:      req.BatchNo,
	}
	rsp, err := l.svcCtx.AssessClient.GetAlgoProfile(l.ctx, profileReq)
	if err != nil {
		l.Logger.Error("profile rpc call error:", err)
		return nil, nil
	}
	p := &types.ProfileRsp{
		Code:      200,
		Msg:       "success",
		Total:     rsp.Total,
		Economy:   GetEconomy(req.ProfileType, rsp.Info),
		Progress:  GetProgress(req.ProfileType, rsp.Info),
		Risk:      GetRisk(req.ProfileType, rsp.Info),
		Assess:    GetAssess(req.ProfileType, rsp.Info),
		Stability: GetStability(req.ProfileType, rsp.Info),
	}

	return p, nil
}

func BuildProfileHead(head *assessservice.ProfileInfo) types.ProfileHead {
	return types.ProfileHead{
		BatchNo:     head.GetBatchNo(),
		AccountId:   head.GetAccoutId(),
		AccountName: head.GetAccountName(),
		Provider:    head.GetProvider(),
		AlgoId:      head.GetAlgoId(),
		AlgoName:    head.GetAlgoName(),
		CreateTime:  head.GetCreateTime(),
		SecId:       head.GetSecId(),
		SecName:     head.GetSecName(),
		AlgoOrderId: head.GetAlgoOrderId(),
		Industry:    head.GetIndustry(),
		FundType:    head.GetFundType(),
		Flowability: head.GetLiquidity(),
	}
}

func GetEconomy(t int32, data []*assessservice.ProfileInfo) (ret []types.EconomyStruct) {
	if t != 1 || len(data) <= 0 {
		return []types.EconomyStruct{}
	}

	for _, v := range data {
		economy := types.EconomyStruct{
			ProfileHead:    BuildProfileHead(v),
			TradeVol:       float64(v.GetTradeVol()) / 10000,
			Profit:         float64(v.GetProfit()) / 10000,
			ProfitRate:     v.GetProfitRate(),
			TotalFee:       float64(v.GetTotalFee()) / 10000,
			CrossFee:       float64(v.GetCrossFee()) / 10000,
			CancelRate:     v.GetCancelRate(),
			MinSplitOrder:  int32(v.GetMinSplitOrder()),
			DealEffi:       v.GetDealEffi(),
			AvgTradePrice:  v.GetAvgTradePrice() / 10000,  // 执行均价
			AvgArrivePrice: v.GetAvgArrivePrice() / 10000, // 到达均价
			PWP:            v.GetPwp() / 10000,
			AlgoDuration:   v.GetAlgoDuration(), // 有效时长
			TWAP:           v.GetTwap(),
			TWAPDev:        v.GetTwapDev(),
			VWAP:           v.GetVwap(),
			VWAPDev:        v.GetVwapDev(),
		}
		ret = append(ret, economy)
	}
	return
}

func GetProgress(t int32, data []*assessservice.ProfileInfo) (ret []types.ProgressStruct) {
	if t != 2 || len(data) <= 0 {
		return []types.ProgressStruct{}
	}

	for _, v := range data {
		progress := types.ProgressStruct{
			ProfileHead:  BuildProfileHead(v),
			Progress:     v.GetProgress(),
			AlgoOrderFit: v.GetAlgoOrderFit(),
			PriceFit:     v.GetPriceFit(),
			TradeVolFit:  v.GetTradeVolFit(),
		}
		ret = append(ret, progress)
	}
	return
}

func GetRisk(t int32, data []*assessservice.ProfileInfo) (ret []types.RiskStruct) {
	if t != 3 || len(data) <= 0 {
		return []types.RiskStruct{}
	}
	for _, v := range data {
		risk := types.RiskStruct{
			ProfileHead:  BuildProfileHead(v),
			MinJointRate: v.GetMinJointRate(),
			ProfitRate:   v.GetProfitRate(),
			WithdrawRate: v.GetWithdrawRate(),
		}
		ret = append(ret, risk)
	}
	return
}

func GetAssess(t int32, data []*assessservice.ProfileInfo) (ret []types.AssessStruct) {
	if t != 4 || len(data) <= 0 {
		return []types.AssessStruct{}
	}

	for _, v := range data {
		assess := types.AssessStruct{
			ProfileHead: BuildProfileHead(v),
			VwapDev:     v.GetVwapDev(),
			ProfitRate:  v.GetAssessProfitRate(),
		}
		ret = append(ret, assess)
	}
	return
}

func GetStability(t int32, data []*assessservice.ProfileInfo) (ret []types.StabilityStruct) {
	if t != 5 || len(data) <= 0 {
		return []types.StabilityStruct{}
	}
	for _, v := range data {
		stability := types.StabilityStruct{
			ProfileHead:       BuildProfileHead(v),
			VwapStdDev:        v.GetVwapStdDev(),
			ProfitRateStd:     v.GetPfRateStdDev(), //收益率标准差
			JointRate:         v.GetMinJointRate(),
			TradeVolFitStdDev: v.GetTradeVolFitStdDev(),
			TimeFitStdDev:     v.GetTimeFitStdDev(),
		}
		ret = append(ret, stability)
	}
	return
}
