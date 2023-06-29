package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"math"
	"math/rand"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WinRatioAnalyseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWinRatioAnalyseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WinRatioAnalyseLogic {
	return &WinRatioAnalyseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WinRatioAnalyseLogic) WinRatioAnalyse(req *types.WinRatioReq) (resp *types.WinRatioRsp, err error) {
	l.Logger.Infof("in WinRatioAnalyse, get req:%+v", *req)
	if l.svcCtx.Config.WorkControl.EnableFakeMsg {
		return BuildVirtualRsp(req)
	}

	var algoId int32
	var algoTypeName, algoName, provider string
	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))

	if req.AlgoName == "" { // 算法名称为空的， 取兜底数据
		l.Logger.Info("algo_name empty, get default algo....")
		dReq := &assessservice.DefaultReq{
			Scene:     4,
			UserId:    req.UserId,
			UserType:  int32(req.UserType),
			StartTime: start,
			EndTime:   end,
		}
		dRsp, err := l.svcCtx.AssessClient.GetDefaultAlgo(l.ctx, dReq)
		if err != nil {
			l.Logger.Error("rpc call GetDefaultAlgo error:", err)
			return &types.WinRatioRsp{
				Code:        350,
				Msg:         err.Error(),
				Head:        types.RatioHeader{},
				WinRatio:    nil,
				OddsRatio:   nil,
				ProfitRatio: nil,
			}, nil
		}
		algoId = dRsp.GetAlgoId()
		algoTypeName = dRsp.GetAlgoTypeName()
		algoName = dRsp.GetAlgoName()
		provider = dRsp.GetProvider()
	} else { // 根据算法名称反查algoId
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
	if algoId <= 0 { // 无数据时直接返回
		l.Logger.Info("no data....")
		return &types.WinRatioRsp{
			Code:        200,
			Msg:         "no data",
			Head:        types.RatioHeader{},
			WinRatio:    nil,
			OddsRatio:   nil,
			ProfitRatio: nil,
		}, nil
	}

	wReq := &proto.WinRatioReq{
		AlgoId:   algoId,
		StartDay: start,
		EndDay:   end,
		UserId:   req.UserId,
		UserType: int32(req.UserType),
	}
	wRsp, err := l.svcCtx.AssessClient.GetWinRatio(l.ctx, wReq)
	if err != nil {
		l.Logger.Error("rpc call GetWinRatio error:", err)
		return &types.WinRatioRsp{
			Code:        368,
			Msg:         err.Error(),
			Head:        types.RatioHeader{},
			WinRatio:    nil,
			OddsRatio:   nil,
			ProfitRatio: nil,
		}, nil
	}
	// 拼返回报文
	var wr []types.WR
	var odds []types.Odds
	var profit []types.Profit
	for _, v := range wRsp.GetInfo() {
		w := types.WR{
			Day:      v.GetDay(),
			WinRatio: v.GetWinRatio(),
		}
		o := types.Odds{
			Day:       v.GetDay(),
			OddsRatio: v.GetOdds(),
		}
		p := types.Profit{
			Day:         v.GetDay(),
			ProfitRatio: v.GetProfit(),
		}
		wr = append(wr, w)
		odds = append(odds, o)
		profit = append(profit, p)
	}
	rsp := &types.WinRatioRsp{
		Code:         200,
		Msg:          "success",
		Provider:     provider,
		AlgoTypeName: algoTypeName,
		AlgoName:     algoName,
		Head: types.RatioHeader{
			AlgoName:         wRsp.GetWinHead().GetAlgoName(),
			StartTime:        wRsp.GetWinHead().GetStartDay(),
			EndTime:          wRsp.GetWinHead().GetEndDay(),
			TradeDays:        int(wRsp.GetWinHead().GetTradeDays()),
			AvgDailyProfit:   wRsp.GetWinHead().GetAvgDailyProfit(),
			AnnualizedProfit: wRsp.GetWinHead().GetAnnualizedProfit(),
			TotalProfit:      wRsp.GetWinHead().GetTotalProfit(),
			MaxWithDraw:      wRsp.GetWinHead().GetMaxWithdraw(),
			DailyProgress:    wRsp.GetWinHead().GetDailyProgress(),
			DailyStocks:      int(wRsp.GetWinHead().GetDailyStocks()),
			ProfitDays:       int(wRsp.GetWinHead().GetProfitDays()),
			ProfitDaysRate:   wRsp.GetWinHead().GetProfitDayRate(),
		},
		WinRatio:    wr,
		OddsRatio:   odds,
		ProfitRatio: profit,
	}

	return rsp, nil
}

func BuildVirtualRsp(req *types.WinRatioReq) (*types.WinRatioRsp, error) {
	// TODO: 返回固定数据
	var an string
	if req.AlgoName != "" {
		an = req.AlgoName
	} else {
		an = "智能委托"
	}
	head := types.RatioHeader{
		AlgoName:         an,
		StartTime:        time.Now().AddDate(0, 0, -14).Format(global.TimeFormatDaySp),
		EndTime:          time.Now().AddDate(0, 0, -1).Format(global.TimeFormatDaySp),
		TradeDays:        14,
		AvgDailyProfit:   2.5,
		AnnualizedProfit: 13.8,
		TotalProfit:      3.5,
		MaxWithDraw:      20,
		MaxWithDrawDays:  3,
		DailyProgress:    92,
		DailyStocks:      25,
		ProfitDays:       8,
		ProfitDaysRate:   60,
	}

	var ws []types.WR
	var ods []types.Odds
	var ps []types.Profit
	for i := -14; i <= -1; i++ {
		d := time.Now().AddDate(0, 0, i).Format(global.TimeFormatDaySp)
		// 拼胜率
		w := types.WR{
			Day:      d,
			WinRatio: RandFloat64(50.0, 70.0),
		}
		ws = append(ws, w)
		// 拼赔率
		o := types.Odds{
			Day:       d,
			OddsRatio: RandFloat64(0.1, 5.0),
		}
		ods = append(ods, o)
		// 拼盈亏
		p := types.Profit{
			Day:         d,
			ProfitRatio: RandFloat64(500.0, 10000.0),
		}
		ps = append(ps, p)
	}

	return &types.WinRatioRsp{
		Code:        200,
		Msg:         "success",
		Head:        head,
		WinRatio:    ws,
		OddsRatio:   ods,
		ProfitRatio: ps,
	}, nil
}

// RandFloat64 生成指定范围内的浮点数随机数，精确到1位
func RandFloat64(min, max float64) float64 {
	return math.Round((min+rand.Float64()*(max-min))*10) / 10
}
