// Package engine
/*
 Author: hawrkchen
 Date: 2022/5/26 16:11
 Desc: 订单结构体，包括买入卖出
*/
package engine

import "encoding/json"

type Order struct {
	Amount uint64 `json:"amount"` // 数量
	Price  uint64 `json:"price"`  // 价格
	Id     string `json:"id"`     // 序列号
	Side   int8   `json:"side"`   // 买卖方向  1-买入   -1 -卖出
}

func (o *Order) FromJson(b []byte) error {
	return json.Unmarshal(b, o)
}

func (o *Order) ToJson() []byte {
	out, _ := json.Marshal(o)
	return out
}
