// Package consumer
/*
 Author: hawrkchen
 Date: 2022/4/18 17:08
 Desc:
*/
package consumer

import (
	"algo_assess/global"
	"strconv"
	"strings"
	"time"
)

func BuildLevel2ToRedisHash(data global.TagQuoteClientLevel2Data) map[string]string {
	m := make(map[string]string)
	m["SecId"] = string(data.SecID[:])
	m["OpenPrice"] = strconv.Itoa(int(data.OpenPrice))
	m["HighPrice"] = strconv.Itoa(int(data.HighPrice))
	m["LowPrice"] = strconv.Itoa(int(data.LowPrice))
	m["LastPrice"] = strconv.Itoa(int(data.LastPrice))
	m["AskPrice"] = Byte32ArrayToString(data.AskPrice)
	m["AskVol"] = Byte64ArrayToString(data.AskVol)
	m["BidPrice"] = Byte32ArrayToString(data.BidPrice)
	m["BidVol"] = Byte64ArrayToString(data.BidVol)
	m["TotalTradeNum"] = strconv.Itoa(int(data.TotalTradeNum))
	m["TotalTradeVol"] = strconv.FormatUint(data.TotalTradeVol, 10)
	m["TotalTradeValue"] = strconv.FormatUint(data.TotalTradeValue, 10)
	m["TotalBidVol"] = strconv.FormatUint(data.TotalBidVol, 10)
	m["TotalAskVol"] = strconv.FormatUint(data.TotalAskVol, 10)

	return m
}

func Byte32ArrayToString(data [10]uint32) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(strconv.Itoa(int(v)))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

func Byte64ArrayToString(data [10]uint64) string {
	var build strings.Builder
	for _, v := range data {
		build.WriteString(strconv.FormatUint(v, 10))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

func BuildMarketTime(u uint64) (out string) {
	date := time.Now().Format(global.TimeFormatDay)
	str := strconv.FormatUint(u, 10)
	if len(str) == 8 {
		out = date + "0" + str[:3]
	} else if len(str) == 9 {
		out = date + str[:4]
	}
	return out
}
