package logic

import (
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"encoding/json"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOrderSummaryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserOrderSummaryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOrderSummaryLogic {
	return &GetUserOrderSummaryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  取用户订单汇总信息(dashboard 用户数量， 交易量，买卖占比，厂商个数，资金占比,完成度--查汇总表）
func (l *GetUserOrderSummaryLogic) GetUserOrderSummary(in *proto.OrderSummaryReq) (*proto.OrderSummaryRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("get req:", in)
	ac, pc, tac, tpc := l.GetAlgoSummary(in)
	// 查一个基础用户表数据，把用户总数量查回来
	totalUserCnt, err := l.svcCtx.UserInfoRepo.GetAccountInfoCounts(l.ctx)
	if err != nil {
		l.Logger.Error("get account count error:", err)
	}
	//取整体数据
	algoSummary, orderSummary, err := l.svcCtx.SummaryRepo.GetUserOrderSummary(l.ctx, in.GetDate(), in.GetUserId(), in.GetUserType()) //根据普通用户统计
	if err != nil {
		l.Logger.Error("query algo_summary error:", err)
		return &proto.OrderSummaryRsp{}, nil
	}
	var rsp proto.OrderSummaryRsp
	var buySum, sellSum, sideSum int64
	var hugeSum, bigSum, middleSum, smallSum, fundSum int64

	// 买买占比，资金占比需要算出来
	for _, v := range algoSummary {
		fundJson := v.FundRateJson
		sideJson := v.TradeDirectJson
		var fund global.FundRate
		var side global.TradeVolDirect

		if err := json.Unmarshal(tools.String2Bytes(fundJson), &fund); err != nil {
			l.Logger.Error("unmarshal fundJson error:", err)
		}
		if err := json.Unmarshal(tools.String2Bytes(sideJson), &side); err != nil {
			l.Logger.Error("unmarshal sideJson error:", err)
		}
		buySum += side.BuyVol
		sellSum += side.SellVol

		hugeSum += fund.Huge
		bigSum += fund.Big
		middleSum += fund.Middle
		smallSum += fund.Small
	}
	fundSum += hugeSum + bigSum + middleSum + smallSum
	sideSum += buySum + sellSum
	l.Logger.Info("get fundSum:", fundSum, ", get sidesum:", sideSum)

	// 汇总数据直接取
	for _, v := range orderSummary {
		rsp.UserNum = v.UserNum
		rsp.TotalUserNum = totalUserCnt
		rsp.TotalTradeAmount = v.TradeAmount
		rsp.OrderNum = v.OrderNum
		rsp.ProviderNum = v.ProviderNum
		if v.EntrustQty > 0 {
			rsp.Progress = float64(v.DealQty) / float64(v.EntrustQty)
		}
	}
	if fundSum > 0 {
		fr := &proto.FundRate{
			Huge:   float64(hugeSum) / float64(fundSum),
			Big:    float64(bigSum) / float64(fundSum),
			Middle: float64(middleSum) / float64(fundSum),
			Small:  float64(smallSum) / float64(fundSum),
		}
		rsp.FundRate = fr
	}
	if sideSum > 0 {
		rsp.BuyRate = float64(buySum) / float64(sideSum)
		rsp.SellRate = float64(sellSum) / float64(sideSum)
	}
	// 填充算法信息
	rsp.AlgoCount = ac
	rsp.ProviderCount = pc
	rsp.TradeAlgoCount = tac
	rsp.TradeProviderCount = tpc

	return &rsp, nil
}

// GetAlgoSummary 统计Dashboard 算法数量和厂商数据信息，为原接口AlgoSummary迁移过来，旧接口作废
func (l *GetUserOrderSummaryLogic) GetAlgoSummary(in *proto.OrderSummaryReq) (int32, int32, int32, int32) {
	// 静态数据
	algoCnt, providerCnt, err := l.svcCtx.AlgoInfoRepo.GetAlgoSummary(l.ctx)
	if err != nil {
		l.Logger.Error("get algo summary data error:", err)
		return 0, 0, 0, 0
	}
	// 如果是超级管理员，则直接查summary表
	if in.GetUserType() == global.UserTypeAdmin {
		tradeAlgoCnt, tradeProviderCnt, err := l.svcCtx.SummaryRepo.GetAlgoSummaryCntByAdmin(l.ctx, in.GetDate())
		if err != nil {
			l.Logger.Error("GetAlgoSummaryCntByAdmin error:", err)
			return 0, 0, 0, 0
		}
		return algoCnt, providerCnt, tradeAlgoCnt, tradeProviderCnt
	} else {
		// 先根据用户ID到账户表反查用户名称和角色权限
		account, err := l.svcCtx.UserInfoRepo.GetAccountInfoByUserId(l.ctx, in.GetUserId())
		if err != nil {
			l.Logger.Error("get account info error:", err)
			return 0, 0, 0, 0
		}
		// 当天有交易的数据
		tradeAlgoCnt, tradeProviderCnt, err := l.svcCtx.SummaryRepo.GetAlgoSummaryCnt(l.ctx, in.GetDate(), in.GetUserId(), account.UserName, account.UserType)
		if err != nil {
			l.Logger.Error("get algo summary trade count error:", err)
			return 0, 0, 0, 0
		}
		return algoCnt, providerCnt, tradeAlgoCnt, tradeProviderCnt
	}
}
