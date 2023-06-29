// Package global
/*
 Author: hawrkchen
 Date: 2022/3/23 10:55
 Desc:
*/
package global

const (
	RedisSep         = ":"
	TimeFormat       = "2006-01-02 15:04:05"
	TimeFormatMin    = "2006-01-02 15:04"
	TimeFormatMinInt = "200601021504"
	TimeFormatSecond = "20060102150405"
	TimeFormatDay    = "20060102"
	TimeFormatDaySp  = "2006.1.2"
	AssessTimeSetKey = "assess-calculate"

	HeartBeatKey = "heartbeat"

	MaxAssessTimeLen = 242      // 单日绩效汇总最大条数
	Level2RedisPrx   = "level2" // 行情十档数据Redis前缀
	//Level2FixRedisPrx = "level2fix:" // 行情十档修复数据Redis前缀  --行情共用一份数据

	AlgoOrderIdPrx = "algoOrderId"

	BasketsPrx    = "basket"
	BasketsFixPrx = "fix:basket"

	AlgoEntrustPrx = "algoEntrust" // 一期母单委托数量
	AlgoDealPrx    = "algoDeal"

	UserAlgoEntrust = "userAlgoEntrust" // 二期母单委托总数量
	UserDealAlgo    = "userDealAlgo"    // 二期已成交数量
	UserCancelAlgo  = "userCancelAlgo"  // 二期撤销数量

	OrderSourceNor = "Nor" // 正常订单前缀
	OrderSourceFix = "Fix" // 数据修复订单前缀
	OrderSourceOri = "Ori" // 原始订单前缀
	OrderSourceImp = "Imp" // 订单导入前缀
	OrderSourceAbt = "Abt" // 异常，未开放的来源类型

	TradeSideBuy  = 1 // 买
	TradeSideSell = 2 // 卖

	// 算法类型
	AlgoTypeT0        = 1 // T0算法
	AlgoTypeSplit     = 2 // 拆单算法  （智能委托）
	AlgoTypeNameT0    = "日内回转"
	AlgoTypeNameSplit = "智能委托"

	// 账户类型
	AccountTypeNormal   = 1 // 普通账户
	AccountTypeProvider = 2 // 算法厂商账户
	AccountTypeMngr     = 3 // 汇总账户
	AccountTypeSuAdmin  = 4 // 超级管理员

	// reload redis 加载的key
	CacheNorProfileKey    = "cacheNorProfile"
	CacheNorProfileSumKey = "cacheNorProfileSum"
	CacheNorTimeLineKey   = "cacheNorTimeline"

	CacheMngrProfileKey    = "cacheNorProfile"
	CacheMngrProfileSumKey = "cacheNorProfileSum"
	CacheMngrTimeLineKey   = "cacheNorTimeline"

	CacheProviderProfileKey    = "cacheNorProfile"
	CacheProviderProfileSumKey = "cacheNorProfileSum"
	CacheProviderTimeLineKey   = "cacheNorTimeline"

	CacheAdminProfileKey    = "cacheNorProfile"
	CacheAdminProfileSumKey = "cacheNorProfileSum"
	CacheAdminTimeLineKey   = "cacheNorTimeline"

	// 落表失败时，写redis失败队列的key
	FailProfileKey    = "FailProfileQueue"
	FailTimeLineKey   = "FailTimeLineQueue"
	FailProfileSumKey = "FailProfileSumQueue"

	// 存储在redis中的自定义Key过期时间
	RedisKeyExpireTime = 86400 // 86400秒， 1 day
	// 通用管理员账户
	AdminUserId = "admin"
	// 用户类型: 0-普通用户   1-管理员
	UserTypeNormal = 0
	UserTypeAdmin  = 1

	// 默认批次号
	DefaultBatchNo = 1

	// 数据来源
	SourceFromBus    = 0 // 总线
	SourceFromFix    = 1 // 数据修复
	SourceFromOrigin = 2 // 原始订单
	SourceFromImport = 3 // 订单导入

	MaxChannelBuffer = 1000000

	ConsisHasNode = "ConsisHashKey"
)

const (
	OrderStatusApAccept  = 0 // 总线接收
	OrderStatusApReject  = 1 // 总线拒绝
	OrderStatusCtAccept  = 2 // 柜台接收
	OrderStatusCtReject  = 3 // 柜台拒绝
	OrderStatusTaAccept  = 4 // 交易所接收
	OrderStatusTaReject  = 5 // 交易所拒绝
	OrderStatusPatiDeal  = 6 // 部分成交
	OrderStatusTotalDeal = 7 // 完全成交
	OrderStatusCancel    = 8 // 撤单
)

var ShOriginTime int64 = 202202280900
var SzOriginTime int64 = 202202280900
