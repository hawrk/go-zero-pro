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
		AlgoOrder: make(map[string]uint64),
	}

	GlobalMarketLevel2 = MarketLevel2Map{
		EntrustVol: make(map[string]uint64),
		TradeVol:   make(map[string]uint64),
		LastPrice:  make(map[string]int64),
	}
)

type MCalAssess struct {
	CalAlgo  map[string]*OrderAssess // key -> transat:algoId:secuId
	OrderMap map[string]struct{}     // key -> 交易时间：子订单号
	sync.RWMutex
}

// 母单委托数量  数据源： 母单交易推送
type AlgoOrderMap struct {
	AlgoOrder map[string]uint64 // key-> transat:algoId:secuId  value -> 总委托数量
	sync.RWMutex
}

type OrderAssess struct {
	AlgorithmType         uint
	AlgorithmId           uint
	UsecurityId           uint
	SecurityId            string
	TimeDimension         int
	TransactAt            int64 //精确到分的时间
	TransactTime          int64
	ArrivedPrice          int64   // 到达价格
	LastPrice             int64   // 最新价格
	SubVwap               uint64  // 单笔vwap
	SubVwapEntrust        uint64  // 单笔委托vwap
	SubVwapArrived        uint64  // 到达价vwap
	Vwap                  float64 // vwap
	DealRate              float64 // 成交量比重
	OrderQty              uint64  // 委托数量
	LastQty               uint64  // 成交数量
	TotalLastQty          uint64  // 累计成交
	CancelQty             uint64  // 撤销数量
	RejectedQty           uint64  // 拒绝数量
	MarketRate            float64 // 市场参与率
	DealProgress          float64 // 成交进度
	VwapDeviation         float64 // vwap 滑点
	ArrivedPriceDeviation float64 // 到达价滑点
	CreateTime            time.Time
}

// 市场委托数量  数据源：十档行情快照
type MarketLevel2Map struct {
	EntrustVol map[string]uint64 // key-> 证券ID:时间，  value -> 委托数量
	TradeVol   map[string]uint64 // key -> 证券ID:时间 value -> 成交数量
	LastPrice  map[string]int64  // key-> 证券ID:时间   value -> 最新价格
	sync.RWMutex
}

var QuoteMap = make(map[string]uint64)
