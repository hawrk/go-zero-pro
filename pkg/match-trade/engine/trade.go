// Package engine
/*
 Author: hawrkchen
 Date: 2022/5/26 16:17
 Desc: 交易信息结构体
*/
package engine

import "encoding/json"

type Trade struct {
	TakerOrderId string `json:"taker_order_id"`
	MakerOrderId string `json:"maker_order_id"`
	Amount       uint64 `json:"amount"`
	Price        uint64 `json:"price"`
}

func (t *Trade) FromJson(b []byte) error {
	return json.Unmarshal(b, t)
}

func (t *Trade) ToJson() []byte {
	out, _ := json.Marshal(t)
	return out
}
