// Package match_trade
/*
 Author: hawrkchen
 Date: 2022/5/26 17:19
 Desc:
*/
package main

import (
	"algo_assess/pkg/match-trade/engine"
	"fmt"
)

func main() {
	book := engine.OrderBook{
		BuyOrders:  make([]engine.Order, 0, 1000),
		SellOrders: make([]engine.Order, 0, 1000),
	}
	// 交易请求源可以从Kafka 来，
	order := engine.Order{
		Amount: 100,
		Price:  1550,
		Id:     "333",
		Side:   1,
	}

	trades := book.Process(order)
	fmt.Println("trades:", trades)
}
