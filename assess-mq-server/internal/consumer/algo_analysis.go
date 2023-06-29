// Package consumer
/*
 Author: hawrkchen
 Date: 2022/6/21 15:37
 Desc: 算法分析逻辑实现
*/
package consumer

import (
	"algo_assess/global"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/threading"
	"sync"
)

// AlgoAnalysisConcurrent 多用户计算并发版本
func (s *AlgoPlatformOrderTrade) AlgoAnalysisConcurrent(orderData *global.ChildOrderData, u AssessUser) {
	var wgUser sync.WaitGroup
	if u.NormalUser != "" {
		wgUser.Add(1)
		threading.GoSafe(func() {
			defer wgUser.Done()
			s.DisPatchAnalysis(u.NormalUser, orderData)
		})
	}

	for _, v := range u.MngrUser {
		wgUser.Add(1)
		go func(user string) {
			defer wgUser.Done()
			s.DisPatchAnalysis(user, orderData)
		}(v)
	}

	if u.ProviderUser != "" {
		wgUser.Add(1)
		threading.GoSafe(func() {
			defer wgUser.Done()
			s.DisPatchAnalysis(u.ProviderUser, orderData)
		})

	}
	//计算一个超级管理员的
	wgUser.Add(1)
	threading.GoSafe(func() {
		defer wgUser.Done()
		s.DisPatchAnalysis(u.AdminUser, orderData)
	})

	wgUser.Wait()
}

func (s *AlgoPlatformOrderTrade) AlgoAnalysis(orderData *global.ChildOrderData, u AssessUser) {
	if u.NormalUser != "" {
		s.DisPatchAnalysis(u.NormalUser, orderData)
	}

	for _, v := range u.MngrUser {
		s.DisPatchAnalysis(v, orderData)
	}

	if u.ProviderUser != "" {
		s.DisPatchAnalysis(u.ProviderUser, orderData)
	}
	//计算一个超级管理员的
	s.DisPatchAnalysis(u.AdminUser, orderData)
}

func (s *AlgoPlatformOrderTrade) DisPatchAnalysis(userId string, orderData *global.ChildOrderData) {
	s.Logger.Info("process user:[", userId, "]......................")
	sourcePrx := orderData.SourcePrx
	// profitKey 按天+证券ID区分
	ProfileKey := fmt.Sprintf("%s:%d:%s:%d:%s:%d", sourcePrx, orderData.CurDate, userId, orderData.AlgoId, orderData.SecId, orderData.AlgoOrderId)
	// profileKey 按天区分
	ProfileSumKey := fmt.Sprintf("%s:%d:%s:%d", sourcePrx, orderData.CurDate, userId, orderData.AlgoId)
	// timeLineKey 按分钟区分
	TimeLineKey := fmt.Sprintf("%s:%d:%s:%d", sourcePrx, orderData.TransTime, userId, orderData.AlgoId)

	s.Logger.Info("get  ProfileKey:", ProfileKey, ", ProfileSumKey:", ProfileSumKey, ", TimeLineKey:", TimeLineKey)
	// 先从redis 取交易订单的一些数量
	qtys := GetRedisOrderQty(s.svcCtx.RedisClient, userId, orderData)
	s.Logger.Infof("get qtys:%+v", qtys)
	// 填充头部基础数据
	head := BuildProfileHead(userId, orderData)

	var pv *global.Profile
	var sv *global.ProfileSum
	var tv *global.ProfileSum
	if head.AccountType == global.AccountTypeNormal {
		pv, sv, tv = s.NormalUserTask(ProfileKey, ProfileSumKey, TimeLineKey, &qtys, &head, orderData)
		s.LoadNorCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey, pv, sv, tv)
	} else if head.AccountType == global.AccountTypeProvider {
		pv, sv, tv = s.ProviderUserTask(ProfileKey, ProfileSumKey, TimeLineKey, &qtys, &head, orderData)
		s.LoadProviderCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey, pv, sv, tv)
	} else if head.AccountType == global.AccountTypeMngr {
		pv, sv, tv = s.MngrUserTask(ProfileKey, ProfileSumKey, TimeLineKey, &qtys, &head, orderData)
		s.LoadMngrCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey, pv, sv, tv)
	} else if head.AccountType == global.AccountTypeSuAdmin {
		pv, sv, tv = s.AdminUserTask(ProfileKey, ProfileSumKey, TimeLineKey, &qtys, &head, orderData)
		s.LoadAdminCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey, pv, sv, tv)
	}

	global.ProfileChan <- *pv
	global.ProfileSumChan <- *sv
	global.TlProfileSumChan <- *tv
}

// Write2FailProfileQueue DB异常处理，落表失败时，写入redis失败队列， 等待定时任务处理
func (s *AlgoPlatformOrderTrade) Write2FailProfileQueue(ProfileKey string, profile *global.Profile) {
	s.Logger.Info("Write2FailProfileQueue, Key:", ProfileKey)
	b, err := json.Marshal(profile)
	if err != nil {
		s.Logger.Error("Marshal profile error:", err)
	}
	if err := s.svcCtx.RedisClient.Hset(global.FailProfileKey, ProfileKey, string(b)); err != nil {
		s.Logger.Error("Write2FailProfileQueue, hset key fail:", err)
	}
}

func (s *AlgoPlatformOrderTrade) Write2FailTimeLineQueue(TimeLineKey string, timeline *global.ProfileSum) {
	s.Logger.Info("Write2FailTimeLineQueue, Key:", TimeLineKey)
	b, err := json.Marshal(timeline)
	if err != nil {
		s.Logger.Error("Marshal timeline  error:", err)
	}
	if err := s.svcCtx.RedisClient.Hset(global.FailTimeLineKey, TimeLineKey, string(b)); err != nil {
		s.Logger.Error("Write2FailTimeLineQueue, hset key fail:", err)
	}
}

func (s *AlgoPlatformOrderTrade) Write2FailProfitQueue(profileSumKey string, profileSum *global.ProfileSum) {
	s.Logger.Info("Write2FailProfitQueue, Key:", profileSumKey)
	b, err := json.Marshal(profileSum)
	if err != nil {
		s.Logger.Error("Marshal profit  error:", err)
	}
	if err := s.svcCtx.RedisClient.Hset(global.FailProfileSumKey, profileSumKey, string(b)); err != nil {
		s.Logger.Error("Write2FailProfitQueue, hset key fail:", err)
	}
}
