// Package consumer
/*
 Author: hawrkchen
 Date: 2023/5/18 16:49
 Desc:
*/
package consumer

import (
	"algo_assess/global"
	"encoding/json"
)

func (s *AlgoPlatformOrderTrade) NormalUserTask(pfk, pfsk, tlk string, qtys *global.OrderTradeQty,
	head *global.ProfileHead, order *global.ChildOrderData) (*global.Profile, *global.ProfileSum, *global.ProfileSum) {

	global.GNorUserProfile.RWMutex.Lock()
	// 计算画像明细
	pv := global.GNorUserProfile.ProfileMap[pfk]
	if pv == nil {
		s.Logger.Info("new User ProfileKey record:", pfk)
		pv = new(global.Profile)
		pv.CreateTime = order.TransTime
		pv.MiniSplitOrder = order.OrderQty
		pv.MiniDealOrder = order.LastQty
	}
	s.Logger.Infof("before cal User Profile, pv:%+v", pv)
	CalculateProfile(head, pv, order, qtys, s.svcCtx.RedisClient)
	s.Logger.Infof("after cal User Profile, pv:%+v", pv)
	global.GNorUserProfile.ProfileMap[pfk] = pv

	// 计算画像汇总
	sv := global.GNorUserProfile.ProfileSumMap[pfsk] // 取当前key 的信息
	if sv == nil {
		s.Logger.Info("new User ProfileSum record, key:", pfsk)
		sv = new(global.ProfileSum)
	}
	s.Logger.Infof(" before  cal User, get profile sum:%+v", sv)
	CalculateProfileSum(head, pfsk, sv, order, global.AccountTypeNormal)
	// 计算算法动态
	CalculateDynamic(sv, order)
	s.Logger.Infof("after cal User, profile sum:%+v", sv)
	global.GNorUserProfile.ProfileSumMap[pfsk] = sv

	// 计算时间线展示
	tv := global.GNorUserProfile.TimeLineMap[tlk]
	if tv == nil {
		s.Logger.Info("new time line record:", tlk)
		tv = new(global.ProfileSum)
	}
	// 计算时间线
	CalculateTimeLine(head, pfsk, tv, order, global.AccountTypeNormal)
	s.Logger.Infof("after timeline cal User, timeline:%+v", tv)
	global.GNorUserProfile.TimeLineMap[tlk] = tv

	global.GNorUserProfile.RWMutex.Unlock()

	return pv, sv, tv
}

func (s *AlgoPlatformOrderTrade) ProviderUserTask(pfk, pfsk, tlk string, qtys *global.OrderTradeQty,
	head *global.ProfileHead, order *global.ChildOrderData) (*global.Profile, *global.ProfileSum, *global.ProfileSum) {

	global.GProviderProfile.RWMutex.Lock()
	// 计算画像明细
	pv := global.GProviderProfile.ProfileMap[pfk]
	if pv == nil {
		s.Logger.Info("new Provider ProfileKey record:", pfk)
		pv = new(global.Profile)
		pv.CreateTime = order.TransTime
		pv.MiniSplitOrder = order.OrderQty
		pv.MiniDealOrder = order.LastQty
	}
	s.Logger.Infof("before cal Provider Profile, pv:%+v", pv)
	CalculateProfile(head, pv, order, qtys, s.svcCtx.RedisClient)
	s.Logger.Infof("after cal Provider Profile, pv:%+v", pv)
	global.GProviderProfile.ProfileMap[pfk] = pv

	// 计算画像汇总
	sv := global.GProviderProfile.ProfileSumMap[pfsk] // 取当前key 的信息
	if sv == nil {
		s.Logger.Info("new ProfileSum record, key:", pfsk)
		sv = new(global.ProfileSum)
	}
	s.Logger.Infof(" before cal Provider, get profile sum:%+v", sv)
	CalculateProfileSum(head, pfsk, sv, order, global.AccountTypeProvider)
	// 计算算法动态
	CalculateDynamic(sv, order)
	s.Logger.Infof("after cal Provider, profile sum:%+v", sv)
	global.GProviderProfile.ProfileSumMap[pfsk] = sv

	// 计算时间线展示
	tv := global.GProviderProfile.TimeLineMap[tlk]
	if tv == nil {
		s.Logger.Info("new time line record:", tlk)
		tv = new(global.ProfileSum)
	}
	// 计算时间线
	CalculateTimeLine(head, pfsk, tv, order, global.AccountTypeProvider)
	s.Logger.Infof("after timeline cal Provider, timeline:%+v", tv)
	global.GProviderProfile.TimeLineMap[tlk] = tv

	global.GProviderProfile.RWMutex.Unlock()

	return pv, sv, tv
}

func (s *AlgoPlatformOrderTrade) MngrUserTask(pfk, pfsk, tlk string, qtys *global.OrderTradeQty,
	head *global.ProfileHead, order *global.ChildOrderData) (*global.Profile, *global.ProfileSum, *global.ProfileSum) {

	global.GMngrProfile.RWMutex.Lock()
	// 计算画像明细
	pv := global.GMngrProfile.ProfileMap[pfk]
	if pv == nil {
		s.Logger.Info("new Mngr security ProfileKey record:", pfk)
		pv = new(global.Profile)
		pv.CreateTime = order.TransTime
		pv.MiniSplitOrder = order.OrderQty
		pv.MiniDealOrder = order.LastQty
	}
	s.Logger.Infof("before cal Mngr Profile, pv:%+v", pv)
	CalculateProfile(head, pv, order, qtys, s.svcCtx.RedisClient)
	s.Logger.Infof("after cal Mngr Profile, pv:%+v", pv)
	global.GMngrProfile.ProfileMap[pfk] = pv

	// 计算画像汇总
	sv := global.GMngrProfile.ProfileSumMap[pfsk] // 取当前key 的信息
	if sv == nil {
		s.Logger.Info("new ProfileSum record, key:", pfsk)
		sv = new(global.ProfileSum)
	}
	s.Logger.Infof(" before cal Mngr, get profile sum:%+v", sv)
	CalculateProfileSum(head, pfsk, sv, order, global.AccountTypeMngr)
	// 计算算法动态
	CalculateDynamic(sv, order)
	s.Logger.Infof("after cal Mngr, profile sum:%+v", sv)
	global.GMngrProfile.ProfileSumMap[pfsk] = sv

	// 计算时间线展示
	tv := global.GMngrProfile.TimeLineMap[tlk]
	if tv == nil {
		s.Logger.Info("new time line record:", tlk)
		tv = new(global.ProfileSum)
	}
	// 计算时间线
	CalculateTimeLine(head, pfsk, tv, order, global.AccountTypeMngr)
	s.Logger.Infof("after timeline cal Mngr, timeline:%+v", tv)
	global.GMngrProfile.TimeLineMap[tlk] = tv

	global.GMngrProfile.RWMutex.Unlock()

	return pv, sv, tv
}

func (s *AlgoPlatformOrderTrade) AdminUserTask(pfk, pfsk, tlk string, qtys *global.OrderTradeQty,
	head *global.ProfileHead, order *global.ChildOrderData) (*global.Profile, *global.ProfileSum, *global.ProfileSum) {

	global.GAdminProfile.RWMutex.Lock()
	// 计算画像明细
	pv := global.GAdminProfile.ProfileMap[pfk]
	if pv == nil {
		s.Logger.Info("new security ProfileKey record:", pfk)
		pv = new(global.Profile)
		pv.CreateTime = order.TransTime
		pv.MiniSplitOrder = order.OrderQty
		pv.MiniDealOrder = order.LastQty
	}
	s.Logger.Infof("before cal Admin Profile, pv:%+v", pv)
	CalculateProfile(head, pv, order, qtys, s.svcCtx.RedisClient)
	s.Logger.Infof("after cal Admin Profile, pv:%+v", pv)
	global.GAdminProfile.ProfileMap[pfk] = pv

	// 计算画像汇总
	sv := global.GAdminProfile.ProfileSumMap[pfsk] // 取当前key 的信息
	if sv == nil {
		s.Logger.Info("new ProfileSum record, key:", pfsk)
		sv = new(global.ProfileSum)
	}
	s.Logger.Infof(" before cal Admin, get profile sum:%+v", sv)
	CalculateProfileSum(head, pfsk, sv, order, global.AccountTypeSuAdmin)
	// 计算算法动态
	CalculateDynamic(sv, order)
	s.Logger.Infof("after cal Admin, profile sum:%+v", sv)
	global.GAdminProfile.ProfileSumMap[pfsk] = sv

	// 计算时间线展示
	tv := global.GAdminProfile.TimeLineMap[tlk]
	if tv == nil {
		s.Logger.Info("new time line record:", tlk)
		tv = new(global.ProfileSum)
	}
	// 计算时间线
	CalculateTimeLine(head, pfsk, tv, order, global.AccountTypeSuAdmin)
	s.Logger.Infof("after timeline cal Admin, timeline:%+v", tv)
	global.GAdminProfile.TimeLineMap[tlk] = tv

	global.GAdminProfile.RWMutex.Unlock()

	return pv, sv, tv
}

// LoadNorCacheBuffer 普通用户
func (s *AlgoPlatformOrderTrade) LoadNorCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey string,
	profile *global.Profile, profileSum *global.ProfileSum, timeline *global.ProfileSum) {
	// 落Redis
	b, err := json.Marshal(profile)
	if err != nil {
		s.Logger.Error("Marshal profile error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheNorProfileKey, ProfileKey, string(b))

	b2, err := json.Marshal(profileSum)
	if err != nil {
		s.Logger.Error("Marshal profit error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheNorProfileSumKey, ProfileSumKey, string(b2))

	b3, err := json.Marshal(timeline)
	if err != nil {
		s.Logger.Error("Marshal time line error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheNorTimeLineKey, TimeLineKey, string(b3))

	s.svcCtx.RedisClient.Expire(global.CacheNorProfileKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheNorProfileSumKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheNorTimeLineKey, global.RedisKeyExpireTime)
}

// LoadMngrCacheBuffer 管理员
func (s *AlgoPlatformOrderTrade) LoadMngrCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey string,
	profile *global.Profile, profileSum *global.ProfileSum, timeline *global.ProfileSum) {
	// 落Redis
	b, err := json.Marshal(profile)
	if err != nil {
		s.Logger.Error("Marshal profile error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheMngrProfileKey, ProfileKey, string(b))

	b2, err := json.Marshal(profileSum)
	if err != nil {
		s.Logger.Error("Marshal profit error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheMngrProfileSumKey, ProfileSumKey, string(b2))

	b3, err := json.Marshal(timeline)
	if err != nil {
		s.Logger.Error("Marshal time line error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheMngrTimeLineKey, TimeLineKey, string(b3))

	s.svcCtx.RedisClient.Expire(global.CacheMngrProfileKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheMngrProfileSumKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheMngrTimeLineKey, global.RedisKeyExpireTime)
}

// LoadProviderCacheBuffer 算法厂商
func (s *AlgoPlatformOrderTrade) LoadProviderCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey string,
	profile *global.Profile, profileSum *global.ProfileSum, timeline *global.ProfileSum) {
	// 落Redis
	b, err := json.Marshal(profile)
	if err != nil {
		s.Logger.Error("Marshal profile error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheProviderProfileKey, ProfileKey, string(b))

	b2, err := json.Marshal(profileSum)
	if err != nil {
		s.Logger.Error("Marshal profit error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheProviderProfileSumKey, ProfileSumKey, string(b2))

	b3, err := json.Marshal(timeline)
	if err != nil {
		s.Logger.Error("Marshal time line error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheProviderTimeLineKey, TimeLineKey, string(b3))

	s.svcCtx.RedisClient.Expire(global.CacheProviderProfileKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheProviderProfileSumKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheProviderTimeLineKey, global.RedisKeyExpireTime)
}

// LoadAdminCacheBuffer  超管
func (s *AlgoPlatformOrderTrade) LoadAdminCacheBuffer(ProfileKey, ProfileSumKey, TimeLineKey string,
	profile *global.Profile, profileSum *global.ProfileSum, timeline *global.ProfileSum) {
	// 落Redis
	b, err := json.Marshal(profile)
	if err != nil {
		s.Logger.Error("Marshal profile error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheAdminProfileKey, ProfileKey, string(b))

	b2, err := json.Marshal(profileSum)
	if err != nil {
		s.Logger.Error("Marshal profit error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheAdminProfileSumKey, ProfileSumKey, string(b2))

	b3, err := json.Marshal(timeline)
	if err != nil {
		s.Logger.Error("Marshal time line error:", err)
	}
	s.svcCtx.RedisClient.Hset(global.CacheAdminTimeLineKey, TimeLineKey, string(b3))

	s.svcCtx.RedisClient.Expire(global.CacheAdminProfileKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheAdminProfileSumKey, global.RedisKeyExpireTime)
	s.svcCtx.RedisClient.Expire(global.CacheAdminTimeLineKey, global.RedisKeyExpireTime)
}
