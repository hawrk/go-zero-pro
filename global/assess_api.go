// Package global
/*
 Author: hawrkchen
 Date: 2022/5/27 15:43
 Desc:
*/
package global

import "errors"

const (
	JwtUserId    = "jwtUserId"  // Token信息中保存用户ID的Key
	UserTokenKey = "user_token" // 保存redis的token前缀

	TokenKey    = "tb_session"
	TokenLen    = 32
	ServerIPLen = 32
	// LoginCount redis 保存登陆次数的Key
	LoginCount = "loginCnt" // 登陆次数:用户ID
)

var TokenErr = errors.New("token check fail")

// 篮子ID， 母单ID，算法ID，算法类型，用户ID，证券代码，母单数量，交易时间
// 开始时间， 结束时间
var AlgoHeader = []string{"basket_id(篮子ID)", "algo_order_id(母单ID)", "algorithm_id(算法ID)",
	"algorithm_type(算法类型)", "user_id(用户ID)", "security_id(证券代码)", "algo_order_qty(母单数量)", "transact_time(交易时间)",
	"start_time(开始时间)", "end_time(结束时间)"}

// 订单ID，母单ID，算法ID，算法类型，用户ID，证券代码，买卖方向，委托数量
// 委托价格， 订单类型，成交价格，成交数量，累计成交数量，手续费，到达价格，订单状态，交易时间
var ChildOrderHeader = []string{"order_id(订单ID)", "algo_order_id(母单ID)", "algorithm_id(算法ID)", "algorithm_type(算法类型)",
	"user_id(用户ID)", "security_id(证券代码)", "side(买卖方向)", "order_qty(委托数量)",
	"order_price(委托价格)", "order_type(订单类型)", "last_price(成交价格)", "last_qty(成交数量)", "cum_qty(累计成交数量)",
	"charge(手续费)", "arrived_price(到达价格)", "order_status(订单状态)", "transact_time(交易时间)"}

// token 解析二进制结构体
type Session struct {
	Id            uint32         // ID
	UuserId       uint32         // 用户ID
	Token         [TokenLen]byte // 会话Token
	ClientType    int8           // 客户端类型
	SessionStatus uint8          // 登陆状态 0:未登录 1正常登录状态 2 网络异常断开 3 会话超时 4 正常登出"
	Socket        int32          // 套接字，内部使用
	ServerId      uint16         // int64 -> uint16
	CreateTime    uint64
	UpdateTime    uint64
	ServerIp      [ServerIPLen]byte // 服务地址端口
}
