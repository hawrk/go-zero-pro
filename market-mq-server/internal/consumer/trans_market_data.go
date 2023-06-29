// Package consumer
/*
 Author: hawrkchen
 Date: 2022/11/28 11:15
 Desc:
*/
package consumer

import (
	"algo_assess/global"
	pb "algo_assess/market-mq-server/proto"
	"algo_assess/pkg/tools"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"
)

// TransSzQuoteData 转换深圳交易所行情数据结构
func TransSzQuoteData(data *global.TagQuoteClientLevel2Data) global.QuoteLevel2Data {
	out := global.QuoteLevel2Data{
		SecID: "sz:" + strings.TrimSpace(tools.Bytes2String(data.SecID[:])),
		//OrigTime:      tools.GetDayByEnv(s.svcCtx.Config.Deployment.Env, data.OrigTime),
		OrigTime:      tools.BuildMarketTime(data.OrigTime, cast.ToString(data.OrigDate)),
		LastPrice:     int64(data.LastPrice) / 100, // 深圳交易所以元为单位乘以100万，这里要转成分，除以10000--hawrk 除以100
		TotalTradeVol: int64(data.TotalTradeVol) / 100,
		AskPrice:      tools.Byte32ArrayToStringForSz(data.AskPrice),
		AskVol:        tools.Byte64ArrayToStringForSz(data.AskVol),
		BidPrice:      tools.Byte32ArrayToStringForSz(data.BidPrice),
		BidVol:        tools.Byte64ArrayToStringForSz(data.BidVol),
		TotalAskVol:   int64(data.TotalAskVol) / 100,
		TotalBidVol:   int64(data.TotalBidVol) / 100,
	}
	return out
}

// TransShQuoteData 转换上交所行情数据结构
func TransShQuoteData(data *global.TagQuoteClientLevel2Data) global.QuoteLevel2Data {
	out := global.QuoteLevel2Data{
		SecID: "sh:" + strings.TrimSpace(tools.Bytes2String(data.SecID[:])),
		//OrigTime:      tools.GetDayByEnv(s.svcCtx.Config.Deployment.Env, data.OrigTime),
		OrigTime:      tools.BuildMarketTime(data.OrigTime, cast.ToString(data.OrigDate)),
		LastPrice:     int64(data.LastPrice) / 100,     // hawrk 20221216 原来 10 -> 100
		TotalTradeVol: int64(data.TotalTradeVol) / 100, // hawrk 20221216 原来 1000 -> 100
		AskPrice:      tools.Byte32ArrayToStringForSh(data.AskPrice),
		AskVol:        tools.Byte64ArrayToStringForSh(data.AskVol),
		BidPrice:      tools.Byte32ArrayToStringForSh(data.BidPrice),
		BidVol:        tools.Byte64ArrayToStringForSh(data.BidVol),
		TotalAskVol:   int64(data.TotalAskVol) / 100,
		TotalBidVol:   int64(data.TotalBidVol) / 100,
	}
	return out
}

func TransQuoteWithPB(data pb.QuoteLevel, quoteTag int) global.QuoteLevel2Data {
	// TODO: 行情数据只有当天时间，没有日期，这里先指定成当天的
	today := time.Now().Format(global.TimeFormatDay)
	var secId string
	if quoteTag == 1 {
		secId = "sz:" + data.SeculityId
	} else {
		secId = "sh:" + data.SeculityId
	}
	out := global.QuoteLevel2Data{
		SecID:         secId,
		OrigTime:      BuildPBMarketTime(data.OrgiTime, today),
		LastPrice:     int64(data.LastPrice) / 10000, // 转成分
		TotalTradeVol: int64(data.TotalTradeVol) / 100,
		AskPrice:      ParsePrice2Yuan(data.AskPrice),
		AskVol:        ParseVol2RealVol(data.AskVol),
		BidPrice:      ParsePrice2Yuan(data.BidPrice),
		BidVol:        ParseVol2RealVol(data.BidVol),
		TotalAskVol:   int64(data.TotalAskVol) / 100,
		TotalBidVol:   int64(data.TotalBidVol) / 100,
	}
	return out
}

// BuildPBMarketTime 十档行情数据时间更新为 200601061504时间格式
func BuildPBMarketTime(u int64, fmDay string) int64 {
	str := strconv.FormatInt(u, 10)
	var out string
	if len(str) == 6 || len(str) == 9 {
		out = fmDay + "0" + str[:3]
	} else if len(str) == 7 || len(str) == 10 {
		out = fmDay + str[:4]
	}
	t := tools.TimeMoveForward(out)
	return t
}

// ParsePrice2Yuan ParsePrice2Real 把18680000,18710000,18720000,18730000,18750000,18770000,18780000,18790000,18800000,18810000
// 类型字符串转换成真实价格的字符串 (除以100万）
func ParsePrice2Yuan(str string) string {
	arr := strings.Split(str, ",")
	var build strings.Builder
	for _, v := range arr {
		build.WriteString(tools.DivMillion(cast.ToUint32(v)))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}

// ParseVol2RealVol 把20000,380000,160000,400000,760000,90000,370000,350000,290000,289888申买申卖量转换成
// 真实成交量（除以100)
func ParseVol2RealVol(str string) string {
	arr := strings.Split(str, ",")
	var build strings.Builder
	for _, v := range arr {
		build.WriteString(tools.DivHundred(cast.ToUint64(v)))
		build.WriteString(",")
	}
	return build.String()[:build.Len()-1]
}
