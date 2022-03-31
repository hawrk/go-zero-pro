// Package job
/*
 Author: hawrkchen
 Date: 2022/3/29 15:35
 Desc:
*/
package job

import (
	"algo_assess/global"
	"algo_assess/mqueue/internal/config"
	"algo_assess/mqueue/internal/svc"
	pb "algo_assess/mqueue/proto/order"
	"context"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

type AssessJober interface {
	CalculateOrders()
	CalculateAssess(transactAt uint64,  assess []pb.ChildOrderPerf) error
}

type AssessJob struct {
	s *svc.ServiceContext
	logx.Logger
}

func NewOrderJob(c config.Config) AssessJober {
	return &AssessJob{
		s: svc.NewServiceContext(c),
		Logger: logx.WithContext(context.Background()),
	}
}

func (o *AssessJob) CalculateOrders() {
	// 开始计算
	// 1. 复制一份本地缓存
	localMap := make(map[uint64][]pb.ChildOrderPerf)
	global.GlobalOrders.RWMutex.RLock()
	for k, v := range global.GlobalOrders.CalOrders {
		localMap[k] = v
	}
	global.GlobalOrders.RWMutex.RUnlock()

	for k, v := range localMap {
		mAssess := make(map[string][]pb.ChildOrderPerf)   // key -> algoid:securityid
		// 单个时间所有订单
		for _, detail := range v {  // 转成根据算法ID+证券ID维度进行计算
			algoId := fmt.Sprintf("%d:%d",detail.GetAlgorithmId(), detail.GetUSecurityId())
			mAssess[algoId] = append(mAssess[algoId], detail)
		}
		// 开始逐个计算一支算法的绩效
		for _, values := range mAssess {
			// 单个时间单支算法所有订单
			o.CalculateAssess(k, values)
		}
		// 4.删除本地缓存的key
		global.GlobalOrders.RWMutex.RLock()
		delete(global.GlobalOrders.CalOrders, k)
		global.GlobalOrders.RWMutex.RUnlock()
		// 5. 更新detail表为已处理
		o.s.OrderDetail.UpdateOrderDetail(k)
	}

}

func (o *AssessJob) CalculateAssess(transactAt uint64, assess []pb.ChildOrderPerf) error {
	var orderQty, lastQty, cancelQty, rejectQty uint64       // 委托数量， 成交数量， 撤销数量， 拒绝数量
	var subVWap, subVWapEntrust, subVWapArrived uint64
	var algoType, algoId, usecurityId uint32
	securityId := ""

	orderMap := make(map[uint32]struct{})
	for _, v := range assess {
		algoType = v.GetAlgorithmType()
		algoId = v.GetAlgorithmId()
		usecurityId = v.GetUSecurityId()
		securityId = v.GetSecurityId()
		subVWap +=  v.GetLastPx() * v.GetLastQty()   // 成交价格 * 成交笔数
		subVWapEntrust +=  v.GetPrice() * v.GetLastQty()     // 委托价格 * 成交笔数
		subVWapArrived += 1 * v.GetLastQty()               // 到达价格 * 成交笔数   //TODO:
		lastQty += v.GetLastQty()
		// 委托数量需判断是否是同一个子订单
		if _, exist := orderMap[v.GetId()]; !exist {
			orderQty += v.GetOrderQty()
		}
		if v.GetChildOrdStatus() == 8 {
			cancelQty += v.GetOrderQty() - v.GetCumQty()
		}
		if v.GetChildOrdStatus() == 1|| v.GetChildOrdStatus() ==3 || v.GetChildOrdStatus() == 5 {
			rejectQty += v.GetOrderQty()
		}
	}
	vWap := cast.ToInt64(subVWap/lastQty)
	vWapDeviation := vWap - cast.ToInt64(subVWapEntrust/lastQty)   // vwap 滑点
	arriveDeviation := vWap - cast.ToInt64(subVWapArrived/lastQty)  // 到达价滑点
	dealRate := lastQty/orderQty       // 成交进度
	marketRate := orderQty/100000     // 市场参与率    委托数量/行情数量    //TODO:
	// 3.落地assess表
	ass := &global.OrderAssess{
		AlgorithmType:         cast.ToUint(algoType),
		AlgorithmId:           cast.ToUint(algoId),
		UsecurityId:           cast.ToUint(usecurityId),
		SecurityId:            securityId,
		TimeDimension:         2,
		TransactTime:          cast.ToInt(transactAt),   //TODO
		ArrivedPrice:          0,
		Vwap:                  cast.ToInt(vWap),
		DealRate:              cast.ToInt(dealRate),
		OrderQty:              cast.ToUint(orderQty),
		LastQty:               cast.ToUint(lastQty),
		CancelQty:             cast.ToInt(cancelQty),
		RejectedQty:           cast.ToInt(rejectQty),
		MarketRate:            cast.ToInt(marketRate),
		VwapDeviation:         cast.ToInt(vWapDeviation),
		ArrivedPriceDeviation: cast.ToInt(arriveDeviation),
		CreateTime:            time.Now(),
	}
	if err := o.s.OrderAssessRepo.CreateOrderAssess(context.Background(), ass); err != nil {
		o.Logger.Error("create order assess fail:", err)
		return err
	}
	return nil
}
