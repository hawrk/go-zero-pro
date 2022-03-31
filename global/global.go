// Package global
/*
 Author: hawrkchen
 Date: 2022/3/23 10:55
 Desc:
*/
package global

import (
	pb "algo_assess/mqueue/proto/order"
	"sync"
	"time"
)

const (
	RedisSep      = ":"
	TimeFormat    = "2006-01-02 15:04:05"
	TimeFormatMin = "2006-01-02 15:04"
	TimeFormatMinInt = "200601021504"
)

var GlobalOrders = MCalOrders{
	CalOrders: make(map[uint64][]pb.ChildOrderPerf),
}

type MCalOrders struct {
	CalOrders map[uint64][]pb.ChildOrderPerf // key -> time
	sync.RWMutex
}

type OrderAssess struct {
	AlgorithmType         uint
	AlgorithmId           uint
	UsecurityId           uint
	SecurityId            string
	TimeDimension         int
	TransactTime          int
	ArrivedPrice          int
	Vwap                  int
	DealRate              int
	OrderQty              uint
	LastQty               uint
	CancelQty             int
	RejectedQty           int
	MarketRate            int
	VwapDeviation         int
	ArrivedPriceDeviation int
	CreateTime            time.Time

}
