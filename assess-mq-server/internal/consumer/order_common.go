// Package consumer
/*
 Author: hawrkchen
 Date: 2022/12/15 10:20
 Desc:
*/
package consumer

import (
	"algo_assess/assess-mq-server/internal/dao"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/global"
	"algo_assess/repo"
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

type AssessUser struct {
	NormalUser   string   // 普通用户
	MngrUser     []string // 普通用户所属管理员用户
	ProviderUser string   // 算法厂商用户
	AdminUser    string   // 超级管理员用户
}

// CheckAlgoOrderParam 母单入参校验
func CheckAlgoOrderParam(data *pb.AlgoOrderPerf) error {
	if data.TransactTime <= 0 {
		return errors.New("transtime invalid")
	}
	return nil
}

// CheckOrderParam 子单入参校验
func CheckOrderParam(v *pb.ChildOrderPerf) error {
	if v.GetAlgoOrderId() == 0 || v.AlgorithmId == 0 { // 过滤普通单，没有绑定算法的不参与计算
		return errors.New("normal order, continue")
	}
	if v.GetTransactTime() <= 0 {
		return errors.New("error field transTime")
	}
	if v.GetOrderQty() <= 0 {
		return errors.New("error field orderQty")
	}
	if v.GetChildOrdStatus() == global.OrderStatusApAccept || v.GetChildOrdStatus() == global.OrderStatusCtAccept ||
		v.GetChildOrdStatus() == global.OrderStatusTaAccept {
		return errors.New("accept status, continue")
	}
	//if childOrderVersion == 0 || int(v.GetVersion()) > childOrderVersion {  // 每个单的版本号独立
	//	childOrderVersion = int(v.GetVersion())
	//} else {
	//	return errors.New(fmt.Sprintf("invalid version, cur version:%d, request version:%d", childOrderVersion, v.GetVersion()))
	//}
	return nil
}

// CheckAlgoOrder 校验子单是否有母单信息
func CheckAlgoOrder(algoOrderKey string, rc *redis.Redis, data *pb.ChildOrderPerf) bool {
	//algoOrderId := fmt.Sprintf("algoOrderId:%d", data.AlgoOrderId)
	logx.Info("get algoOrderId key:", algoOrderKey)
	algoArrived := false
	for i := 0; i < 20; i++ {
		exist, err := rc.Exists(algoOrderKey)
		if err != nil {
			logx.Error("query redis algo order error:", err)
		}
		if !exist {
			logx.Info("algo order not found, waiting ......")
			time.Sleep(time.Millisecond * 100)
			continue
		}
		algoArrived = true
		break
	}
	if !algoArrived {
		//  找不到母单时，不再参与计算
		logx.Error("oops: algo order not found :", data.AlgoOrderId)
	}
	return algoArrived
}

// CheckQuote 校验行情信息
func CheckQuote(rc *redis.Redis, orderData *global.ChildOrderData) {
	QuoteKey := fmt.Sprintf("%s:%s:%d", global.Level2RedisPrx, orderData.SecId, orderData.TransTime)
	quoteArrived := false
	for i := 0; i < 20; i++ {
		exist, err := rc.Hexists(QuoteKey, "lastprice")
		if err != nil {
			logx.Error(" query redis quote data error:", err)
		}
		if !exist {
			logx.Info("quote data not found, waiting ......")
			time.Sleep(time.Millisecond * 200)
			continue
		}
		quoteArrived = true
		break
	}
	if !quoteArrived {
		logx.Error("oops: quote data not found:", QuoteKey)
	}

}

// GetArrviPrice 取行情价格
func GetArrviPrice(szRepo repo.MarketLevelRepo, shRepo repo.SHMarketLevelRepo, redis *redis.Redis, secId string, transTime int64) int64 {
	strTransTime := cast.ToString(transTime)
	quoteKey := fmt.Sprintf("%s:%d", secId, transTime)
	//logx.Info("get quoteKey:", quoteKey)
	var arriPrice int64
	// 1. 先取本地缓存数据
	global.GQuotes.RWMutex.Lock()
	if _, exist := global.GQuotes.Quotes[quoteKey]; exist {
		arriPrice = global.GQuotes.Quotes[quoteKey]
		//logx.Info("get local cache.")
	} else { // 本地缓存没有，取redis的
		// 当天的交易，如果行情推送有延迟，则当前的交易订单在Redis中无数据，并且在DB中也无该行情数据，
		// 最多倒推五分钟，取最近的一笔行情数据
		times, _ := time.ParseInLocation(global.TimeFormatMinInt, strTransTime, time.Local)
		for i := 0; i < 5; i++ {
			stMoveBackward := time.Unix(times.Unix(), 0).Add(-time.Minute * time.Duration(i)).Format(global.TimeFormatMinInt)
			redisQuoteKey := fmt.Sprintf("%s:%s:%s", global.Level2RedisPrx, secId, stMoveBackward)
			//logx.Info("get redis key:", redisQuoteKey)
			lastPrice, _ := redis.Hget(redisQuoteKey, "lastprice")
			if lastPrice != "" {
				arriPrice = cast.ToInt64(lastPrice)
				// 再同步到本地缓存中--缓存一条  (redis中在的，一般是当天数据，实时交易)
				global.GQuotes.Quotes[quoteKey] = arriPrice
				break
			}
		}

		if arriPrice == 0 { // 倒推五分钟后仍然为0，则从DB中取
			logx.Info("get db.")
			if secId[0:1] == "0" || secId[0:1] == "3" { //深市
				out, err := szRepo.GetSzMarketLevelBySecId(context.Background(), secId, strTransTime[:8])
				if err != nil {
					logx.Error("GetSzMarketLevelBySecId error:", err)
				}
				// 把DB查到的当天该证券的所有数据load到缓存中
				for _, v := range out {
					if v.OrgiTime == transTime {
						arriPrice = v.LastPrice // 找到当前这条，先赋值
					}
					bufkey := fmt.Sprintf("%s:%d", secId, v.OrgiTime)
					global.GQuotes.Quotes[bufkey] = v.LastPrice
				}
				if arriPrice == 0 { // 对15：00 后面的交易做下兼容，取最后一笔，即收盘价
					closingPrice, err := szRepo.GetSzClosingPrice(context.Background(), secId, strTransTime[:8])
					if err != nil {
						logx.Error("GetSzClosingPrice error:", err)
					}
					arriPrice = closingPrice
				}
			} else if secId[0:1] == "6" { // 沪市
				out, err := shRepo.GetShMarketLevelBySecId(context.Background(), secId, strTransTime[:8])
				if err != nil {
					logx.Error("GetShMarketLevelBySecId error:", err)
				}
				for _, v := range out {
					if v.OrgiTime == transTime {
						arriPrice = v.LastPrice
					}
					bufkey := fmt.Sprintf("%s:%d", secId, v.OrgiTime)
					global.GQuotes.Quotes[bufkey] = v.LastPrice
				}
				if arriPrice == 0 { // 对15：00 后面的交易做下兼容，取最后一笔，即收盘价
					closingPrice, err := shRepo.GetShClosingPrice(context.Background(), secId, strTransTime[:8])
					if err != nil {
						logx.Error("GetShClosingPrice error:", err)
					}
					arriPrice = closingPrice
				}
			}
		}
	}
	global.GQuotes.RWMutex.Unlock()
	return arriPrice * 100
}

// GetAllUsers 根据输入参数返回所属的所有用户名称
func GetAllUsers(userId string, algoId int) AssessUser {
	// 根据账户维度进行计算，分普通用户， 管理员用户，算法用户
	// 其中算法用户可以理解为一个虚拟用户，就是该算法下所有用户的统计，默认为 0
	// 如果是普通用户，计算该普通用户下的绩效，还需要找到其管理员，计算该管理员下的绩效，同时计算虚拟用户， 共3条记录
	// 如果是管理员用户， 计算该管理员下的绩效 和虚拟用户绩效 共2条记录
	// 如果就是算法用户，则只有一条记录
	var u AssessUser
	dao.GAccountMap.RWMutex.RLock()
	if _, exist := dao.GAccountMap.Account[userId]; exist {
		if dao.GAccountMap.Account[userId].UserType == 1 { // 普通账户，还需要找到其管理员账户
			u.NormalUser = userId
			u.MngrUser = dao.GAccountMap.Account[userId].ParUserId
		} else if dao.GAccountMap.Account[userId].UserType == 3 { // 管理员账户，直接取其user_id
			u.MngrUser = []string{userId}
		} else if dao.GAccountMap.Account[userId].UserType == 2 { // 算法厂商用户
			u.ProviderUser = userId
		}
	} else {
		logx.Info(" userId not found:", userId)
		// 本地缓存找不到账户信息时，只能计算所有用户的汇总绩效
	}
	dao.GAccountMap.RWMutex.RUnlock()

	if u.ProviderUser == "" { // 上面如果是普通用户||管理员的话，还要找到其所属算法厂商
		dao.GProviderMap.RWMutex.RLock()
		if _, exist := dao.GProviderMap.Provider[algoId]; exist {
			u.ProviderUser = dao.GProviderMap.Provider[algoId]
		} else {
			//providerUser = "0" // 无法找到算法厂商用户ID时，不再统计
			logx.Info("algo provider not found...")
		}
		dao.GProviderMap.RWMutex.RUnlock()
	}
	u.AdminUser = global.AdminUserId

	return u
}

// GetRedisOrderQty 取redis数据，返回 委托数量， 已成交数量， 撤销数量
// fixFlag 是否为数据修复的标志   0-正常交易 1- 数据修复标识
func GetRedisOrderQty(redis *redis.Redis, userId string, data *global.ChildOrderData) global.OrderTradeQty {
	//date := cast.ToString(data.TransTime)[:8]
	var userAlgoKey, userDealKey, userCancelKey string
	userAlgoKey = fmt.Sprintf("%s:%s:%d:%s:%d:%s:%d",
		data.SourcePrx, global.UserAlgoEntrust, data.CurDate, userId, data.AlgoId, data.SecId, data.AlgoOrderId) // 用户委托总量key(母单下发时已存）
	userDealKey = fmt.Sprintf("%s:%s:%d:%s:%d:%s:%d",
		data.SourcePrx, global.UserDealAlgo, data.CurDate, userId, data.AlgoId, data.SecId, data.AlgoOrderId) // 用户已成交数量key
	userCancelKey = fmt.Sprintf("%s:%s:%d:%s:%d:%s:%d",
		data.SourcePrx, global.UserCancelAlgo, data.CurDate, userId, data.AlgoId, data.SecId, data.AlgoOrderId) // 用户已撤销数量 key

	logx.Info("userAlgoKey:", userAlgoKey, ", userDealKey:", userDealKey, ", userCancelKey:", userCancelKey)
	// 取总委托数量
	entrustQty, err := redis.Get(userAlgoKey)
	if err != nil {
		logx.Error("get redis user algo order entrust qty err:", err)
		return global.OrderTradeQty{}
	}
	userEntrustQty := cast.ToInt64(entrustQty)
	// 累加当前成交量
	logx.Info("Get Redis OrderQty, incr deal qty:", data.LastQty)
	userDealQty, err := redis.Incrby(userDealKey, data.LastQty)
	if err != nil {
		logx.Error("incrby user deal order err:", err)
	}
	if err := redis.Expire(userDealKey, global.RedisKeyExpireTime); err != nil {
		logx.Error("expire algoDealKey :", userDealKey, " error:", err)
	}
	// 累加当前撤单量
	var userCancelQty int64
	var cancelerr error
	if data.ChildOrderStatus == global.OrderStatusCancel || data.ChildOrderStatus == global.OrderStatusApReject ||
		data.ChildOrderStatus == global.OrderStatusCtReject || data.ChildOrderStatus == global.OrderStatusTaReject {
		userCancelQty, cancelerr = redis.Incrby(userCancelKey, data.OrderQty-data.ComQty)
		if err != nil {
			logx.Error("incrby user cancel qty err:", cancelerr)
		}
	} else { // 正常交易的话，则取历史的撤单数量
		cancelQty, err := redis.Get(userCancelKey)
		if err != nil {
			logx.Error("get user cancel qty err:", cancelerr)
		}
		userCancelQty = cast.ToInt64(cancelQty)
	}
	q := global.OrderTradeQty{
		EntrustQty: userEntrustQty,
		DealQty:    userDealQty,
		CancelQty:  userCancelQty,
	}
	return q
}

func BuildProfileHead(userId string, orderData *global.ChildOrderData) global.ProfileHead {
	var head global.ProfileHead

	head.Date = orderData.CurDate
	head.TransAt = orderData.TransTime
	head.BatchNo = orderData.BatchNo
	head.SourceFrom = orderData.SourceFrom
	head.SourcePrx = orderData.SourcePrx
	// 填充算法信息
	head.AlgoId = orderData.AlgoId
	head.AlgoType = orderData.AlgorithmType
	head.SecId = orderData.SecId
	head.AlgoOrderId = orderData.AlgoOrderId

	dao.GAlgoBaseMap.RWMutex.RLock()
	head.AlgoName = dao.GAlgoBaseMap.AlgoBase[orderData.AlgoId].AlgoName
	head.Provider = dao.GAlgoBaseMap.AlgoBase[orderData.AlgoId].Provider
	dao.GAlgoBaseMap.RWMutex.RUnlock()
	// 填充证券信息
	dao.GSecurityMap.RWMutex.RLock()
	securityInfo := dao.GSecurityMap.SecurityBase[orderData.SecId]
	head.SecName = securityInfo.SecurityName
	head.Industry = securityInfo.Industry
	head.Liquidity = securityInfo.Liquidity
	head.FundType = securityInfo.FundType
	dao.GSecurityMap.RWMutex.RUnlock()

	// 填充账户信息
	head.AccountId = userId
	dao.GAccountMap.RWMutex.RLock()
	if _, exist := dao.GAccountMap.Account[userId]; exist {
		head.AccountName = dao.GAccountMap.Account[userId].UserName
		if dao.GAccountMap.Account[userId].UserType == 1 { // 普通用户
			head.AccountType = global.AccountTypeNormal
		} else if dao.GAccountMap.Account[userId].UserType == 2 { // 算法厂商用户
			head.AccountType = global.AccountTypeProvider
		} else if dao.GAccountMap.Account[userId].UserType == 3 { // 管理员
			head.AccountType = global.AccountTypeMngr
		}
	} else if userId == global.AdminUserId {
		head.AccountName = "管理员"
		head.AccountType = global.AccountTypeSuAdmin // 超级管理员
	} else {
		logx.Info("error: userId:", userId, " not found!!!!")
	}
	dao.GAccountMap.RWMutex.RUnlock()

	return head
}
