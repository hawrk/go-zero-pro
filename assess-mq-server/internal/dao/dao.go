// Package dao
/*
 Author: hawrkchen
 Date: 2022/6/24 14:33
 Desc: 启动加载 基础数据
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/config"
	"algo_assess/assess-mq-server/internal/svc"
	"fmt"
)

type LoadTask struct {
	account *AccountDb
}

func NewLoadTask() *LoadTask {
	return &LoadTask{
		account: new(AccountDb),
	}
}

func (l *LoadTask) StartLoadTask(svcContext *svc.ServiceContext) {
	l.account.LoadAccountInfo(svcContext)
}

func Run(c config.Config, svcContext *svc.ServiceContext) error {
	// 加载 用户基础信息
	if err := loadAccountInfo(svcContext); err != nil {
		return err
	}
	// 加载 算法基础信息--加载绩效的
	if err := loadAlgoInfo(svcContext); err != nil {
		return err
	}
	// 加载证券基础信息
	if err := loadSecurityInfo(svcContext); err != nil {
		return err
	}
	// 加载算法厂商基础信息，这里需要从tb_algo_info算法基础表中找到所有算法厂商名称，然后到用户基础信息表中匹配对应的UserId,
	// 用作二期绩效根据算法厂商算画像评分
	if err := loadProviderInfo(svcContext); err != nil {
		return err
	}
	// 加载基础数据表
	if err := loadBusiConfig(svcContext); err != nil {
		return err
	}
	return nil
}

func Reload(c config.Config, svcContext *svc.ServiceContext) error {
	// 统一加载
	fmt.Println("reload cache data...")
	if c.WorkProcesser.EnableReloadCache {
		if err := reloadCacheData(svcContext); err != nil {
			return err
		}
	}
	return nil
}
