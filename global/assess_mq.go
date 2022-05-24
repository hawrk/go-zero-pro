// Package global
/*
 Author: hawrkchen
 Date: 2022/4/19 17:20
 Desc:
*/
package global

import (
	pb "algo_assess/assess-mq-server/proto/order"
	"sync"
	"time"
)

//----------ticker job
var GlobalOrders = MCalOrders{
	CalOrders: make(map[int64][]*pb.ChildOrderPerf),
}

type MCalOrders struct {
	CalOrders map[int64][]*pb.ChildOrderPerf // key -> time
	sync.RWMutex
}

//-------------real time calculate-------
var (
	GlobalAssess = MCalAssess{
		CalAlgo:  make(map[string]*OrderAssess),
		OrderMap: make(map[string]struct{}),
	}

	GlobalAlgoOrder = AlgoOrderMap{
		AlgoOrder:     make(map[string]int64),
		DealAlgoOrder: make(map[string]int64),
	}

	GlobalMarketLevel2 = MarketLevel2Map{
		EntrustVol: make(map[string]int64),
		TradeVol:   make(map[string]int64),
		LastPrice:  make(map[string]int64),
	}

	GlobalMarketVwap = MarketVwap{
		LastVol:     make(map[string]int64),
		TotalPrxCal: make(map[string]int64),
		TotalVol:    make(map[string]int64),
		MVwap:       make(map[string]float64),
	}
)

// 用于加快DB处理
var QuateKeyMap sync.Map

type MCalAssess struct {
	CalAlgo  map[string]*OrderAssess // key -> transat:algoId:secuId
	OrderMap map[string]struct{}     // key -> 子单ID， 用于判断同一个时间区内，同一支算法证券ID的委托数量
	sync.RWMutex
}

// 母单委托数量  数据源： 母单交易推送
type AlgoOrderMap struct {
	AlgoOrder     map[string]int64 // key-> 算法ID：证券ID  value -> 总委托数量
	DealAlgoOrder map[string]int64 // key -> 算法ID：证券ID   value -> 已成交数量
	sync.RWMutex
}

// 市场委托数量  数据源：十档行情快照
type MarketLevel2Map struct {
	EntrustVol map[string]int64 // key-> 证券ID:时间，  value -> 委托数量
	TradeVol   map[string]int64 // key -> 证券ID:时间 value -> 成交数量
	LastPrice  map[string]int64 // key-> 证券ID:时间   value -> 最新价格
	sync.RWMutex
}

// 计算市场vwap
type MarketVwap struct {
	LastVol     map[string]int64   // key-> 证券ID， value -> 上一次的成交总量，用于计算当前时间点的成交增量
	TotalPrxCal map[string]int64   // key -> 证券ID:时间   value -> 当前成交总量和最新价格的和， 即 ∑(订单成交数量 *成交价格)
	TotalVol    map[string]int64   // key -> 证券ID:时间   value -> 当前成交总量    即 ∑订单成交数量
	MVwap       map[string]float64 // key -> 证券ID:时间   value -> 市场vwap
	sync.RWMutex
}

// -------------------计算结果输出----------------------
// OrderAssess 绩效计算最终数据
type OrderAssess struct {
	AlgorithmType  uint
	AlgorithmId    uint
	UsecurityId    uint
	SecurityId     string
	TimeDimension  int
	TransactAt     int64   // 精确到分的时间 格式： 202204241532
	ArrivedPrice   int64   // 到达价格
	LastPrice      int64   // 最新价格
	SubVwap        int64   // 单笔vwap
	SubVwapEntrust int64   // 单笔委托vwap
	SubVwapArrived int64   // 到达价vwap
	Vwap           float64 // 市场vwap
	TradeVwap      float64 // vwap
	DealRate       float64 // 成交量比重
	OrderQty       int64   // 委托数量
	LastQty        int64   // 成交数量
	// TotalLastQty          uint64  // 累计成交
	CancelQty             int64   // 撤销数量
	RejectedQty           int64   // 拒绝数量
	MarketRate            float64 // 市场参与率
	DealProgress          float64 // 成交进度
	VwapDeviation         float64 // vwap 滑点
	ArrivedPriceDeviation float64 // 到达价滑点
	CreateTime            time.Time
}

// -------------------总线平台-> 绩效平台结构体转换--------
//母单下发
type MAlgoOrder struct {
	AlgoId       int
	AlgorithmId  int
	UsecId       int
	SecId        string
	AlgoOrderQty int64
	TransTime    int64
}

// 子单下发
type ChildOrderData struct {
	OrderId          int64  // 子单ID
	AlgoOrderId      int64  // 母单ID
	AlgorithmType    uint   // 算法类型
	AlgoId           uint   // 算法ID
	UsecId           uint   // 证券码
	SecId            string // 证券ID
	OrderQty         int64  // 委托数量 （实际数量）
	Price            int64  // 委托价格 （以分为单位）
	OrderType        uint   // 订单类型
	LastPx           int64  // 成交价格 （以分为单位）
	LastQty          int64  // 成交数量
	ComQty           int64  // 总成交数量 （实际数量）
	ArrivePrice      int64  // 到达价格
	ChildOrderStatus uint   // 订单交易状态
	TransTime        int64  // 交易时间 格式：202205050930
}
