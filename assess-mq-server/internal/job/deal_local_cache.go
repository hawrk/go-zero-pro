// Package job
/*
 Author: hawrkchen
 Date: 2022/10/27 10:27
 Desc:
*/
package job

import (
	"algo_assess/global"
	"time"
)

// ClearLocalCache  清除二期本地缓存的map数据
// 时间控制在晚上8点清空一次
// 目前已处理的结构体：   1. global.GAlgoProfile
// 2. redis中用于reload的数据
// 3. 清加载行情的数据 global.GQuotes
func (o *AssessJob) ClearLocalCache() {
	o.Logger.Info("ClearLocalCache begin....")
	t := time.Now()

	ClearNorUserCache()
	ClearMngrUserCache()
	ClearProviderUserCache()
	ClearAdminUserCache()

	// 清redis   -- redis 的key 过期时间统一在新增时设置
	/*
		n, err := o.s.RedisClient.Del(global.CacheProfileKey, global.CacheProfitKey, global.CacheTimeLineKey)
		if err != nil {
			o.Logger.Error("del key error:", err)
		}
		o.Logger.Info("del cache key num:", n)
	*/
	// 清行情数据
	global.GQuotes.RWMutex.Lock()
	for k := range global.GQuotes.Quotes {
		delete(global.GQuotes.Quotes, k)
		//o.Logger.Info("delete gquote key:", k)
	}
	global.GQuotes.RWMutex.Unlock()

	o.Logger.Info("finished clear local cache..,latency:", time.Since(t))
	//o.Logger.Info("ClearLocalCache end....")
}

func ClearNorUserCache() {
	global.GNorUserProfile.RWMutex.Lock()
	// 清profile
	for k := range global.GNorUserProfile.ProfileMap {
		delete(global.GNorUserProfile.ProfileMap, k)
		//o.Logger.Info("delete profile key:", k)
	}
	// 清ProfileSum
	for k := range global.GNorUserProfile.ProfileSumMap {
		delete(global.GNorUserProfile.ProfileSumMap, k)
		//o.Logger.Info("delete profile sum key:", k)
	}
	// 清time_line
	for k := range global.GNorUserProfile.TimeLineMap {
		delete(global.GNorUserProfile.TimeLineMap, k)
		//o.Logger.Info("delete time line key:", k)
	}

	global.GNorUserProfile.RWMutex.Unlock()
}

func ClearMngrUserCache() {
	global.GMngrProfile.RWMutex.Lock()
	// 清profile
	for k := range global.GMngrProfile.ProfileMap {
		delete(global.GMngrProfile.ProfileMap, k)
		//o.Logger.Info("delete profile key:", k)
	}
	// 清ProfileSum
	for k := range global.GMngrProfile.ProfileSumMap {
		delete(global.GMngrProfile.ProfileSumMap, k)
		//o.Logger.Info("delete profile sum key:", k)
	}
	// 清time_line
	for k := range global.GMngrProfile.TimeLineMap {
		delete(global.GMngrProfile.TimeLineMap, k)
		//o.Logger.Info("delete time line key:", k)
	}

	global.GMngrProfile.RWMutex.Unlock()
}

func ClearProviderUserCache() {
	global.GProviderProfile.RWMutex.Lock()
	// 清profile
	for k := range global.GProviderProfile.ProfileMap {
		delete(global.GProviderProfile.ProfileMap, k)
		//o.Logger.Info("delete profile key:", k)
	}
	// 清ProfileSum
	for k := range global.GProviderProfile.ProfileSumMap {
		delete(global.GProviderProfile.ProfileSumMap, k)
		//o.Logger.Info("delete profile sum key:", k)
	}
	// 清time_line
	for k := range global.GProviderProfile.TimeLineMap {
		delete(global.GProviderProfile.TimeLineMap, k)
		//o.Logger.Info("delete time line key:", k)
	}

	global.GProviderProfile.RWMutex.Unlock()
}

func ClearAdminUserCache() {
	global.GAdminProfile.RWMutex.Lock()
	// 清profile
	for k := range global.GAdminProfile.ProfileMap {
		delete(global.GAdminProfile.ProfileMap, k)
		//o.Logger.Info("delete profile key:", k)
	}
	// 清ProfileSum
	for k := range global.GAdminProfile.ProfileSumMap {
		delete(global.GAdminProfile.ProfileSumMap, k)
		//o.Logger.Info("delete profile sum key:", k)
	}
	// 清time_line
	for k := range global.GAdminProfile.TimeLineMap {
		delete(global.GAdminProfile.TimeLineMap, k)
		//o.Logger.Info("delete time line key:", k)
	}

	global.GAdminProfile.RWMutex.Unlock()
}
