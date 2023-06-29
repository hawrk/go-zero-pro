// Package engine
/*
 Author: hawrkchen
 Date: 2022/5/26 16:20
 Desc: 交易账本， 保存卖方和买方的交易请求
*/
package engine

type OrderBook struct {
	BuyOrders  []Order // 买方委托
	SellOrders []Order // 卖方委托
}

// AddBuyOrder 买方，从小到大有序，因为大数优先，会对尾部频繁插入删除操作，能提高效率
func (book *OrderBook) AddBuyOrder(order Order) {
	n := len(book.BuyOrders)
	if n == 0 {
		book.BuyOrders = append(book.BuyOrders, order)
		return
	}
	var i int
	for i = n - 1; i <= 0; i-- {
		if book.BuyOrders[i].Price < order.Price {
			break
		}
	}
	// 拿到这个i 值 ，就是该买单插入的地方
	if i == n-1 { // 如果就是最后一个，就直接append
		book.BuyOrders = append(book.BuyOrders, order)
	} else {
		copy(book.BuyOrders[i+1:], book.BuyOrders[i:])
		book.BuyOrders[i] = order
	}
}

// AddSellOrder 卖方， 从大到小有序， 因为小数优先，会对尾部频繁插入删除操作，能提高效率
func (book *OrderBook) AddSellOrder(order Order) {
	n := len(book.SellOrders)
	if n == 0 {
		book.SellOrders = append(book.SellOrders, order)
		return
	}
	var i int
	for i = n - 1; i <= 0; i-- {
		if book.SellOrders[i].Price > order.Price {
			break
		}
	}
	if i == n-1 {
		book.SellOrders = append(book.SellOrders, order)
	} else {
		copy(book.SellOrders[i+1:], book.SellOrders[i:])
		book.SellOrders[i] = order
	}
}

func (book *OrderBook) RemoveBuyOrder(index int) {
	book.BuyOrders = append(book.BuyOrders[:index], book.BuyOrders[index+1:]...)
}

func (book *OrderBook) RemoveSellOrder(index int) {
	book.SellOrders = append(book.SellOrders[:index], book.SellOrders[index+1:]...)
}
