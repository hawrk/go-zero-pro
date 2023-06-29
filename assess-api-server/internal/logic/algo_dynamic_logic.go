package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"context"
	"encoding/json"
	"github.com/spf13/cast"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlgoDynamicLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoDynamicLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoDynamicLogic {
	return &AlgoDynamicLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AlgoDynamic 算法动态页面
func (l *AlgoDynamicLogic) AlgoDynamic(req *types.DynamicReq) (resp *types.DynamicRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in AlgoDynamic, get req:%+v", *req)
	// 算法动态数据获取
	dynamicCh := make(chan *assessservice.DynamicRsp)
	tlCh := make(chan *assessservice.TimeLineRsp)

	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))
	crossDayFlag := false
	if start != end {
		crossDayFlag = true
	}

	var algoId int32
	var algoTypeName, algoName, provider string
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
		algoName = req.AlgoName
	} else {
		if req.SourceFrom == global.SourceFromBus || req.SourceFrom == global.SourceFromFix { // 总线数据根据算法ID返回数据
			l.Logger.Info("algo_name empty, get default algo....")
			dReq := &assessservice.DefaultReq{
				Scene:     1,
				UserId:    req.UserId,
				UserType:  int32(req.UserType),
				StartTime: start,
				EndTime:   end,
			}
			dRsp, err := l.svcCtx.AssessClient.GetDefaultAlgo(l.ctx, dReq)
			if err != nil {
				l.Logger.Error("rpc call GetDefaultAlgo error:", err)
				return &types.DynamicRsp{
					Code:           350,
					Msg:            err.Error(),
					CrossDay:       crossDayFlag,
					Dimension:      []types.DimensionInfo{},
					CompositeScore: 0,
					Ranking:        0,
					MarketRate:     []types.DNMarketRateInfo{},
					Side:           types.DNTradeSide{},
					PriceType:      []types.StockPriceType{},
					VolType:        []types.TradeVol{},
					AssessLine:     types.DemensionLine{},
					ProgressLine:   types.DemensionLine{},
				}, nil
			}
			algoId = dRsp.GetAlgoId()
			algoTypeName = dRsp.GetAlgoTypeName()
			algoName = dRsp.GetAlgoName()
			provider = dRsp.GetProvider()
		} else { // 订单导入或原始订单有多个算法ID计算时，需要根据批次号默认返回一个算法
			l.Logger.Info("into multi algo, get default algo....")
			dReq := &assessservice.DefaultReq{
				Scene:     5,
				UserId:    req.UserId,
				UserType:  int32(req.UserType),
				StartTime: start,
				EndTime:   end,
				BatchNo:   req.BatchNo,
			}
			dRsp, _ := l.svcCtx.AssessClient.GetDefaultAlgo(l.ctx, dReq)
			algoId = dRsp.GetAlgoId()
			algoTypeName = dRsp.GetAlgoTypeName()
			algoName = dRsp.GetAlgoName()
			provider = dRsp.GetProvider()
		}
	}

	// 获取算法动态（五个维度，综合评分，评分描述，资金占比， 买卖方向，股价类型，交易量）
	go func() {
		dReq := &assessservice.DynamicReq{
			AlgoId:       algoId,
			UserId:       req.UserId,
			UserType:     int32(req.UserType),
			StartTime:    start,
			EndTime:      end,
			CrossDayFlag: crossDayFlag,
			SourceFrom:   req.SourceFrom,
			BatchNo:      req.BatchNo,
		}
		dRsp, err := l.svcCtx.AssessClient.GetAlgoDynamic(l.ctx, dReq)
		if err != nil {
			l.Logger.Error(" call assess rpc GetAlgoDynamic error:", err)
			dynamicCh <- &assessservice.DynamicRsp{}
			return
		}
		dynamicCh <- dRsp
	}()
	// 获取时间线图:绩效 ,完成度
	go func() {
		tlReq := &assessservice.TimeLineReq{
			LineType:     12,
			StartTime:    start,
			EndTime:      end,
			UserId:       req.UserId,
			UserType:     int32(req.UserType),
			AlgoId:       algoId,
			CrossDayFlag: crossDayFlag,
			SourceFrom:   req.SourceFrom,
			BatchNo:      req.BatchNo,
		}
		tlRsp, err := l.svcCtx.AssessClient.GetAlgoTimeLine(l.ctx, tlReq)
		if err != nil {
			l.Logger.Error("call assess rcp GetAlgoTimeLine error :", err)
			tlCh <- &assessservice.TimeLineRsp{}
			return
		}
		tlCh <- tlRsp
	}()

	dyRsp, tRsp := <-dynamicCh, <-tlCh
	resp = l.BuildDynamicRsp(dyRsp, tRsp, crossDayFlag, provider, algoTypeName, algoName)
	return resp, nil
}

func (l *AlgoDynamicLogic) BuildDynamicRsp(dRsp *assessservice.DynamicRsp, tRsp *assessservice.TimeLineRsp, cdflag bool,
	provider, algoTypeName, algoName string) *types.DynamicRsp {
	var di []types.DimensionInfo
	for _, v := range dRsp.Dimension {
		d := types.DimensionInfo{
			ProfileType: v.GetProfileType(),
			Score:       v.GetProfileScore(),
			Desc:        v.GetProfileDesc(),
		}
		di = append(di, d)
	}
	aTL, pTL := l.ParseTimeLine(tRsp.Line)

	out := &types.DynamicRsp{
		Code:           200,
		Msg:            "success",
		CrossDay:       cdflag,
		Provider:       provider,
		AlgoTypeName:   algoTypeName,
		AlgoName:       algoName,
		Dimension:      di,
		CompositeScore: dRsp.TotalScore,
		Ranking:        dRsp.Ranking,
		MarketRate:     l.ParseMarketRate(dRsp.FundRate),
		Side:           l.ParseTradeSide(dRsp.TradeSide),
		PriceType:      l.ParseStockType(dRsp.StockType),
		VolType:        l.ParseTradeVol(dRsp.TradeVol),
		AssessLine:     aTL,
		ProgressLine:   pTL,
		AlgoNameList:   dRsp.AlgoNames,
	}
	return out
}

// ParseMarketRate 计算资金占比
func (l *AlgoDynamicLogic) ParseMarketRate(js []string) []types.DNMarketRateInfo {
	if js == nil {
		return []types.DNMarketRateInfo{}
	}
	var totalSum, hugeSum, bigSum, middleSum, smallSum int64
	var sli []types.DNMarketRateInfo
	for _, v := range js {
		var fr global.FundRate
		if err := json.Unmarshal([]byte(v), &fr); err != nil {
			l.Logger.Error("Unmarshal FundRate error:", err, ", input:", v)
			return nil
		}
		hugeSum += fr.Huge
		bigSum += fr.Big
		middleSum += fr.Middle
		smallSum += fr.Small
	}
	totalSum = hugeSum + bigSum + middleSum + smallSum
	if totalSum <= 0 {
		return []types.DNMarketRateInfo{}
	}
	fSum := float64(totalSum)
	mk1 := types.DNMarketRateInfo{
		MkName: "超大市值",
		Rate:   float64(hugeSum) / fSum * 100,
	}
	mk2 := types.DNMarketRateInfo{
		MkName: "大市值",
		Rate:   float64(bigSum) / fSum * 100,
	}
	mk3 := types.DNMarketRateInfo{
		MkName: "中等市值",
		Rate:   float64(middleSum) / fSum * 100,
	}
	mk4 := types.DNMarketRateInfo{
		MkName: "小市值",
		Rate:   float64(smallSum) / fSum * 100,
	}
	sli = append(sli, mk1, mk2, mk3, mk4)

	return sli
}

// ParseTradeSide 计算买卖方向占比
func (l *AlgoDynamicLogic) ParseTradeSide(js []string) types.DNTradeSide {
	if js == nil {
		return types.DNTradeSide{}
	}
	var totalSum, buySum, sellSum int64
	for _, v := range js {
		var ts global.TradeVolDirect
		if err := json.Unmarshal([]byte(v), &ts); err != nil {
			l.Logger.Error("Unmarshal TradeSide error:", err, ", input:", v)
			return types.DNTradeSide{}
		}
		buySum += ts.BuyVol
		sellSum += ts.SellVol
	}
	totalSum = buySum + sellSum
	if totalSum <= 0 {
		return types.DNTradeSide{}
	}
	fSum := float64(totalSum)
	out := types.DNTradeSide{
		BuyRate:  float64(buySum) / fSum * 100,
		SellRate: float64(sellSum) / fSum * 100,
	}
	return out
}

// ParseStockType 计算股价类型占比
func (l *AlgoDynamicLogic) ParseStockType(js []string) []types.StockPriceType {
	if js == nil {
		return []types.StockPriceType{}
	}
	var totalSum, redSum, orangeSum, yellowSum, greenSum, cyanSum, blueSum, purpleSum int64
	var sli []types.StockPriceType
	for _, v := range js {
		var st global.StockType
		if err := json.Unmarshal([]byte(v), &st); err != nil {
			l.Logger.Error("Unmarshal stockType error:", err, ", input:", v)
			return nil
		}
		redSum += st.Red
		orangeSum += st.Orange
		yellowSum += st.Yellow
		greenSum += st.Green
		cyanSum += st.Cyan
		blueSum += st.Blue
		purpleSum += st.Purple
	}
	totalSum = redSum + orangeSum + yellowSum + greenSum + cyanSum + blueSum + purpleSum
	if totalSum <= 0 {
		return []types.StockPriceType{}
	}
	fSum := float64(totalSum)
	co1 := types.StockPriceType{
		TypeName: "普通股",
		Rate:     float64(redSum) / fSum * 100,
	}
	co2 := types.StockPriceType{
		TypeName: "蓝筹股",
		Rate:     float64(orangeSum) / fSum * 100,
	}
	co3 := types.StockPriceType{
		TypeName: "龙头股",
		Rate:     float64(yellowSum) / fSum * 100,
	}
	co4 := types.StockPriceType{
		TypeName: "黑马股",
		Rate:     float64(greenSum) / fSum * 100,
	}
	co5 := types.StockPriceType{
		TypeName: "ST股",
		Rate:     float64(cyanSum) / fSum * 100,
	}
	co6 := types.StockPriceType{
		TypeName: "*ST股",
		Rate:     float64(blueSum) / fSum * 100,
	}
	co7 := types.StockPriceType{
		TypeName: "PT股",
		Rate:     float64(purpleSum) / fSum * 100,
	}
	sli = append(sli, co1, co2, co3, co4, co5, co6, co7)
	return sli
}

// ParseTradeVol 计算交易量占比
func (l *AlgoDynamicLogic) ParseTradeVol(js []string) []types.TradeVol {
	if js == nil {
		return []types.TradeVol{}
	}
	var totalSum, billionSum, millionSum, thousandSum int64
	var sli []types.TradeVol
	for _, v := range js {
		var tv global.TradeVolRate
		if err := json.Unmarshal([]byte(v), &tv); err != nil {
			l.Logger.Error("Unmarshal TradeVol error:", err, ", input:", v)
			return nil
		}
		billionSum += tv.Billion
		millionSum += tv.Million
		thousandSum += tv.Thousand
	}
	totalSum = billionSum + millionSum + thousandSum
	if totalSum <= 0 {
		return []types.TradeVol{}
	}
	fSum := float64(totalSum)
	tv1 := types.TradeVol{
		VolName: "百万以上",
		Rate:    float64(billionSum) / fSum * 100,
	}
	tv2 := types.TradeVol{
		VolName: "百万以下",
		Rate:    float64(millionSum) / fSum * 100,
	}
	tv3 := types.TradeVol{
		VolName: "万元以下",
		Rate:    float64(thousandSum) / fSum * 100,
	}
	sli = append(sli, tv1, tv2, tv3)
	return sli
}

func (l *AlgoDynamicLogic) ParseTimeLine(line []*assessservice.TimeLine) (types.DemensionLine, types.DemensionLine) {
	var assTL []types.TimeLine
	var pgTL []types.TimeLine
	for _, v := range line {
		atl := types.TimeLine{
			TimePoint: v.TimePoint,
			Score:     v.AssessScore,
		}
		ptl := types.TimeLine{
			TimePoint: v.TimePoint,
			Score:     int32(v.Progress),
		}
		assTL = append(assTL, atl)
		pgTL = append(pgTL, ptl)
	}
	assess := types.DemensionLine{
		ProfileType: 4,
		Point:       assTL,
	}

	progress := types.DemensionLine{
		ProfileType: 2,
		Point:       pgTL,
	}
	return assess, progress
}
