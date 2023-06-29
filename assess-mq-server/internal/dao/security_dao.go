// Package dao
/*
 Author: hawrkchen
 Date: 2022/7/6 16:17
 Desc:
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/pkg/tools"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

var GSecurityMap = SecurityMap{
	SecurityBase: make(map[string]SecurityInfo),
}

type SecurityMap struct {
	SecurityBase map[string]SecurityInfo // key -> secrityId
	sync.RWMutex
}

type SecurityInfo struct {
	SecuritySource  string
	SecurityName    string
	PreClosePrice   float64
	Status          int
	UpperLimitPrice int64
	LowerLimitPrice int64
	FundType        int // 市值
	StockType       int
	Liquidity       int    // 流动性
	Industry        string // 行业类型
}

func loadSecurityInfo(svcContext *svc.ServiceContext) error {
	infos, err := svcContext.SecurityRepo.GetSecurityInfos(context.Background())
	if err != nil {
		return err
	}
	GSecurityMap.RWMutex.Lock()
	defer GSecurityMap.RWMutex.Unlock()

	for _, v := range infos {
		info := SecurityInfo{
			SecuritySource:  v.SecuritySource,
			SecurityName:    tools.RMu0000(v.SecurityName),
			PreClosePrice:   v.PreClosePx,
			Status:          v.Status,
			UpperLimitPrice: v.UpperLimitPrice,
			LowerLimitPrice: v.LowerLimitPrice,
			FundType:        v.FundType,
			StockType:       v.StockType,
			Liquidity:       v.Liquidity,
			Industry:        v.Industry,
		}
		GSecurityMap.SecurityBase[v.SecurityId] = info
	}
	if len(GSecurityMap.SecurityBase) <= 0 {
		logx.Info("security info load no data")
		return nil
	}
	logx.Info("load security info success, len:", len(GSecurityMap.SecurityBase))
	return nil
}
