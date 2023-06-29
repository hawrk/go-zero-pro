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
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type AssessJober interface {
	GetAssessWithCache() // 实时，读kafka 已计算完，直接取出
	DealDBAssess()       // 异常处理， 处理程序奔溃重新启动
	ClearLocalCache()    // 清除当天内存中的数据
	DealDBException()    // 异常处理， 定时处理二期落表失败的数据
	DoProfileWork()
	DoProfileSumWork()
	DoTimeLineWork()

	//
	KeepHeartBeat() // 心跳更新，用于主备机检测
	SyncHeartBeat()
}

type AssessJob struct {
	ctx context.Context
	s   *svc.ServiceContext
	logx.Logger
}

func NewOrderJob(c config.Config) AssessJober {
	return &AssessJob{
		ctx:    context.Background(),
		s:      svc.NewServiceContext(c),
		Logger: logx.WithContext(context.Background()),
	}
}

func (o *AssessJob) GetAssessWithCache() {
	// 1.读取缓存数据
	o.Logger.Info("order job begin.....")
	localMap := make(map[string]*global.OrderAssess)
	elapseMin := o.s.Config.WorkProcesser.ElapseMin
	now := cast.ToInt64(time.Now().Add(-time.Minute * time.Duration(elapseMin)).Format(global.TimeFormatMinInt))
	global.GlobalAssess.RWMutex.RLock()
	for key, val := range global.GlobalAssess.CalAlgo {
		o.Logger.Info("get now:", now, ",transtime:", val.TransactAt)
		if now <= val.TransactAt {
			continue
		}
		localMap[key] = val
	}
	global.GlobalAssess.RWMutex.RUnlock()
	o.Logger.Info("process localMap len:", len(localMap))
	if len(localMap) <= 0 {
		o.Logger.Info("no data, return")
		return
	}
	// 2.落地DB
	timeKey := make(map[int64]struct{})
	for key, val := range localMap {
		//o.Logger.Info("insert db key:", key)
		// 判断该key是否是fix key前缀，如果带fix,则是修复的数据，直接更新
		if strings.Contains(key, "Fix") {
			id, err := o.s.OrderAssessRepo.GetAlgoAssessByKey(context.Background(), val.UserId, val.TransactAt, val.AlgorithmId, val.SecurityId)
			if err != nil {
				o.Logger.Error("GetAlgoAssessByKey error:", err)
				return
			}
			if id <= 0 { // 无记录，则新增一条记录
				if err := o.s.OrderAssessRepo.CreateOrderAssess(context.Background(), val); err != nil {
					o.Logger.Error("create order assess fail:", err, "key:", key)
					// 插入失败时，一般为数据库出现异常了，这里直接退出本次操作，由下次启动Job时处理
					//continue
					return
				}
			} else { // 有记录，更新
				if err := o.s.OrderAssessRepo.UpdateAlgoAssess(context.Background(), val); err != nil {
					o.Logger.Error("update algo_assess error:", err)
					return
				}
			}

		} else {
			if err := o.s.OrderAssessRepo.CreateOrderAssess(context.Background(), val); err != nil {
				o.Logger.Error("create order assess fail:", err, "key:", key)
				// 插入失败时，一般为数据库出现异常了，这里直接退出本次操作，由下次启动Job时处理
				//continue
				return
			}
		}

		if _, exist := timeKey[val.TransactAt]; !exist {
			timeKey[val.TransactAt] = struct{}{}
		}
	}
	for key := range timeKey {
		// 3.删除redis
		o.s.RedisClient.Srem(global.AssessTimeSetKey, key)
		// 4.更新 Detail 表为已处理
		if err := o.s.OrderDetailRepo.UpdateOrderDetailStatus(key); err != nil {
			o.Logger.Error("update order detail fail:", err)
		}
	}
	// 5. 删除本地缓存
	global.GlobalAssess.RWMutex.Lock()
	for key := range localMap { // 只删除已落表的数据
		//o.Logger.Info("del cache key:", key)
		delete(global.GlobalAssess.CalAlgo, key)
	}
	global.GlobalAssess.RWMutex.Unlock()
	// 6. 删除 本地buffer
	for key := range localMap {
		delete(localMap, key)
	}
}
