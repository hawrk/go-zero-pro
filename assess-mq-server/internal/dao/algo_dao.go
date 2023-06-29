// Package dao
/*
 Author: hawrkchen
 Date: 2022/7/6 15:24
 Desc:
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/svc"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

var GAlgoBaseMap = AlgoMap{
	AlgoBase: make(map[int]AlgoInfo),
}

type AlgoMap struct {
	AlgoBase map[int]AlgoInfo // key -> algoid
	sync.RWMutex
}

type AlgoInfo struct {
	AlgoName     string
	AlgoType     int
	AlgoTypeName string // T0, 拆单
	Provider     string // 算法厂商名称
}

func loadAlgoInfo(svcContext *svc.ServiceContext) error {
	// 切换到总线基础算法表
	/*
		rsp, err := svcContext.MornanoClient.GetAlgoInfo(context.Background(), &proto.AlgoInfoReq{OperType: 1})
		if err != nil {
			return err
		}
		infos := rsp.GetInfos()

		GAlgoBaseMap.RWMutex.Lock()
		defer GAlgoBaseMap.RWMutex.Unlock()

		for _, v := range infos {
			info := AlgoInfo{
				AlgoName:     v.AlgoName,
				AlgoType:     int(v.AlgoType),
				AlgoTypeName: v.AlgoTypeName,
				Provider:     v.Provider,
			}
			//logx.Infof("get info:%+v", info)
			GAlgoBaseMap.AlgoBase[int(v.AlgoId)] = info
		}
	*/
	//if len(GAlgoBaseMap.AlgoBase) <= 0 {
	//logx.Info("mornano-rpc-server algo info not found, reload local...")
	// 兼容，如果在总线DB中找不到算法基础信息，则从绩效本地加载
	assInfos, err := svcContext.AlgoBaseRepo.GetAlgoBase(context.Background())
	if err != nil {
		return err
	}
	for _, v := range assInfos {
		info := AlgoInfo{
			AlgoName:     v.AlgoName,
			AlgoType:     v.AlgoType,
			AlgoTypeName: v.AlgoTypeName,
			Provider:     v.Provider,
		}
		GAlgoBaseMap.AlgoBase[v.AlgoId] = info
	}
	if len(GAlgoBaseMap.AlgoBase) <= 0 {
		return errors.New(" query tb_algo_info no data")
	}
	//}

	logx.Info("load algo info success, len:", len(GAlgoBaseMap.AlgoBase))
	return nil
}
