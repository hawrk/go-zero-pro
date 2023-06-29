// Package dao
/*
 Author: hawrkchen
 Date: 2022/8/30 16:39
 Desc:
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/pkg/tools"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync"
)

var GProviderMap = ProviderMap{
	Provider: make(map[int]string),
}

type ProviderMap struct {
	Provider map[int]string // key -> 算法ID， value  -> 算法厂商用户ID, 二期绩效根据算法ID直接找到算法厂商的用户ID
	sync.RWMutex
}

func loadProviderInfo(svcContext *svc.ServiceContext) error {
	infos, _ := svcContext.AccountRepo.GetProviderAccountInfo(context.Background())
	GProviderMap.RWMutex.Lock()
	defer GProviderMap.RWMutex.Unlock()
	for _, v := range infos {
		GProviderMap.Provider[v.AlgoId] = strings.TrimSpace(tools.RMu0000(v.UserId))
	}
	if len(GProviderMap.Provider) <= 0 {
		return errors.New("algo info no data")
	}
	logx.Info("load account provider info success, len:", len(GProviderMap.Provider))
	return nil
}
