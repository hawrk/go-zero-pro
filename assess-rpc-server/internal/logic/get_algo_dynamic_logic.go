package logic

import (
	"algo_assess/assess-rpc-server/internal/config"
	"algo_assess/global"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"
	"math"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoDynamicLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoDynamicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoDynamicLogic {
	return &GetAlgoDynamicLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAlgoDynamic  算法动态
func (l *GetAlgoDynamicLogic) GetAlgoDynamic(in *proto.DynamicReq) (*proto.DynamicRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in GetAlgoDynamic, get req:", in)
	// 先根据用户ID到账户表反查用户名称和角色权限--超管不需要
	account, err := l.svcCtx.UserInfoRepo.GetAccountInfoByUserId(l.ctx, in.GetUserId())
	if err != nil {
		l.Logger.Error("get account info error:", err)
	}
	// 取综合评分列表
	role := account.UserType
	if role == 0 { // 默认能取到的 1-普通用户 2-算法厂商 3-管理员，如果都不是，那就给一个超管的角色
		role = 4
	}

	var rsp proto.DynamicRsp
	if in.GetCrossDayFlag() { // 跨天处理
		var out []*models.TbAlgoSummary
		var algoNameList []string
		if in.GetSourceFrom() == global.SourceFromOrigin || in.GetSourceFrom() == global.SourceFromImport {
			out, err = l.svcCtx.SummaryOrigRepo.GetCrossDaySummaryByAlgoIdBatch(l.ctx, in.GetAlgoId(), in.GetBatchNo(), in.GetUserId(), in.GetUserType())
			if err != nil {
				l.Logger.Error("in GetCrossDaySummaryByAlgoIdBatch error:", err)
				return &proto.DynamicRsp{}, nil
			}
		} else {
			out, err = l.svcCtx.SummaryRepo.GetCrossDaySummaryByAlgoId(l.ctx, in.GetAlgoId(), in.GetStartTime(), in.GetEndTime(), in.GetUserId(), in.GetUserType())
			if err != nil {
				l.Logger.Error("in GetCrossDaySummaryByAlgoIds, error:", err)
				return &proto.DynamicRsp{}, nil
			}
		}

		var ecoList, progressList, riskList, assessList, stableList, comsumList []int
		var fundList, tradedirectList, stocktypeList, tradevolList []string
		for _, v := range out {
			ecoList = append(ecoList, v.EconomyScore)
			progressList = append(progressList, v.ProgressScore)
			riskList = append(riskList, v.RiskScore)
			assessList = append(assessList, v.AssessScore)
			stableList = append(stableList, v.StableScore)
			comsumList = append(comsumList, v.CumsumScore)
			fundList = append(fundList, v.FundRateJson)
			tradedirectList = append(tradedirectList, v.TradeDirectJson)
			stocktypeList = append(stocktypeList, v.StockTypeJson)
			tradevolList = append(tradevolList, v.TradeVolJson)
		}
		ecoScore := GetIntAvg(ecoList)
		progressScore := GetIntAvg(progressList)
		riskScore := GetIntAvg(riskList)
		assessScore := GetIntAvg(assessList)
		stableScore := GetIntAvg(stableList)
		comsumScore := GetIntAvg(comsumList)

		// 查分析的算法列表
		algoNameList, err = l.svcCtx.SummaryOrigRepo.GetAlgoNameByBatchNo(l.ctx, 0, in.GetBatchNo())
		if err != nil {
			l.Logger.Error("GetAlgoNameByBatchNo error:", err)
			return &proto.DynamicRsp{}, nil
		}

		rsp = proto.DynamicRsp{
			Code:       200,
			Msg:        "success",
			Dimension:  BuildAlgoDimension(ecoScore, progressScore, riskScore, assessScore, stableScore),
			TotalScore: int32(comsumScore),
			Ranking:    1,
			FundRate:   fundList,
			TradeSide:  tradedirectList,
			StockType:  stocktypeList,
			TradeVol:   tradevolList,
			AlgoNames:  algoNameList,
		}

	} else { // 当天
		var result []*models.TbAlgoSummary
		var scoreList []int
		var algoNameList []string
		if in.SourceFrom == global.SourceFromImport || in.SourceFrom == global.SourceFromOrigin { // 订单导入按照批次号来查
			result, err = l.svcCtx.SummaryOrigRepo.GetImportAlgoSummary(l.ctx, in.GetUserId(), in.GetUserType(), in.GetStartTime(), in.GetBatchNo())
			if err != nil {
				l.Logger.Error("GetImportAlgoSummary error:", err)
				return &proto.DynamicRsp{}, nil
			}
			// 查分析的算法列表
			algoNameList, err = l.svcCtx.SummaryOrigRepo.GetAlgoNameByBatchNo(l.ctx, in.GetStartTime(), in.GetBatchNo())
			if err != nil {
				l.Logger.Error("GetAlgoNameByBatchNo error:", err)
				return &proto.DynamicRsp{}, nil
			}
			// 查排名
			scoreList, err = l.svcCtx.SummaryOrigRepo.GetCumsumList(l.ctx, in.GetStartTime(), role)
			if err != nil {
				l.Logger.Error("GetCumsumList error:", err)
				return &proto.DynamicRsp{
					Code: 205,
					Msg:  err.Error(),
				}, nil
			}

		} else { // 其他先按算法ID查询
			result, err = l.svcCtx.SummaryRepo.GetAlgoSummary(l.ctx, in.GetUserId(), in.GetUserType(), in.GetStartTime(), in.GetAlgoId())
			if err != nil {
				l.Logger.Error("get algo summary error:", err)
				return &proto.DynamicRsp{}, nil
			}

			// 查排名
			scoreList, err = l.svcCtx.SummaryRepo.GetCumsumList(l.ctx, in.GetStartTime(), role)
			if err != nil {
				l.Logger.Error("GetCumsumList error:", err)
				return &proto.DynamicRsp{
					Code: 205,
					Msg:  err.Error(),
				}, nil
			}
			for _, v := range result {
				algoNameList = append(algoNameList, v.AlgoName)
			}
		}
		if len(result) >= 2 { // 不允许有超过2条的数据
			l.Logger.Error("more than 1 record.....")
		}

		for _, v := range result {
			rsp = proto.DynamicRsp{
				Code:       200,
				Msg:        "success",
				Dimension:  BuildAlgoDimension(v.EconomyScore, v.ProgressScore, v.RiskScore, v.AssessScore, v.StableScore),
				TotalScore: int32(v.CumsumScore),
				Ranking:    tools.GetRanking(scoreList, v.CumsumScore),
				FundRate:   []string{v.FundRateJson},
				TradeSide:  []string{v.TradeDirectJson},
				StockType:  []string{v.StockTypeJson},
				TradeVol:   []string{v.TradeVolJson},
				AlgoNames:  algoNameList,
			}
		}
	}
	return &rsp, nil
}

// GetIntAvg 求整型列表里的平均数，返回四舍五入的平均值
func GetIntAvg(l []int) int {
	if len(l) == 0 {
		return 0
	}
	n := len(l)
	var sum int
	for _, v := range l {
		sum += v
	}
	return int(math.Round(float64(sum) / float64(n)))
}

func BuildAlgoDimension(ecoScore, progressScore, riskScore, assessScore, stableScore int) []*proto.AlgoDimension {
	var ds []*proto.AlgoDimension
	d1 := &proto.AlgoDimension{ //经济性
		ProfileType:  1,
		ProfileScore: int32(ecoScore),
		ProfileDesc:  config.GetEconomyDesc(ecoScore),
	}
	d2 := &proto.AlgoDimension{ // 完成度
		ProfileType:  2,
		ProfileScore: int32(progressScore),
		ProfileDesc:  config.GetProgressDesc(progressScore),
	}
	d3 := &proto.AlgoDimension{ // 风险度
		ProfileType:  3,
		ProfileScore: int32(riskScore),
		ProfileDesc:  config.GetRiskDesc(riskScore),
	}
	d4 := &proto.AlgoDimension{ // 绩效
		ProfileType:  4,
		ProfileScore: int32(assessScore),
		ProfileDesc:  config.GetAssessDesc(assessScore),
	}
	d5 := &proto.AlgoDimension{ // 稳定性
		ProfileType:  5,
		ProfileScore: int32(stableScore),
		ProfileDesc:  config.GetStabilityDesc(stableScore),
	}
	ds = append(ds, d1, d2, d3, d4, d5)
	return ds
}
