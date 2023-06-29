package provider

import (
	bytesUtil "algo_assess/assess-rpc-server/internal/util"
	"algo_assess/global"
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

func BuildTagTagQuoteClientLevel2Data(retData global.TagQuoteClientLevel2Data) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, &retData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ParseFileLine(line string, market int) *global.TagQuoteClientLevel2Data {
	options := strings.Split(line, "|")
	var SecIDStr string
	if market == 1 {
		SecIDStr = "sh:" + options[0]
	} else if market == 2 {
		SecIDStr = "sz:" + options[0]
	}
	//PTimestampStr := options[1]
	OrigTimeStr := options[2]
	LastPriceStr := options[3]
	AskPriceStr := options[4]
	AskVolStr := options[5]
	BidPriceStr := options[6]
	BidVolStr := options[7]
	TotalTradeNumStr := options[8]
	TotalTradeVolStr := options[9]
	TotalTradeValueStr := options[10]
	TradeStatusStr := options[11]
	MsgIDStr := options[12]
	fmt.Printf("MsgIDStr:%v\n:", MsgIDStr)
	secIdArr := bytesUtil.BytesExtend(bytesUtil.StringToBytes(SecIDStr), 8)
	var SecID [8]byte
	copy(SecID[:], secIdArr[:8])

	OrigTime, _ := strconv.ParseInt(OrigTimeStr, 10, 64)
	LastPrice, _ := strconv.ParseInt(LastPriceStr, 10, 64)

	AskPriceStrs := strings.Split(AskPriceStr, ",")
	var AskPrice [10]uint32
	for i, str := range AskPriceStrs {
		u, _ := strconv.ParseInt(str, 10, 32)
		AskPrice[i] = uint32(u)
	}

	AskVolStrs := strings.Split(AskVolStr, ",")
	var AskVol [10]uint64
	for i, str := range AskVolStrs {
		u, _ := strconv.ParseInt(str, 10, 64)
		AskVol[i] = uint64(u)
	}

	BidPriceStrs := strings.Split(BidPriceStr, ",")
	var BidPrice [10]uint32
	for i, str := range BidPriceStrs {
		u, _ := strconv.ParseInt(str, 10, 32)
		BidPrice[i] = uint32(u)
	}

	BidVolStrs := strings.Split(BidVolStr, ",")
	var BidVol [10]uint64
	for i, str := range BidVolStrs {
		u, _ := strconv.ParseInt(str, 10, 64)
		BidVol[i] = uint64(u)
	}

	TotalTradeNum, _ := strconv.ParseInt(TotalTradeNumStr, 10, 64)

	TotalTradeVol, _ := strconv.ParseInt(TotalTradeVolStr, 10, 64)

	TotalTradeValue, _ := strconv.ParseInt(TotalTradeValueStr, 10, 64)

	TradeStatus, _ := strconv.ParseInt(TradeStatusStr, 10, 1)

	MsgIDStr = strings.Trim(MsgIDStr, "\x00")
	MsgID, err := strconv.ParseInt(MsgIDStr, 10, 64)
	if err != nil {
	}

	return &global.TagQuoteClientLevel2Data{
		QuoteHead: global.QuoteHead{
			QuoteID:   uint32(MsgID),
			QuoteType: 4,
		},
		MsgID:           uint32(MsgID),
		SecID:           SecID,
		Market:          uint8(market),
		OrigTime:        uint64(OrigTime),
		TradeStatus:     uint8(TradeStatus),
		LastPrice:       uint32(LastPrice),
		AskPrice:        AskPrice,
		AskVol:          AskVol,
		BidPrice:        BidPrice,
		BidVol:          BidVol,
		TotalTradeNum:   uint32(TotalTradeNum),
		TotalTradeVol:   uint64(TotalTradeVol),
		TotalTradeValue: uint64(TotalTradeValue),
		TotalBidVol:     0,
		TotalAskVol:     0,
	}
}
