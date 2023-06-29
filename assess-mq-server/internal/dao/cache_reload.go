// Package dao
/*
 Author: hawrkchen
 Date: 2022/10/27 14:48
 Desc: 重新加载DB的数据，防止服务panic后内存数据清空
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/global"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

func reloadCacheData(svcContext *svc.ServiceContext) error {
	// 从redis reload数据
	logx.Info("reload cache buffer....")
	reloadNorUser(svcContext)
	reloadMngrUser(svcContext)
	reloadProviderUser(svcContext)
	reloadAdminUser(svcContext)

	return nil

}

func reloadNorUser(svcContext *svc.ServiceContext) error {
	profile, err := svcContext.RedisClient.Hgetall(global.CacheNorProfileKey)
	if err != nil {
		logx.Error("reload Hgetall CacheProfileKey error:", err)
	}
	profit, err := svcContext.RedisClient.Hgetall(global.CacheNorProfileSumKey)
	if err != nil {
		logx.Error("reload Hgetall CacheProfitKey error:", err)
	}
	timeLine, err := svcContext.RedisClient.Hgetall(global.CacheNorTimeLineKey)
	if err != nil {
		logx.Error("reload Hgetall CacheTimeLineKey error:", err)
	}

	global.GNorUserProfile.RWMutex.Lock()
	for k, v := range profile {
		//logx.Info("get key:", k, "\n value:", v)
		var pf global.Profile
		json.Unmarshal([]byte(v), &pf)
		global.GNorUserProfile.ProfileMap[k] = &pf
	}
	for k, v := range profit {
		var ps global.ProfileSum
		json.Unmarshal([]byte(v), &ps)
		global.GNorUserProfile.ProfileSumMap[k] = &ps
	}
	for k, v := range timeLine {
		var tl global.ProfileSum
		json.Unmarshal([]byte(v), &tl)
		global.GNorUserProfile.TimeLineMap[k] = &tl
	}
	logx.Info("reload1 ProfileMap len:", len(global.GNorUserProfile.ProfileMap),
		",reload1 profileSumMap len:", len(global.GNorUserProfile.ProfileSumMap),
		",reload1 timeLineMap len:", len(global.GNorUserProfile.TimeLineMap))
	global.GNorUserProfile.RWMutex.Unlock()
	return nil
}

func reloadMngrUser(svcContext *svc.ServiceContext) error {
	profile, err := svcContext.RedisClient.Hgetall(global.CacheMngrProfileKey)
	if err != nil {
		logx.Error("reload Hgetall CacheMngrProfileKey error:", err)
	}
	profit, err := svcContext.RedisClient.Hgetall(global.CacheMngrProfileSumKey)
	if err != nil {
		logx.Error("reload Hgetall CacheMngrProfileSumKey error:", err)
	}
	timeLine, err := svcContext.RedisClient.Hgetall(global.CacheMngrTimeLineKey)
	if err != nil {
		logx.Error("reload Hgetall CacheMngrTimeLineKey error:", err)
	}

	global.GMngrProfile.RWMutex.Lock()
	for k, v := range profile {
		//logx.Info("get key:", k, "\n value:", v)
		var pf global.Profile
		json.Unmarshal([]byte(v), &pf)
		global.GMngrProfile.ProfileMap[k] = &pf
	}
	for k, v := range profit {
		var ps global.ProfileSum
		json.Unmarshal([]byte(v), &ps)
		global.GMngrProfile.ProfileSumMap[k] = &ps
	}
	for k, v := range timeLine {
		var tl global.ProfileSum
		json.Unmarshal([]byte(v), &tl)
		global.GMngrProfile.TimeLineMap[k] = &tl
	}
	logx.Info("reload2 ProfileMap len:", len(global.GMngrProfile.ProfileMap),
		",reload2 profileSumMap len:", len(global.GMngrProfile.ProfileSumMap),
		",reload2 timeLineMap len:", len(global.GMngrProfile.TimeLineMap))

	global.GMngrProfile.RWMutex.Unlock()
	return nil
}

func reloadProviderUser(svcContext *svc.ServiceContext) error {
	profile, err := svcContext.RedisClient.Hgetall(global.CacheProviderProfileKey)
	if err != nil {
		logx.Error("reload Hgetall CacheProviderProfileKey error:", err)
	}
	profit, err := svcContext.RedisClient.Hgetall(global.CacheProviderProfileSumKey)
	if err != nil {
		logx.Error("reload Hgetall CacheProviderProfileSumKey error:", err)
	}
	timeLine, err := svcContext.RedisClient.Hgetall(global.CacheProviderTimeLineKey)
	if err != nil {
		logx.Error("reload Hgetall CacheProviderTimeLineKey error:", err)
	}

	global.GProviderProfile.RWMutex.Lock()
	for k, v := range profile {
		//logx.Info("get key:", k, "\n value:", v)
		var pf global.Profile
		json.Unmarshal([]byte(v), &pf)
		global.GProviderProfile.ProfileMap[k] = &pf
	}
	for k, v := range profit {
		var ps global.ProfileSum
		json.Unmarshal([]byte(v), &ps)
		global.GProviderProfile.ProfileSumMap[k] = &ps
	}
	for k, v := range timeLine {
		var tl global.ProfileSum
		json.Unmarshal([]byte(v), &tl)
		global.GProviderProfile.TimeLineMap[k] = &tl
	}
	logx.Info("reload3 ProfileMap len:", len(global.GProviderProfile.ProfileMap),
		",reload3 profileSumMap len:", len(global.GProviderProfile.ProfileSumMap),
		",reload3 timeLineMap len:", len(global.GProviderProfile.TimeLineMap))
	global.GProviderProfile.RWMutex.Unlock()
	return nil
}

func reloadAdminUser(svcContext *svc.ServiceContext) error {
	profile, err := svcContext.RedisClient.Hgetall(global.CacheAdminProfileKey)
	if err != nil {
		logx.Error("reload Hgetall CacheProfileKey error:", err)
	}
	profit, err := svcContext.RedisClient.Hgetall(global.CacheAdminProfileSumKey)
	if err != nil {
		logx.Error("reload Hgetall CacheProfitKey error:", err)
	}
	timeLine, err := svcContext.RedisClient.Hgetall(global.CacheAdminTimeLineKey)
	if err != nil {
		logx.Error("reload Hgetall CacheTimeLineKey error:", err)
	}

	global.GAdminProfile.RWMutex.Lock()
	for k, v := range profile {
		//logx.Info("get key:", k, "\n value:", v)
		var pf global.Profile
		json.Unmarshal([]byte(v), &pf)
		global.GAdminProfile.ProfileMap[k] = &pf
	}
	for k, v := range profit {
		var ps global.ProfileSum
		json.Unmarshal([]byte(v), &ps)
		global.GAdminProfile.ProfileSumMap[k] = &ps
	}
	for k, v := range timeLine {
		var tl global.ProfileSum
		json.Unmarshal([]byte(v), &tl)
		global.GAdminProfile.TimeLineMap[k] = &tl
	}
	logx.Info("reload4 ProfileMap len:", len(global.GAdminProfile.ProfileMap),
		",reload4 profileSumMap len:", len(global.GAdminProfile.ProfileSumMap),
		",reload4 timeLineMap len:", len(global.GAdminProfile.TimeLineMap))
	global.GAdminProfile.RWMutex.Unlock()
	return nil
}
