// Package consumer
/*
 Author: hawrkchen
 Date: 2022/10/24 17:20
 Desc:
*/
package consumer

import (
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"math"
	"strings"
)

type ProfileResults struct {
	TradeVol   int64 // 总交易量
	LastQty    int64 // 总成交数量
	Charge     int64 // 总手续费
	CrossFee   int64 // 总流量费
	BuyVol     int64 // 总买入量
	SellVol    int64 // 总卖出量
	EntrustQty int64 // 总委托量
	DealQty    int64 // 总成交量
	CancelQty  int64 // 总撤单量

	MinJointRate float64 // 最小贴合度
	DealEffi     float64 // 成交效率
	AlgoOrderFit float64 // 母单贴合度
	TradeVolFit  float64 // 成交量贴合度

	AssessFactor float64 // 绩效收益率

	ProfitAmount  float64 // 盈亏金额
	CostAmount    int64   // 智能委托交易成本
	CostT0Amount  int64   // 日内回转交易成本
	VwapSum       float64 // 滑点总和
	VwapStdDevSum float64 // 滑点标准差总和

	PfRateStdDevSum   float64 // 收益率标准差总和
	TradeVolFitStdSum float64 // 成交量贴合度标准差总和
	TimeFitStdSum     float64 // 时间贴合度标准差总和

	Withdraw float64 // 回撤比例

	Count          int64 // 总个数 (统计单个算法下的股票交易个数)
	TradeCountPlus int64 // 正盈亏交易次数

}

/*
计算标准差一般分四步：  计算平均值， 计算方差，计算平均方差，计算标准差
如针对集合{2, 3, 4，5，7，9}， 计算平均值  sum(2+3+4+5+7+9)/6 = 5
计算方差 :  (2-5)^2= -3^3 = 9
    (3-5)^2 = 4
    (4-5)^2 = 1
	(5-5)^2 = 0
    (7-5)^2 = 4
    (9-5)^2 = 16

计算平均方差： （9+4+1+0+4+16）/6 = 34/6 = 5.6666
计算标准差： 直接开根号 5.6666开根号 = 2.38

*/

// CalVariance 计算方差，
// 入参：[]float为参与计算的列表，{a1,a2,a3,...,an},  float64为元素的平均值，avg
// 出参： float64 返回方差值
// 计算公式 (a1-avg)^2 + (a2-avg)^2 + ...(an-avg)^2
func CalVariance(an []float64, avg float64) float64 {
	var sum float64
	for _, v := range an {
		sum += math.Pow(v-avg, 2)
	}
	return sum
}

// CalFactor 计算绩效因子
func CalFactor(vwapStdDev, arrivePrice, avgPrice float64, side int) float64 {
	var factor float64
	t := math.Abs(arrivePrice - avgPrice)
	if vwapStdDev == 0 || t == 0 { // in.VwapStdDev 为分子， t 为分母， 当两个任一为0时，此时绩效因子固定为1
		factor = 1
		return factor
	} else {
		// 区分买方和卖方
		cv := vwapStdDev / t // 滑点变异系数
		if cv > 0 {
			if side == global.TradeSideBuy { // 买方
				if arrivePrice > avgPrice { // 市场均价 大于 成交均价， 绩效因子 1/cv
					factor = 1 / cv
				} else {
					factor = cv
				}
			} else if side == global.TradeSideSell { // 卖方
				if arrivePrice < avgPrice { // 市场均价 小于 成交均价， 绩效因子 1/cv
					factor = 1 / cv
				} else {
					factor = cv
				}
			} else { // 未知的买卖方向
				logx.Error("unknown trade side")
			}
		} else {
			factor = 1
		}
	}
	return factor
}

// ProfileContains 判别单个证券画像是否在算法中
// s 为画像的key,  sub 为汇总的key----加前缀
// Nor:1:20230425:1:2:000001:4
// Nor:20230425:1:2
// 原始订单 Ori:12345:20230425:1:20:000001:4
// Ori:12345:20230425:1:20
func ProfileContains(s, sub string) bool {
	if len(s) == 0 || len(sub) == 0 {
		return false
	}
	arr := strings.Split(s, ":")
	arrs := strings.Split(sub, ":")
	if len(arr) != 7 || len(arrs) != 5 {
		return false
	}
	if arr[0] == arrs[0] && arr[1] == arrs[1] && arr[2] == arrs[2] && arr[3] == arrs[3] && arr[4] == arrs[4] {
		return true
	}
	return false
}

// GetTradeVolRate 计算交易量
func GetTradeVolRate(amount int64, vol *global.TradeVolRate) {
	if amount >= 10000000000 { // 百万以上
		vol.Billion += amount
	} else if amount > 100000000 && amount < 10000000000 { // 万元以上百万以下
		vol.Million += amount
	} else { // 万元以下
		vol.Thousand += amount
	}
}

// GetOrderNum 取订单数量，本地Cache取母单保存的订单数据
// fixFlag = 0 为正常数据，  fixFlag = 1 为修复数据
func GetOrderNum(redisClient *redis.Redis, fixFlag int, date int64, userId string, algoId int) int {
	var basketKey string
	if fixFlag == 0 {
		basketKey = fmt.Sprintf("%s:%d:%s:%d", global.BasketsPrx, date, userId, algoId)
	} else {
		basketKey = fmt.Sprintf("%s:%d:%s:%d", global.BasketsFixPrx, date, userId, algoId)
	}
	orderNum, err := redisClient.Scard(basketKey)
	if err != nil {
		logx.Error("redis Scard  error:", err)
	}
	/*
		global.GAlgoOrderBasket.RWMutex.Lock()
		orderNum := len(global.GAlgoOrderBasket.Baskets[basketKey])
		global.GAlgoOrderBasket.RWMutex.Unlock()
	*/

	return int(orderNum)
}

// GetMinJointRate 计算最小贴合度
// 附加返回一个母单委托时间，计算成交效率用到
func GetMinJointRate(redis *redis.Redis, progress float64, data *global.ChildOrderData) (float64, int64) {
	// 算最小贴合度
	var theoryProgress float64
	//durTime := tools.GetCurDurationTime()
	durTime := tools.GetDurationTime(data.UnixTime)
	// redis 取母单开始时间和结束时间
	algoOrderId := fmt.Sprintf("%s:%s:%d", data.SourcePrx, global.AlgoOrderIdPrx, data.AlgoOrderId)
	str, err := redis.Hmget(algoOrderId, "startTime", "endTime", "transTime")
	if err != nil {
		logx.Error("redis get algo starttime &endtime error:", err)
		return 0.00, 0
	}
	if len(str) < 2 {
		logx.Error("algo starttime &endtime lack")
		return 0.00, 0
	}
	theoryStartTime, theoryEndTime, enTime := cast.ToInt64(str[0]), cast.ToInt64(str[1]), cast.ToInt64(str[2])

	entrustTime := tools.GetDurationTime(cast.ToInt64(enTime)) // 母单实际委托时间
	logx.Info("in GetMinJointRate, startTime:", theoryStartTime, ", endTime:", theoryEndTime,
		",transTime:", durTime, ",entrustTime:", entrustTime, ", progress:", progress)
	// 开始时间取母单理论开始时间和母单委托时间的最大值
	startTime := tools.MaxInt(theoryStartTime, entrustTime)
	divTime := theoryEndTime - startTime
	if divTime > 0 {
		theoryProgress = float64(durTime-startTime) / float64(divTime)
	}
	var curMinJointRate, retMinJointRate float64
	// progress 传进来都是*100了，这里需要再除以100
	if theoryProgress != 0 {
		curMinJointRate = (progress / 100) / theoryProgress
	}
	if curMinJointRate > 1 { // 比理论的快，该值大于1，再除1，使其小于1
		retMinJointRate = float64(1) / curMinJointRate
	} else {
		retMinJointRate = math.Abs(curMinJointRate) // curMinJointRate 有可能出现负数
	}
	logx.Info("get retMinJointRate:", retMinJointRate)
	return retMinJointRate, enTime
}

// GetFundRate 计算资金占比
func GetFundRate(fundType int, amount int64, fund *global.FundRate, data *global.ChildOrderData) {
	// 根据证券ID取其市值属性
	switch fundType {
	case 1:
		fund.Huge += amount
	case 2:
		fund.Big += amount
	case 3:
		fund.Middle += amount
	case 4:
		fund.Small += amount
	}
}

// GetTradeVolDirect 计算买卖方向
func GetTradeVolDirect(amount int64, dirt *global.TradeVolDirect, data *global.ChildOrderData) {
	if data.TradeSide == global.TradeSideBuy {
		dirt.BuyVol += amount
	} else {
		dirt.SellVol += amount
	}
}

// GetStockType 计算股价类型
func GetStockType(stockType int, t *global.StockType, data *global.ChildOrderData) {
	switch stockType {
	case 1:
		t.Red++
	case 2:
		t.Orange++
	case 3:
		t.Yellow++
	case 4:
		t.Green++
	case 5:
		t.Cyan++
	case 6:
		t.Blue++
	case 7:
		t.Purple++
	default:
		t.Red++
	}
}
