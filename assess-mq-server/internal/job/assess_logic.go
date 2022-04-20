// Package job
/*
 Author: hawrkchen
 Date: 2022/3/29 15:35
 Desc:
*/
package job

import (
	"algo_assess/assess-mq-server/internal/config"
	"algo_assess/assess-mq-server/internal/svc"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"runtime"
	"sync"
	"time"
)

type AssessJober interface {
	GetAssessWithCalcute() // 定时任务，需要计算
	GetAssessWithCache()   // 实时，读kafka 已计算完，直接取出
	CalculateAssess(transactAt int64, assess []*pb.ChildOrderPerf) error
	DealDBAssess()
}

type AssessJob struct {
	s *svc.ServiceContext
	logx.Logger
}

func NewOrderJob(c config.Config) AssessJober {
	return &AssessJob{
		s:      svc.NewServiceContext(c),
		Logger: logx.WithContext(context.Background()),
	}
}

func (o *AssessJob) GetAssessWithCalcute() {
	defer func() {
		if x := recover(); x != nil {
			o.Logger.Error("recover:", x)
			i := 0
			funcName, file, line, ok := runtime.Caller(i)

			for ok {
				o.Logger.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
				i++
				funcName, file, line, ok = runtime.Caller(i)
			}
		}
	}()
	// 开始计算
	// 1. 复制一份本地缓存
	localMap := make(map[int64][]*pb.ChildOrderPerf)
	// 计算当前时间，一分钟之内的请求暂不计算
	elapseMin := o.s.Config.WorkProcesser.ElapseMin
	now := cast.ToInt64(time.Now().Add(-time.Minute * time.Duration(elapseMin)).Format(global.TimeFormatMinInt))
	global.GlobalOrders.RWMutex.RLock()
	for k, v := range global.GlobalOrders.CalOrders {
		if now <= k {
			continue
		}
		localMap[k] = v
	}
	global.GlobalOrders.RWMutex.RUnlock()
	o.Logger.Info("in CalculateOrders job, get logmap:len:", len(localMap))

	for k, v := range localMap {
		o.Logger.Info("get key:", k, ", slice len:", len(v))
		mAssess := make(map[string][]*pb.ChildOrderPerf) // key -> algoid:securityid
		// 单个时间所有订单
		for _, detail := range v { // 转成根据算法ID+证券ID维度进行计算
			algoId := fmt.Sprintf("%d:%d", detail.GetAlgorithmId(), detail.GetUSecurityId())
			mAssess[algoId] = append(mAssess[algoId], detail)
		}
		var wg sync.WaitGroup
		var maxCh int
		if o.s.Config.WorkProcesser.GorontineNum > 1000 {
			maxCh = 1000
		} else {
			maxCh = o.s.Config.WorkProcesser.GorontineNum
		}
		workNum := make(chan struct{}, maxCh)
		for algoKey, values := range mAssess {
			wg.Add(1)
			workNum <- struct{}{}
			// 单个时间单支算法所有订单
			go func(transactAt int64, algoKey string, ass []*pb.ChildOrderPerf) {
				defer wg.Done()
				if err := o.CalculateAssess(k, ass); err != nil {
					o.Logger.Error("insert db fail:", err)
				}
				<-workNum
			}(k, algoKey, values)
		}
		wg.Wait()
		close(workNum)
		// 4.删除本地缓存的key
		global.GlobalOrders.RWMutex.Lock()
		delete(global.GlobalOrders.CalOrders, k)
		global.GlobalOrders.RWMutex.Unlock()
		// 5.redis 删除已处理的key
		o.s.RedisClient.Srem(global.AssessTimeSetKey, k)
		// 6.更新detail表为已处理
		if err := o.s.OrderDetailRepo.UpdateOrderDetail(k); err != nil {
			o.Logger.Error("update order detail fail:", err)
		}
		for key := range localMap {
			delete(localMap, key)
		}
	}
}

func (o *AssessJob) GetAssessWithCache() {
	defer func() {
		if x := recover(); x != nil {
			o.Logger.Error("recover:", x)
			i := 0
			funcName, file, line, ok := runtime.Caller(i)

			for ok {
				o.Logger.Errorf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
				i++
				funcName, file, line, ok = runtime.Caller(i)
			}
		}
	}()
	// 1.读取缓存数据
	localMap := make(map[string]*global.OrderAssess)
	elapseMin := o.s.Config.WorkProcesser.ElapseMin
	now := cast.ToInt64(time.Now().Add(-time.Minute * time.Duration(elapseMin)).Format(global.TimeFormatMinInt))
	global.GlobalAssess.RWMutex.RLock()
	for key, val := range global.GlobalAssess.CalAlgo {
		if now <= val.TransactAt {
			continue
		}
		localMap[key] = val
	}
	global.GlobalAssess.RWMutex.RUnlock()
	// 2.落地DB
	timeKey := make(map[int64]struct{})
	for key, val := range localMap {
		if err := o.s.OrderAssessRepo.CreateOrderAssess(context.Background(), val); err != nil {
			o.Logger.Error("create order assess fail:", err, "key:", key)
			continue
		}
		if _, exist := timeKey[val.TransactAt]; !exist {
			timeKey[val.TransactAt] = struct{}{}
		}
	}
	for key := range timeKey {
		// 3.删除redis
		o.s.RedisClient.Srem(global.AssessTimeSetKey, key)
		// 4.更新 Detail 表为已处理
		if err := o.s.OrderDetailRepo.UpdateOrderDetail(key); err != nil {
			o.Logger.Error("update order detail fail:", err)
		}
	}
	// 5. 删除 本地buffer
	for key := range localMap {
		delete(localMap, key)
	}
	// 6. 删除本地缓存
	global.GlobalAssess.RWMutex.Lock()
	for key := range global.GlobalAssess.CalAlgo {
		delete(global.GlobalAssess.CalAlgo, key)
	}
	global.GlobalAssess.RWMutex.Unlock()

}

func (o *AssessJob) CalculateAssess(transactAt int64, assess []*pb.ChildOrderPerf) error {
	var orderQty, lastQty, cancelQty, rejectQty uint64 // 委托数量， 成交数量， 撤销数量， 拒绝数量
	var subVWap, subVWapEntrust, subVWapArrived uint64
	var algoType, algoId, usecurityId uint32
	securityId := ""

	orderMap := make(map[uint32]struct{})
	for _, v := range assess {
		o.Logger.Info("get ID:", v.GetId())
		algoType = v.GetAlgorithmType()
		algoId = v.GetAlgorithmId()
		usecurityId = v.GetUSecurityId()
		securityId = v.GetSecurityId()
		subVWap += v.GetLastPx() * v.GetLastQty()       // 成交价格 * 成交笔数
		subVWapEntrust += v.GetPrice() * v.GetLastQty() // 委托价格 * 成交笔数
		subVWapArrived += 1 * v.GetLastQty()            // 到达价格 * 成交笔数   //TODO:
		lastQty += v.GetLastQty()
		// 委托数量需判断是否是同一个子订单
		if _, exist := orderMap[v.GetId()]; !exist {
			orderQty += v.GetOrderQty()
			orderMap[v.GetId()] = struct{}{}
		}
		if v.GetChildOrdStatus() == 8 {
			cancelQty += v.GetOrderQty() - v.GetCumQty()
		}
		if v.GetChildOrdStatus() == 1 || v.GetChildOrdStatus() == 3 || v.GetChildOrdStatus() == 5 {
			rejectQty += v.GetOrderQty()
		}
	}
	// check
	if lastQty <= 0 || orderQty <= 0 {
		return errors.New("divisor zero, cal assess fail!!!")
	}
	vWap := float64(subVWap) / float64(lastQty)
	vWapDeviation := vWap - float64(subVWapEntrust)/float64(lastQty)   // vwap 滑点
	arriveDeviation := vWap - float64(subVWapArrived)/float64(lastQty) // 到达价滑点
	dealRate := float64(lastQty) / float64(orderQty)                   // 成交进度
	marketRate := float64(orderQty) / 100000                           // 市场参与率    委托数量/行情数量    //TODO:
	// 3.落地assess表
	ass := &global.OrderAssess{
		AlgorithmType:         uint(algoType),
		AlgorithmId:           uint(algoId),
		UsecurityId:           uint(usecurityId),
		SecurityId:            securityId,
		TimeDimension:         2,
		TransactTime:          transactAt,
		ArrivedPrice:          0,
		Vwap:                  vWap,
		DealRate:              dealRate,
		OrderQty:              orderQty,
		LastQty:               lastQty,
		CancelQty:             cancelQty,
		RejectedQty:           rejectQty,
		MarketRate:            marketRate,
		VwapDeviation:         vWapDeviation,
		ArrivedPriceDeviation: arriveDeviation,
		CreateTime:            time.Now(),
	}
	if err := o.s.OrderAssessRepo.CreateOrderAssess(context.Background(), ass); err != nil {
		o.Logger.Error("create order assess fail:", err)
		return err
	}
	return nil
}
