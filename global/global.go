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
	TimeFormatDay    = "20060102"
	AssessTimeSetKey = "assess-calculate"
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
