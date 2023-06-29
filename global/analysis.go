// Package global
/*
 Author: hawrkchen
 Date: 2022/7/20 16:37
 Desc: 二期算法动态的一些基本结构体
*/
package global

import "sync"

var GAlgoOrderBasket = AlgoOrderQty{
	Baskets: make(map[string][]int),
}

type AlgoOrderQty struct {
	Baskets map[string][]int // key-> 用户ID:算法ID     value:  篮子ID列表，用来计算订单数（订单数以篮子为单位)
	sync.RWMutex
}

var GQuotes = AQuotes{
	Quotes: make(map[string]int64),
}

// AQuotes 保存所有A股证券的信息
type AQuotes struct {
	Quotes map[string]int64 // key-> 000001:202305050929, value -> lastprice
	sync.RWMutex
}
