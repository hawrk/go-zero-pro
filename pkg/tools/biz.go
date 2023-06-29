// Package tools
/*
 Author: hawrkchen
 Date: 2022/4/24 19:30
 Desc:
*/
package tools

import (
	"algo_assess/global"
	"fmt"
	"github.com/spf13/cast"
	"math"
	"strconv"
	"strings"
	"time"
)

// TimeMoveForward 时间归到下一分钟
func TimeMoveForward(s string) int64 {
	times, _ := time.ParseInLocation(global.TimeFormatMinInt, s, time.Local)
	// 时间归后，比如9:30:00 的数据，统一归到9:31:00 里，要加上一分钟
	stMoveFoward := time.Unix(times.Unix(), 0).Add(time.Minute).Format(global.TimeFormatMinInt)
	t, _ := strconv.ParseInt(stMoveFoward, 10, 64)
	return t
}

// BuildMarketTime 十档行情数据时间更新为 200601061504时间格式
func BuildMarketTime(u uint64, fmDay string) int64 {
	str := strconv.FormatUint(u, 10)
	var out string
	if len(str) == 5 || len(str) == 8 {
		out = fmDay + "0" + str[:3]
	} else if len(str) == 6 || len(str) == 9 {
		out = fmDay + str[:4]
	}
	t := TimeMoveForward(out)
	return t
}

// Byte32ArrayToStringForSz 行情十档数据32位的整型数组转换为[,]分隔的字符串
func Byte32ArrayToStringForSz(data [10]uint32) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(DivTenThousand(v))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

// Byte64ArrayToStringForSz 行情十档数据64位的整型数组转换为[,]分隔的字符串
func Byte64ArrayToStringForSz(data [10]uint64) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(DivHundred(v))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

// Byte32ArrayToStringForSh
func Byte32ArrayToStringForSh(data [10]uint32) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(DivTenThousand(v))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

func Byte64ArrayToStringForSh(data [10]uint64) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(DivHundred(v))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

// DivMillion 整数除以100万转换成浮点数
func DivMillion(d uint32) string {
	out := fmt.Sprintf("%.2f", float64(d)/1000000)
	return out
}

func DivHundred(d uint64) string {
	return strconv.FormatInt(int64(d/100), 10)
}

// DivTenThousand 整数除以1万转换成浮点数
func DivTenThousand(d uint32) string {
	out := fmt.Sprintf("%.2f", float64(d)/10000)
	return out
}

// DivThousandAccount 整数除以1000 转换成浮点数   金额转换
func DivThousandAccount(d uint32) string {
	out := fmt.Sprintf("%.2f", float64(d)/1000)
	return out
}

// DivThousandCount 整数除以1000 转换成浮点数   数量转换
func DivThousandCount(d uint64) string {
	return strconv.FormatInt(int64(d/1000), 10)
}

// DivTenThousandCount 整数除以10000， 数量转换
func DivTenThousandCount(d uint64) string {
	return strconv.FormatInt(int64(d/10000), 10)
}

// GetTimePoint  把 202206211530 格式的时间转换为15:30 的格式
func GetTimePoint(s string) string {
	if len(s) < 12 {
		return ""
	}
	o := s[8:10] + ":" + s[10:]

	if o == "11:30" || o == "13:00" { // 两个时间点需要合并
		o = "11:30/13:00"
	}
	return o
}

// GetTimePointByDay 把20220621 格式的时间转换为06/21
func GetTimePointByDay(s string) string {
	if len(s) < 8 {
		return ""
	}
	return s[4:6] + "/" + s[6:]
}

// Time2String time 格式转换成字符串
func Time2String(time2 time.Time) string {
	if time2.IsZero() {
		return ""
	}
	return time2.Format(global.TimeFormat)
}

// TimeMini2String 把202206211530 格式时间转换成 2022-06-21 15:30:00 格式的字符串
func TimeMini2String(t int64) string {
	times, _ := time.ParseInLocation("200601021504", cast.ToString(t), time.Local)
	return times.Format(global.TimeFormat)
}

// TimeStr2TimeMicro 把202206211530 格式时间转换成 unix时间戳1655796600000000
func TimeStr2TimeMicro(t string) int64 {
	times, _ := time.ParseInLocation("200601021504", t, time.Local)
	return times.UnixMicro()
}

// TimeStr2TimeUnix 把202206211530 格式时间转换成unix时间戳 1655796600  （精确到秒)
func TimeStr2TimeUnix(t string) int64 {
	times, _ := time.ParseInLocation("200601021504", t, time.Local)
	return times.Unix()
}

// TimeDay2string 20220505转换成 2022-05-05格式
func TimeDay2string(t string) string {
	times, _ := time.ParseInLocation("20060102", t, time.Local)
	return times.Format("2006-01-02")
}

// GetCurDurationTime 取当天已持续的时间，计算当天0点过后经过的秒数，如当前时间为凌晨一点时，持续的时间为 3600
func GetCurDurationTime() int64 {
	// 当天0点的时间戳
	t := time.Now()
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	// 当前时间
	curTime := time.Now().Unix()
	dur := curTime - addTime
	return dur
}

// GetDurationTime 取当天指定时间持续的时间，算当天0点过后经过的秒数
func GetDurationTime(curTime int64) int64 {
	// 当天0点的时间戳
	t := time.Now()
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	dur := curTime - addTime
	return dur
}

// ScoreRound 四舍五入
// 浮点型转换成整型
func ScoreRound(score float64) int32 {
	return int32(math.Round(score))
}

// GetOrderTradeSide 买入卖出统一转换
// BUY = '1',               // 普通买入
// SELL = '2'               // 普通卖出
// MTRADE_BUY = '3',                              // 证券、基金、债券融资买入
// MTRADE_SELL = '4',                            // 证券、基金、债券融券卖出
// MTRADE_MARGIN_BUY = '5',             // 证券、基金、债券担保品融资买入
// MTRADE_MARGIN_SELL = '6',            // 证券、基金、债券担保品融券卖出
// MTRADE_BUYREPAY = '7',                  // 买券还券
// MTRADE_SELLREPAY = '8',                 // 卖券还款
// MTRADE_SELLREPAY_SNO = '9',       // 卖券还款指定合约
// NEWSECURITY_BUY = 'a',                   // 新股申购
func GetOrderTradeSide(side uint32) int {
	if side == 1 || side == 49 || side == 51 || side == 53 || side == 55 || side == 58 {
		return 1
	} else if side == 2 || side == 50 || side == 52 || side == 54 || side == 56 || side == 57 {
		return 2
	} else {
		return 0
	}
}
