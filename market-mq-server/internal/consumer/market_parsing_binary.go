package consumer

import (
	"algo_assess/global"
	"bytes"
	"encoding/binary"
	"log"
)

func ParseQuoteHead(v string) (head global.QuoteHead, err error) {
	buf := bytes.NewReader([]byte(v))
	if err := binary.Read(buf, binary.LittleEndian, &head); err != nil {
		return global.QuoteHead{}, err
	}
	return head, nil
}

// 解析十档行情数据
func parsTagTagQuoteClientLevel2Data(val string) (retData global.TagQuoteClientLevel2Data, err error) {
	buf := bytes.NewReader([]byte(val))
	if err := binary.Read(buf, binary.LittleEndian, &retData); err != nil {
		return global.TagQuoteClientLevel2Data{}, err
	}
	return retData, nil
}

// 解析指数行情
func parsTagQuoteClientIndexData(lhsPBuffer *bytes.Buffer, head *global.QuoteHead) (interface{}, error) {
	log.Println("指数行情解析")
	var retData global.TagQuoteClientIndexData
	retData.QuoteID = head.QuoteID
	retData.SecBitType = head.SecBitType
	retData.QuoteType = head.QuoteType
	// 读数据
	err := binary.Read(lhsPBuffer, binary.LittleEndian, &retData.MsgID)
	if err != nil {
		return nil, err
	}
	SecIDTemp := make([]byte, global.SecIdLen)
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &SecIDTemp)
	if err != nil {
		return nil, err
	}
	retData.SecID = string(SecIDTemp[:])
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Market)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrigTime)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OpenIndex)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.HighIndex)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LowIndex)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LastIndex)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalVol)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Turnover)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.PreCloseIndex)
	if err != nil {
		return nil, err
	}
	return retData, nil
}

// 解析五档行情数据
func parsTagQuoteClientLevel1Data(lhsPBuffer *bytes.Buffer, head *global.QuoteHead) (interface{}, error) {
	log.Println("五档行情解析")
	var retData global.TagQuoteClientLevel1Data
	retData.QuoteID = head.QuoteID
	retData.SecBitType = head.SecBitType
	retData.QuoteType = head.QuoteType
	// 读数据
	err := binary.Read(lhsPBuffer, binary.LittleEndian, &retData.MsgID)
	if err != nil {
		return nil, err
	}
	// 交易所证券代码
	lSecIDTemp := make([]byte, global.SecIdLen)
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &lSecIDTemp)
	if err != nil {
		return nil, err
	}
	retData.SecID = string(lSecIDTemp[:])
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Market)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrigTime)
	if err != nil {
		return nil, err
	}

	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TradeStatus)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.PreClosePrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OpenPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.HighPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LowPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LastPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.AskPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.AskVol)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BidPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BidVol)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalTradeNum)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalTradeVol)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalTradeValue)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalBidVol)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalAskVol)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.WeightAvgBidPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.WeightAvgAskPrice)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.IOPV)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.YieldToMaturity)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.HighLimited)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LowLimited)
	if err != nil {
		return nil, err
	}
	return retData, nil
}

// 解析逐笔委托行情数据
/*
func parsTagQuoteClientEntrustData(lhsPBuffer *bytes.Buffer, head *global.QuoteHead) interface{} {
	log.Println("逐笔委托行情行情解析")
	var retData global.TagQuoteClientEntrustData
	retData.QuoteID = head.QuoteID
	retData.SecBitType = head.SecBitType
	retData.QuoteType = head.QuoteType
	// 读数据
	err := binary.Read(lhsPBuffer, binary.LittleEndian, &retData.MsgID)
	if err != nil {
		log.Fatal("read buffer MsgID failed\n")
	}
	// 交易所证券代码
	lSecIDTemp := make([]byte, global.SecIdLen)
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &lSecIDTemp)
	if err != nil {
		log.Fatal("read buffer SecID failed\n")
	}
	retData.SecID = string(lSecIDTemp[:])

	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Market)
	if err != nil {
		log.Fatal("read buffer Market failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrigTime)
	if err != nil {
		log.Fatal("read buffer OrigTime failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TradeChannel)
	if err != nil {
		log.Fatal("read buffer TradeChannel failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Order)
	if err != nil {
		log.Fatal("read buffer Order failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Side)
	if err != nil {
		log.Fatal("read buffer Side failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Price)
	if err != nil {
		log.Fatal("read buffer Price failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Volume)
	if err != nil {
		log.Fatal("read buffer Volume failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrderType)
	if err != nil {
		log.Fatal("read buffer OrderType failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrderNo)
	if err != nil {
		log.Fatal("read buffer OrderNo failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BizIndex)
	if err != nil {
		log.Fatal("read buffer BizIndex failed\n")
	}
	return retData
}

// 解析逐笔成交行情数据
func parsTagQuoteClientTurnoverData(lhsPBuffer *bytes.Buffer, head *global.QuoteHead) interface{} {
	log.Println("逐笔成交行情解析")
	var retData global.TagQuoteClientTurnoverData
	retData.QuoteID = head.QuoteID
	retData.SecBitType = head.SecBitType
	retData.QuoteType = head.QuoteType
	// 读数据
	err := binary.Read(lhsPBuffer, binary.LittleEndian, &retData.MsgID)
	if err != nil {
		log.Fatal("read buffer MsgID failed\n")
	}
	// 交易所证券代码
	lSecIDTemp := make([]byte, global.SecIdLen)
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &lSecIDTemp)
	if err != nil {
		log.Fatal("read buffer SecID failed\n")
	}
	retData.SecID = string(lSecIDTemp[:])
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Market)
	if err != nil {
		log.Fatal("read buffer Market failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrigTime)
	if err != nil {
		log.Fatal("read buffer OrigTime failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TradeChannel)
	if err != nil {
		log.Fatal("read buffer TradeChannel failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TransactionIndex)
	if err != nil {
		log.Fatal("read buffer TransactionIndex failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Price)
	if err != nil {
		log.Fatal("read buffer Price failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Volume)
	if err != nil {
		log.Fatal("read buffer Volume failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BSFlag)
	if err != nil {
		log.Fatal("read buffer BSFlag failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TransactionType)
	if err != nil {
		log.Fatal("read buffer TransactionType failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.AskOrder)
	if err != nil {
		log.Fatal("read buffer AskOrder failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BidOrder)
	if err != nil {
		log.Fatal("read buffer BidOrder failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BizIndex)
	if err != nil {
		log.Fatal("read buffer BizIndex failed\n")
	}
	return retData
}

// 解析委托队列行情数据
func parsTagQuoteClientEntrustQueueData(lhsPBuffer *bytes.Buffer, head *global.QuoteHead) *global.TagQuoteClientEntrustQueueData {
	log.Println("委托队列行情解析")
	var retData global.TagQuoteClientEntrustQueueData
	retData.QuoteID = head.QuoteID
	retData.SecBitType = head.SecBitType
	retData.QuoteType = head.QuoteType
	// 读数据
	err := binary.Read(lhsPBuffer, binary.LittleEndian, &retData.MsgID)
	if err != nil {
		log.Fatal("read buffer MsgID failed\n")
	}
	// 交易所证券代码
	lSecIDTemp := make([]byte, global.SecIdLen)
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &lSecIDTemp)
	if err != nil {
		log.Fatal("read buffer SecID failed\n")
	}
	retData.SecID = string(lSecIDTemp[:])
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Market)
	if err != nil {
		log.Fatal("read buffer Market failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrigTime)
	if err != nil {
		log.Fatal("read buffer OrigTime failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Side)
	if err != nil {
		log.Fatal("read buffer Side failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Price)
	if err != nil {
		log.Fatal("read buffer Price failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Orders)
	if err != nil {
		log.Fatal("read buffer Orders failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrderVolCount)
	if err != nil {
		log.Fatal("read buffer OrderVolCount failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.VOrder)
	if err != nil {
		log.Fatal("read buffer VOrder failed\n")
	}
	return &retData

}

// 解析期权行情数据
func parsTagQuoteClientOptionsDataSh(lhsPBuffer *bytes.Buffer, head *global.QuoteHead) interface{} {
	log.Println("期权行情解析")
	var retData global.TagQuoteClientOptionsDataSh
	retData.QuoteID = head.QuoteID
	retData.SecBitType = head.SecBitType
	retData.QuoteType = head.QuoteType
	// 读数据
	err := binary.Read(lhsPBuffer, binary.LittleEndian, &retData.MsgID)
	if err != nil {
		log.Fatal("read buffer MsgID failed\n")
	}
	// 交易所证券代码
	lSecIDTemp := make([]byte, global.SecIdLen)
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &lSecIDTemp)
	if err != nil {
		log.Fatal("read buffer SecID failed\n")
	}
	retData.SecID = string(lSecIDTemp[:])
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.Market)
	if err != nil {
		log.Fatal("read buffer Market failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OrigTime)
	if err != nil {
		log.Fatal("read buffer OrigTime failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TradeStatus)
	if err != nil {
		log.Fatal("read buffer TradeStatus failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.PreOpenInterest)
	if err != nil {
		log.Fatal("read buffer PreOpenInterest failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.PreClosePrice)
	if err != nil {
		log.Fatal("read buffer PreClosePrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.PreSettlePrice)
	if err != nil {
		log.Fatal("read buffer PreSettlePrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OpenPrice)
	if err != nil {
		log.Fatal("read buffer OpenPrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.HighPrice)
	if err != nil {
		log.Fatal("read buffer HighPrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LowPrice)
	if err != nil {
		log.Fatal("read buffer LowPrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LastPrice)
	if err != nil {
		log.Fatal("read buffer LastPrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.AskPrice)
	if err != nil {
		log.Fatal("read buffer AskPrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.AskVol)
	if err != nil {
		log.Fatal("read buffer AskVol failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BidPrice)
	if err != nil {
		log.Fatal("read buffer BidPrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.BidVol)
	if err != nil {
		log.Fatal("read buffer BidVol failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalTradeNum)
	if err != nil {
		log.Fatal("read buffer TotalTradeNum failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalTradeVol)
	if err != nil {
		log.Fatal("read buffer TotalTradeVol failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.TotalTradeValue)
	if err != nil {
		log.Fatal("read buffer TotalTradeValue failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.OpenInterest)
	if err != nil {
		log.Fatal("read buffer OpenInterest failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.ClosePrice)
	if err != nil {
		log.Fatal("read buffer ClosePrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.SettlePrice)
	if err != nil {
		log.Fatal("read buffer SettlePrice failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.PreDelta)
	if err != nil {
		log.Fatal("read buffer PreDelta failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.CurrDelta)
	if err != nil {
		log.Fatal("read buffer CurrDelta failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.HighLimited)
	if err != nil {
		log.Fatal("read buffer HighLimited failed\n")
	}
	err = binary.Read(lhsPBuffer, binary.LittleEndian, &retData.LowLimited)
	if err != nil {
		log.Fatal("read buffer LowLimited failed\n")
	}
	return retData
}
*/

// 通用解析入口函数
/*
func parsGeneral(lhsByte []byte) interface{} {
	var head global.Head
	// 创建写缓冲
	buffer := bytes.NewBuffer([]byte{})
	// 将kafka接收的分区数据写入写缓冲中
	err := binary.Write(buffer, binary.LittleEndian, lhsByte)
	if err != nil {
		log.Fatal("write buffer failed\n")
	}
	// 解析行情数据前三个字段判断行情数据类型
	_ = binary.Read(buffer, binary.LittleEndian, &head.QuoteID)
	_ = binary.Read(buffer, binary.LittleEndian, &head.SecBitType)
	_ = binary.Read(buffer, binary.LittleEndian, &head.QuoteType)

	switch {
	// 指数快照
	case head.QuoteType == uint16(global.SubscribeQuoteTypeIndex):
		return parsTagQuoteClientIndexData(buffer, &head)
	// 五档行情快照
	case head.QuoteType == uint16(global.SubscribeQuoteTypeLevel1):
		return parsTagQuoteClientLevel1Data(buffer, head.QuoteID, head.SecBitType, head.QuoteType)
	// 十档行情快照
	case head.QuoteType == uint16(global.SubscribeQuoteTypeLevel2):
		return parsTagTagQuoteClientLevel2Data(buffer, head)
	// 逐笔委托
	case head.QuoteType == uint16(global.SubscribeQuoteTypeEntrust):
		return parsTagQuoteClientEntrustData(buffer, head.QuoteID, head.SecBitType, head.QuoteType)
	// 逐笔成交
	case head.QuoteType == uint16(global.SubscribeQuoteTypeTurnover):
		return parsTagQuoteClientTurnoverData(buffer, head.QuoteID, head.SecBitType, head.QuoteType)
	// 委托队列
	case head.QuoteType == uint16(global.SubscribeQuoteTypeEntrustQueue):
		return parsTagQuoteClientEntrustQueueData(buffer, head.QuoteID, head.SecBitType, head.QuoteType)
	// 期权
	case head.QuoteType == uint16(global.SubscribeQuoteTypeAll):
		return parsTagQuoteClientOptionsDataSh(buffer, head.QuoteID, head.SecBitType, head.QuoteType)
	default:
		fmt.Println("不存在的行情数据类型")
		return nil
	}
	return nil
}
*/
