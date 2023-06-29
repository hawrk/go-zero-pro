// Package global
/*
 Author: hawrkchen
 Date: 2022/7/13 9:38
 Desc:
*/
package global

import "sync"

// 异步写表channel
var OrderChan = make(chan MAlgoOrder, MaxChannelBuffer)
var OrderDetailChan = make(chan ChildOrderData, MaxChannelBuffer)
var ProfileChan = make(chan Profile, MaxChannelBuffer)
var ProfileSumChan = make(chan ProfileSum, MaxChannelBuffer)
var TlProfileSumChan = make(chan ProfileSum, MaxChannelBuffer)

// GNorUserProfile 普通用户画像结构体
var GNorUserProfile = MAlgoProfile{
	ProfileSumMap: make(map[string]*ProfileSum),
	TimeLineMap:   make(map[string]*ProfileSum),
	ProfileMap:    make(map[string]*Profile),
}

// GMngrProfile 管理用户画像结构体
var GMngrProfile = MAlgoProfile{
	ProfileSumMap: make(map[string]*ProfileSum),
	TimeLineMap:   make(map[string]*ProfileSum),
	ProfileMap:    make(map[string]*Profile),
}

// GProviderProfile 算法厂商画像结构体
var GProviderProfile = MAlgoProfile{
	ProfileSumMap: make(map[string]*ProfileSum),
	TimeLineMap:   make(map[string]*ProfileSum),
	ProfileMap:    make(map[string]*Profile),
}

// GAdminProfile 超级管理员画像结构体
var GAdminProfile = MAlgoProfile{
	ProfileSumMap: make(map[string]*ProfileSum),
	TimeLineMap:   make(map[string]*ProfileSum),
	ProfileMap:    make(map[string]*Profile),
}

type MAlgoProfile struct {
	ProfileSumMap map[string]*ProfileSum // key -> date + userid + algoId
	TimeLineMap   map[string]*ProfileSum // key -> transTime + userid + algoId     // 时间线展示
	ProfileMap    map[string]*Profile    // key -> date+ userid + algoId + secId+algoOrderId    // 根据证券ID母单号计算盈亏
	sync.RWMutex
}

// FundRate 资金占比结构体
type FundRate struct {
	Huge   int64 `json:"huge"`
	Big    int64 `json:"big"`
	Middle int64 `json:"middle"`
	Small  int64 `json:"small"`
}

// TradeVolDirect 买入卖出成本
type TradeVolDirect struct {
	BuyVol  int64 `json:"buy_vol"`  // 买入总量
	SellVol int64 `json:"sell_vol"` // 卖出总量
}

// StockType   股价类型
type StockType struct {
	Red    int64 `json:"red"`
	Orange int64 `json:"orange"`
	Yellow int64 `json:"yellow"`
	Green  int64 `json:"green"`
	Cyan   int64 `json:"cyan"`
	Blue   int64 `json:"blue"`
	Purple int64 `json:"purple"`
}

// TradeVolRate 交易量占比
type TradeVolRate struct {
	Billion  int64 `json:"billion"`
	Million  int64 `json:"million"`
	Thousand int64 `json:"thousand"`
}

type ProfileHead struct {
	Date        int64  // 交易日期  （精确到天:20220621)
	BatchNo     int64  // 批次号
	AccountId   string // 交易账户
	AccountName string // 户名
	AccountType int    // 用户类型 1-普通账户，2-汇总账户
	Provider    string // 算法厂商名称
	AlgoType    int    // 算法类型
	AlgoId      int    // 算法ID
	AlgoName    string // 算法名称
	SecId       string // 证券代码
	SecName     string // 证券名称
	Industry    string // 行业类型
	FundType    int    // 市值类型
	Liquidity   int    // 流动性
	AlgoOrderId int64  // 母单号
	TransAt     int64  // 交易时间 （精确到分钟）
	SourceFrom  int    // 来源  1-总线  2-导入
	SourcePrx   string // 源前缀
}

// OrderTradeQty 交易订单的一些数量结构体
type OrderTradeQty struct {
	EntrustQty int64 // 总委托数量
	DealQty    int64 // 总成交数量
	CancelQty  int64 // 总撤销数量
}

// Profile 计算各用户算法证券母单的盈亏信息
type Profile struct {
	ProfileHead
	IndexCount int // 计算次数

	//SecId          string  in head
	LastQty       int64 // 成交数量
	TotalTradeVol int64 // 总成交量,总交易成本  （sum(成交价格*成交数量))
	TotalCharge   int64 // 总手续费
	TotalCrossFee int64 // 流量费
	TotalBuyVol   int64 // 总买入金额
	TotalSellVol  int64 // 总卖出金额

	CancelRate     float64 // 撤单率
	MiniSplitOrder int64   // 最小拆单单位
	MiniDealOrder  int64   // 最小成交单位

	EntrustQty   int64   // 总委托数量
	DealQty      int64   // 总成交数量
	CancelQty    int64   // 总撤销数量
	Progress     float64 // 完成度
	MinJointRate float64 // 最小贴合度

	DealCount      int64   // 有成交回执个数 (一次成功回执算一次）
	TradeCount     int64   // 交易次数 (一个母单算一次交易）
	TradeCountPlus int64   // 正盈亏交易次数
	DealEffi       float64 // 成交效率
	AlgoOrderFit   float64 // 母单贴合度
	//PriceFit       float64 // 价格贴合度
	TradeVolFit float64 // 成交量贴合度
	PWP         float64 // pwp价格
	// 计算TWAP指标
	TwapStartTimePoint   int64   // 开始时间点 （时间戳，精确到秒)
	TwapDurTime          int64   // 经过时间  （第二笔成交和第一笔成交的时间差)
	TwapTotalTrade       int64   // 计算TWAP成交价格的分母值
	TwapTotalDur         int64   // 计算TWAP分子值    (该值也是母单的有效时长)
	TWAP                 float64 // TWAP成交值
	TwapTotalMarketTrade int64   // TWAP市场价格的分母值
	TWAPMarket           float64 // TWAP市场值
	TwapDev              float64 // TWAP滑点

	// T0 计算用
	BuyCost      int64 // 买入成本
	SellCost     int64 // 卖出成本
	AvgBuyPrice  int64 // 开仓价（买入成交均价）
	AvgSellPrice int64 // 平仓价 (卖出成交均价)
	TotalBuyQty  int64 // 买入总量
	TotalSellQty int64 // 卖出总量
	TotalT0Fee   int64 // T0 交易总手续费
	TotalT0Cost  int64 // T0 总交易成本   (买入+卖出)/2  即双边总交易额

	// 拆单计算用
	TotalEntrustCost int64   // 普通交易总额
	AvgEntrustPrice  float64 // 普通交易均价   (到达价均价)
	//TotalTradeCost   int64 // 成交总额
	AvgTradePrice float64 // 成交均价  （母单执行均价）
	TradeCost     int64   // 总交易成本
	TotalSplitFee int64   // 拆单总手续费

	ProfitAmount  float64 // 盈亏金额
	ProfitRate    float64 // 收益率
	WithdrawRate  float64 // 回撤比例
	MaxLossAmount float64 // 最大亏损金额
	AvgCost       float64 // 成交均价成本
	// for 滑点相关指标计算
	VwapDeal    int64   // sum（成交价格*成交数量）
	VwapEntrust int64   // sum(委托价格*成交数量)
	VWAP        float64 // vwap值
	VwapDev     float64 // vwap 滑点值
	VwapDevSum  float64 // 所有滑点值之和
	//VwapDevAvg   float64   // 滑点平均值
	VwapDevList []float64 // 滑点总列表（ 需要保存所有的滑点值，再算方差）
	//VwapVariance float64   // vwap 方差
	VwapStdDev float64 // vwap 滑点标准差
	// for 收益率标准差
	PfRateSum float64 // 收益率之和
	//PfRateAvg      float64   //  收益率平均值
	//PfRateVariance float64   // 收益率方差
	PfRateList   []float64 // 收益率总列表（需要保存所有的收益率，才能算方差)
	PfRateStdDev float64   // 收益率标准差
	// for 成交量贴合度标准差
	TradeVolFitList   []float64 // 母单贴合度列表
	TradeVolFitSum    float64   // 贴合度之和
	TradeVolFitStdDev float64   // 成交量贴合度标准差

	// for 时间贴合度标准差
	TimeFitList   []float64 // 贴合度列表
	TimeFitSum    float64   // 贴合度之和
	TimeFitStdDev float64   // 时间贴合度标准差

	AssessFactor float64 // 绩效收益率
	CreateTime   int64   // 创建时间
}

// ProfileSum 算法画像最终结构体
type ProfileSum struct {
	ProfileHead
	IndexCount int // 计算次数
	// 公共字段
	ProfitRate float64 // 收益率
	OrderNum   int64   // 订单数量

	TotalBuyVol  int64 // 总买入金额
	TotalSellVol int64 // 总卖出金额
	// 计算拆单盈亏金额
	ArriveCost      int64   // 到达价成本
	TotalLastCost   int64   // 总成交金额
	TotalArriveCost int64   // 普通交易总价
	AvgArriveCost   int64   // 普通交易均价
	AvgCost         float64 // 成交均价成本
	LastQty         int64   // 总成交数量
	TradeCountPlus  int64   // 盈利交易次数

	// for 经济性
	TradeVol       int64   // 交易量
	ProfitAmount   float64 // 盈亏金额
	TotalFee       int64   // 总手续费
	CrossFee       int64   // 流量费
	CancelRate     float64 // 撤单率
	MiniSplitOrder int64   // 最小拆单单位
	DealEffi       float64 // 成交效率
	// for 完成度
	EntrustQty   int64   // 总委托数量
	DealQty      int64   // 总成交数量
	CancelQty    int64   // 总撤销数量
	Progress     float64 // 完成度
	AlgoOrderFit float64 // 母单贴合度
	TradeVolFit  float64 // 成交量贴合度
	// for 风险度
	MinJointRate float64 // 最小贴合度
	WithdrawRate float64 // 回撤比例
	// for 算法绩效
	VwapDeal         int64     // sum（成交价格*成交数量）
	VwapEntrust      int64     // sum(委托价格*成交数量)
	VwapDev          float64   // vwap 滑点值
	VwapDevSum       float64   // 所有滑点值之和
	VwapDevCnt       int64     // 滑点个数 （交易次数）
	VwapDevAvg       float64   // 滑点平均值
	VwapDevList      []float64 // 滑点总列表（ 需要保存所有的滑点值，再算方差）
	VwapVariance     float64   // vwap 方差
	AssessProfitRate float64   // 绩效收益率
	// for 稳定性
	VwapStdDev        float64 // vwap 滑点标准差
	PfRateStdDev      float64 // 收益率标准差
	CV                float64 // 变异系数
	AssessFactor      float64 // 绩效因子
	TradeVolFitStdDev float64 // 成交量贴合度标准差
	TimeFitStdDev     float64 // 时间贴合度标准差
	// 评分
	EconomyScore   int // 经济性评分
	ProgressScore  int // 完成度评分
	RiskScore      int // 风险度评分
	AssessScore    int // 算法绩效评分
	StabilityScore int // 稳定性评分
	TotalScore     int // 综合评分

	FundPercent  FundRate       // 资金占比
	TradeVolDict TradeVolDirect // 买卖方向
	StockTypeVal StockType      // 股价类型
	TradeVolVal  TradeVolRate   // 交易量占比
}
