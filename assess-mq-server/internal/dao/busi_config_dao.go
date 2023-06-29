// Package dao
/*
 Author: hawrkchen
 Date: 2022/7/15 10:24
 Desc:
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
)

func loadBusiConfig(svcContext *svc.ServiceContext) error {
	infos, err := svcContext.BusiConfigRepo.GetAllBusiConfig(context.Background())
	if err != nil {
		return err
	}
	// 初始化 画像基础数据
	m := make(map[int]string)
	for _, v := range infos {
		if v.BusiType == 4 { // 取画像的权重配置
			m[v.SecType] = v.Params
		}
	}
	GScoreConf = NewWeights(m)
	// 初始化资金市值基础数据
	logx.Info("loadBusiConfig success")
	return nil
}
