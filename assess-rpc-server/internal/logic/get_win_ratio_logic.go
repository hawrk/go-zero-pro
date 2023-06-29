package logic

import (
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"github.com/spf13/cast"
	"strconv"
	"time"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWinRatioLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWinRatioLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWinRatioLogic {
	return &GetWinRatioLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetWinRatio 高阶分析： 胜率分析
func (l *GetWinRatioLogic) GetWinRatio(in *proto.WinRatioReq) (*proto.WinRatioRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in GetWinRatio, get req:", in)
	// 查profile表把胜率信息查回来
	out, err := l.svcCtx.ProfileRepo.GetWinRatioProfile(l.ctx, in.GetUserId(), in.GetUserType(), in.GetAlgoId(),
		in.GetStartDay(), in.GetEndDay())
	if err != nil {
		l.Logger.Error("GetWinRatioProfile error:", err)
		return &proto.WinRatioRsp{
			Code:    367,
			Msg:     err.Error(),
			WinHead: nil,
			Info:    nil,
		}, nil
	}
	stockCount, err := l.svcCtx.ProfileRepo.GetDailyStock(l.ctx, in.GetUserId(), in.GetUserType(), in.GetAlgoId(), in.GetStartDay(), in.GetEndDay())
	if err != nil {
		l.Logger.Error("GetDailyStock error:", err)
		return &proto.WinRatioRsp{
			Code:    367,
			Msg:     err.Error(),
			WinHead: nil,
			Info:    nil,
		}, nil
	}
	if len(out) == 0 { //无数据，直接返回
		l.Logger.Info("win ratio no data, return....")
		return &proto.WinRatioRsp{
			Code:    368,
			Msg:     "success",
			WinHead: nil,
			Info:    nil,
		}, nil
	}
	//dateSet := make(map[int]struct{})   // 统计交易天数
	n := len(out)                      // 交易天数
	var profitSum, progressSum float64 // 总收益率
	var profileSum float64             // 总盈亏金额
	var tradeCost int64                // 总交易金额
	var profitDays int32
	var maxWithdraw float64 // 最大回撤
	var algoName string

	var arr []*proto.WinRatioInfo
	strStartDay := cast.ToString(in.GetStartDay())
	y := cast.ToInt(strStartDay[:4])
	m := cast.ToInt(strStartDay[4:6])
	mk := time.Month(m)
	d := cast.ToInt(strStartDay[6:8])
	for i := 0; i <= 360; i++ {
		startDay := time.Date(y, mk, d, 0, 0, 0, 0, time.Local).AddDate(0, 0, i).Format(global.TimeFormatDay)
		iDay := cast.ToInt(startDay)
		for _, v := range out {
			//填充无数据的日期
			var a *proto.WinRatioInfo
			if cast.ToInt64(startDay) > in.GetEndDay() {
				break
			}
			if iDay < v.Date { // 当前日期在DB中无数据，填充0
				a = &proto.WinRatioInfo{
					Day:      TransDay(iDay), // 20220410 转成 2022.04.10 格式
					WinRatio: 0.00,
					Odds:     0.00,
					Profit:   0.00,
				}
				arr = append(arr, a)
				break
			} else if iDay == v.Date {
				//dateSet[v.Date] = struct{}{}
				profitSum += v.ProfitRate
				maxWithdraw = getMaxWithdraw(maxWithdraw, v.WithdrawRate)
				progressSum += v.ProgressRate
				if v.ProfitAmount > 0 {
					profitDays++
				}
				tradeCost += v.TradeCost
				profileSum += v.ProfitAmount
				algoName = v.AlgoName
				// 拼明细
				a = &proto.WinRatioInfo{
					Day:      TransDay(v.Date),                                        // 20220410 转成 2022.04.10 格式
					WinRatio: float64(v.TradeCountPlus) / float64(v.TradeCount) * 100, // 胜率
					Odds:     getOdds(v.TradeCountPlus, v.TradeCount, v.ProfitRate),   // 赔率 胜率 * 收益率/(1-胜率）
					Profit:   float64(v.ProfitAmount) / 10000,
				}
				arr = append(arr, a)
				break
			}
		}

	}
	avgProfitRate := float64(tradeCost) / profileSum * 100 // 日均收益
	//avgProfitRate := profitSum / float64(n) * 100             // 日均收益
	avgAnnualized := avgProfitRate * (float64(n) / 250) * 100 //年化收益
	avgProgress := progressSum / float64(n)                   // 日均完成率
	profitDayRate := float64(profitDays) / float64(n) * 100   // 盈利天占比

	rsp := proto.WinRatioRsp{
		Code: 200,
		Msg:  "success",
		WinHead: &proto.WinRatioHead{
			AlgoName:         algoName,
			StartDay:         TransDay(int(in.GetStartDay())),
			EndDay:           TransDay(int(in.GetEndDay())),
			TradeDays:        int32(n),
			AvgDailyProfit:   avgProfitRate,
			AnnualizedProfit: avgAnnualized,
			//TotalProfit:      0.00, //TODO: 从总线取--先不做
			MaxWithdraw:   maxWithdraw,
			DailyProgress: avgProgress,
			DailyStocks:   tools.ScoreRound(float64(stockCount) / float64(n)),
			ProfitDays:    profitDays,
			ProfitDayRate: profitDayRate,
		},
		Info: arr,
	}
	l.Logger.Infof("get rsp:%+v", rsp)
	return &rsp, nil
}

func getMaxWithdraw(max, e float64) float64 {
	if max >= e {
		return max
	}
	return e
}

// getOdds 胜率 * 收益率/(1-胜率）
func getOdds(plus, total int, profitRate float64) float64 {
	if total <= 0 {
		return 0.00
	}
	winRate := float64(plus) / float64(total)
	deno := 1 - winRate
	if deno == 0 {
		return 0.00
	}
	odds := winRate * profitRate / deno
	return odds
}

func TransDay(d int) string {
	if d <= 0 {
		return ""
	}
	day := strconv.Itoa(d)
	t, err := time.Parse(global.TimeFormatDay, day)
	if err != nil {
		return ""
	}
	out := t.Format(global.TimeFormatDaySp)
	return out
}
