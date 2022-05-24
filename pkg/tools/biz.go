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
	"strconv"
	"strings"
	"time"
)

// BuildMarketTime 十档行情数据时间更新为 200601061504时间格式
func BuildMarketTime(u uint64, fmDay string) int64 {
	str := strconv.FormatUint(u, 10)
	var out string
	if len(str) == 5 || len(str) == 8 {
		out = fmDay + "0" + str[:3]
	} else if len(str) == 6 || len(str) == 9 {
		out = fmDay + str[:4]
	}
	st, _ := strconv.ParseInt(out, 10, 64)
	return st
}

func GetDayByEnv(env string, t uint64) (origTime int64) {
	if env == global.EnvPro || env == global.EnvGray { // 取当天日期
		today := time.Now().Format(global.TimeFormatDay)
		origTime = BuildMarketTime(t, today)
	} else if env == global.EnvTest || env == global.EvnDev { // 取前一天
		preday := time.Now().AddDate(0, 0, -1).Format(global.TimeFormatDay)
		origTime = BuildMarketTime(t, preday)
	}
	return
}

// Byte32ArrayToString 行情十档数据32位的整型数组转换为[,]分隔的字符串
func Byte32ArrayToStringForSz(data [10]uint32) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(DivMillion(v))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

// Byte64ArrayToString 行情十档数据64位的整型数组转换为[,]分隔的字符串
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
		build.WriteString(DivThousandAccount(v))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

func Byte64ArrayToStringForSh(data [10]uint64) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(DivThousandCount(v))
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

// DivThousand 整数除以1000 转换成浮点数   金额转换
func DivThousandAccount(d uint32) string {
	out := fmt.Sprintf("%.2f", float64(d)/1000)
	return out
}

// DivThousand 整数除以1000 转换成浮点数   数量转换
func DivThousandCount(d uint64) string {
	return strconv.FormatInt(int64(d/1000), 10)
}
