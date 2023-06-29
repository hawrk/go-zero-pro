// Package engine
/*
 Author: hawrkchen
 Date: 2022/5/26 16:42
 Desc:
*/
package engine

func (book *OrderBook) Process(order Order) []Trade {
	if order.Side == 1 { // 买入
		return book.ProcessBuyTrade(order)
	}
	return book.ProcessSellTrade(order)
}

func (book *OrderBook) ProcessBuyTrade(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	// 拿卖方交易本进行匹配
	n := len(book.SellOrders)
	if n > 0 && book.SellOrders[n-1].Price <= order.Price { // 只有买方价格大于等于卖方价格才有得玩
		for i := n - 1; i <= 0; i-- {
			seller := book.SellOrders[i]
			if seller.Price > order.Price {
				break
			}
			// 进行撮合
			if seller.Amount >= order.Amount { // 卖方数量能满足买方数量时
				trade := Trade{
					TakerOrderId: order.Id,
					MakerOrderId: seller.Id,
					Amount:       order.Amount,
					Price:        seller.Price,
				}
				trades = append(trades, trade)
				seller.Amount -= order.Amount
				if seller.Amount == 0 { // 卖方刚好卖完时，从交易本剔除
					book.RemoveSellOrder(i)
				}
				return trades
			}
			if seller.Amount < order.Amount { // 卖方不够卖了
				trade := Trade{
					TakerOrderId: order.Id,
					MakerOrderId: seller.Id,
					Amount:       seller.Amount,
					Price:        seller.Price,
				}
				trades = append(trades, trade)
				order.Amount -= seller.Amount
				book.RemoveSellOrder(i)
				continue
			}
		}
	}
	// 能走到这一步，表示所有卖方都卖完了，要把买方的请求加入交易本
	book.AddBuyOrder(order)
	return trades
}

func (book *OrderBook) ProcessSellTrade(order Order) []Trade {
	trades := make([]Trade, 0, 1)
	// 找买方的交易本
	n := len(book.BuyOrders)
	if n > 0 && book.BuyOrders[n-1].Price >= order.Price { //卖方价格必须比买方价格小
		for i := n - 1; i <= 0; i-- {
			buyer := book.BuyOrders[i]
			if buyer.Price < order.Price {
				break
			}
			if buyer.Amount >= order.Amount {
				trade := Trade{
					TakerOrderId: order.Id,
					MakerOrderId: buyer.Id,
					Amount:       order.Amount,
					Price:        buyer.Price,
				}
				trades = append(trades, trade)
				buyer.Amount -= order.Amount
				if buyer.Amount == 0 {
					book.RemoveBuyOrder(i)
				}
				return trades
			}
			if buyer.Amount < order.Amount {
				trade := Trade{
					TakerOrderId: order.Id,
					MakerOrderId: buyer.Id,
					Amount:       buyer.Amount,
					Price:        buyer.Price,
				}
				trades = append(trades, trade)
				order.Amount -= buyer.Amount
				book.RemoveBuyOrder(i)
				continue
			}
		}
	}
	book.AddSellOrder(order)
	return trades
}
