package global

const (
	Level1Level           = 5
	Level2Level           = 10
	OptionLevel           = 5
	MaxSupportMarketCount = 16
	SecIdLen              = 8
	VOrderLen             = 50
)

//  SubscribeQuoteType 订阅行情数据类型
type SubscribeQuoteType uint16

const (
	SubscribeQuoteTypeUnknown      SubscribeQuoteType = 0x00 // 未知
	SubscribeQuoteTypeIndex        SubscribeQuoteType = 0x01 // 指数快照
	SubscribeQuoteTypeLevel1       SubscribeQuoteType = 0x02 // 五档行情快照(Level1)
	SubscribeQuoteTypeLevel2       SubscribeQuoteType = 0x04 // 十档行情快照(Level2)
	SubscribeQuoteTypeEntrust      SubscribeQuoteType = 0x08 // 逐笔委托
	SubscribeQuoteTypeTurnover     SubscribeQuoteType = 0x10 // 逐笔成交
	SubscribeQuoteTypeEntrustQueue SubscribeQuoteType = 0x20 // 委托队列
	SubscribeQuoteTypeOption       SubscribeQuoteType = 0x40 // 期权
	SubscribeQuoteTypeAll          SubscribeQuoteType = 0x7F // 所有类型
)

type QuoteHead struct {
	QuoteID    uint32 // 行情ID 不同消息类型从0开始递增
	SecBitType uint16 // 业务类型 参考SecBitTypeInfo
	QuoteType  uint16 // 行情数据类型
}

//  TagQuoteClientIndexData 指数数据结构
type TagQuoteClientIndexData struct {
	QuoteHead            //  行情头部消息
	MsgID         uint64 // 消息ID 不同消息类型从0开始递增
	SecID         string // 交易所证券代码（六位十六进制表示）
	Market        uint16 // 市场
	OrigTime      uint64 // 数据生成时间
	OpenIndex     uint64 // 今开指数N(6)/ N(5)
	HighIndex     uint64 // 最高指数N(6)/ N(5)
	LowIndex      uint64 // 最低指数N(6)/ N(5)
	LastIndex     uint64 // 最新指数N(6)/ N(5)
	TotalVol      uint64 // 参与计算的证券总交易数量（单位：手）N(2)/ N(5)
	Turnover      uint64 // 参与计算的证券总成交金额（单位：百元）N(4)/
	PreCloseIndex uint64 // 昨收指数N(6)/ N(5)
}

//  TagQuoteClientLevel1Data 五档行情数据结构
type TagQuoteClientLevel1Data struct {
	QuoteHead            //  行情头部消息
	MsgID         uint64 // 消息ID 不同消息类型从0开始递增
	SecID         string // 交易所证券代码（六位十六进制表示）
	Market        uint16 // 市场
	OrigTime      uint64 // 数据生成时间
	TradeStatus   uint8  // 交易状态 见字典
	PreClosePrice uint32 // 昨收盘价
	OpenPrice     uint32 // 开盘价
	HighPrice     uint32 // 最高价
	LowPrice      uint32 // 最低价
	LastPrice     uint32 // 最新价N(6)/ N(3)

	AskPrice [Level1Level]uint32 // 申卖价N(6)/ N(3)  LEVEL1_LEVEL 32
	AskVol   [Level1Level]uint64 // 申卖量N(2)/ N(3)  LEVEL1_LEVEL 64
	BidPrice [Level1Level]uint32 // 申买价N(6)/ N(3))  LEVEL1_LEVEL 32
	BidVol   [Level1Level]uint64 // 申买量N(2)/ N(3)  LEVEL1_LEVEL 64

	TotalTradeNum     uint32 // 成交笔数
	TotalTradeVol     uint64 // 成交总量N(2)/ N(3)
	TotalTradeValue   uint64 // 成交总金额（单位：元）N(4)/ N(5)
	TotalBidVol       uint64 // 委托买入总量N(2)/ N(3)
	TotalAskVol       uint64 // 委托卖出总量N(2)/ N(3)
	WeightAvgBidPrice uint32 // 加权平均委托买入价
	WeightAvgAskPrice uint32 // 加权平均委托卖出价
	IOPV              uint32 // IOPV净值估值
	YieldToMaturity   int32  // 到期收益率
	HighLimited       uint32 // 涨停价
	LowLimited        uint32 // 跌停价
}

//  TagQuoteClientLevel2Data 十档行情数据结构
type TagQuoteClientLevel2Data struct {
	QuoteHead                    //  行情头部消息
	MsgID         uint32         // 消息ID 不同消息类型从0开始递增
	SecID         [SecIdLen]byte // 交易所证券代码（六位十六进制表示）
	Market        uint8          // 市场
	OrigTime      uint64         // 数据生成时间
	TradeStatus   uint8          // 交易状态 见字典
	PreClosePrice uint32         // 昨收盘价
	OpenPrice     uint32         // 开盘价
	HighPrice     uint32         // 最高价
	LowPrice      uint32         // 最低价
	LastPrice     uint32         // 最新价N(6)/ N(3)

	AskPrice          [Level2Level]uint32 // 申卖价N(6)/ N(3) LEVEL2_LEVEL 32
	AskVol            [Level2Level]uint64 // 申卖量N(2)/ N(3) LEVEL2_LEVEL 64
	BidPrice          [Level2Level]uint32 // 申买价N(6)/ N(3) LEVEL2_LEVEL 32
	BidVol            [Level2Level]uint64 // 申买量N(2)/ N(3) LEVEL2_LEVEL 64
	TotalTradeNum     uint32              // 成交笔数
	TotalTradeVol     uint64              // 成交总量N(2)/ N(3)
	TotalTradeValue   uint64              // 成交总金额（单位：元）N(4)/ N(5)
	TotalBidVol       uint64              // 委托买入总量N(2)/ N(3)
	TotalAskVol       uint64              // 委托卖出总量N(2)/ N(3)
	WeightAvgBidPrice uint32              // 加权平均委托买入价
	WeightAvgAskPrice uint32              // 加权平均委托卖出价
	IOPV              uint32              // IOPV净值估值
	YieldToMaturity   int32               // 到期收益率
	HighLimited       uint32              // 涨停价
	LowLimited        uint32              // 跌停价
	OrigDate          uint32              // 日期
}

//  TagQuoteClientEntrustData 逐笔委托数据结构
type TagQuoteClientEntrustData struct {
	QuoteHead           //  行情头部消息
	MsgID        uint64 // 消息ID 不同消息类型从0开始递增
	SecID        string // 交易所证券代码（六位十六进制表示）
	Market       uint16 // 市场
	OrigTime     uint64 // 数据生成时间
	TradeChannel uint32 // 交易通道号
	Order        uint32 // 委托号
	Side         byte   // 买卖方向
	Price        uint32 // 委托价格N(4)/ N(3)
	Volume       uint32 // 委托数量N(2)/ N(3)
	OrderType    uint8  // 委托类别
	OrderNo      uint64 // （沪市）原始订单号
	BizIndex     uint64 // （沪市）业务序列号，与逐笔成交统一编号
}

//  TagQuoteClientTurnoverData 逐笔成交数据结构
type TagQuoteClientTurnoverData struct {
	QuoteHead               //  行情头部消息
	MsgID            uint64 // 消息ID 不同消息类型从0开始递增
	SecID            string // 交易所证券代码（六位十六进制表示）
	Market           uint16 // 市场
	OrigTime         uint64 // 数据生成时间
	TradeChannel     uint32 // 交易通道号
	TransactionIndex uint32 // 消息记录号
	Price            uint32 // 成交价N(4)/ N(3)
	Volume           uint32 // 成交数量N(2)/ N(3)
	BSFlag           uint8  // 买卖方向
	TransactionType  uint8  // 成交类别
	AskOrder         uint32 // 卖方委托序号（深市）卖方委托号，（沪市）卖方原始订单号
	BidOrder         uint32 // 买方委托序号（深市）买方委托号，（沪市）买方原始订单号
	BizIndex         uint64 // （沪市）业务序列号，与逐笔委托统一编号
}

//  TagQuoteClientEntrustQueueData 委托队列数据结构
type TagQuoteClientEntrustQueueData struct {
	QuoteHead                       //  行情头部消息
	MsgID         uint64            // 消息ID 不同消息类型从0开始递增
	SecID         string            // 交易所证券代码（六位十六进制表示）
	Market        uint16            // 市场
	OrigTime      uint64            // 数据生成时间
	Side          byte              // 买卖方向
	Price         uint32            // 委托价格N(6)/ N(3)
	Orders        uint32            // 委托总笔量N(2)/ N(3)
	OrderVolCount uint32            // 委托总笔数
	VOrder        [VOrderLen]uint32 // 委托明细(最多展示50个)  len=50  32
}

//  TagQuoteClientOptionsDataSh 期权数据结构
type TagQuoteClientOptionsDataSh struct {
	QuoteHead                           //  行情头部消息
	MsgID           uint64              // 消息ID 不同消息类型从0开始递增
	SecID           string              // 交易所证券代码（六位十六进制表示）
	Market          uint16              // 市场
	OrigTime        uint64              // 数据生成时间
	TradeStatus     uint8               // 交易状态  见字典
	PreOpenInterest uint64              // 昨持仓
	PreClosePrice   uint32              // 昨收盘价
	PreSettlePrice  uint32              // 昨结算价
	OpenPrice       uint32              // 开盘价
	HighPrice       uint32              // 最高价
	LowPrice        uint32              // 最低价
	LastPrice       uint32              // 最新价N(6)/
	AskPrice        [OptionLevel]uint32 // 申卖价N(6)/  OPTION_LEVEL 32
	AskVol          [OptionLevel]uint64 // 申卖量N(2)/  OPTION_LEVEL 64
	BidPrice        [OptionLevel]uint32 // 申买价N(6)/  OPTION_LEVEL 32
	BidVol          [OptionLevel]uint64 // 申买量N(2)/  OPTION_LEVEL 64
	TotalTradeNum   uint32              // 成交笔数（张）
	TotalTradeVol   uint64              // 成交总量N(2)/
	TotalTradeValue uint64              // 成交总金额（单位：元）N(4)/
	OpenInterest    int64               // 持仓总量N(2)/
	ClosePrice      uint32              // 今收盘价
	SettlePrice     uint32              // 今结算价
	PreDelta        int32               // 昨虚实度
	CurrDelta       int32               // 今虚实度
	HighLimited     uint32              // 涨停价
	LowLimited      uint32              // 跌停价
}
