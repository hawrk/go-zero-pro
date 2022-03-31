// Package listen
/*
 Author: hawrkchen
 Date: 2022/3/24 10:09
 Desc:
*/
package listen

import (
	"algo_assess/mqueue/internal/config"
	"algo_assess/mqueue/internal/logic"
	"algo_assess/mqueue/internal/svc"
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

func KqMqs(ctx context.Context, c config.Config, svcContext *svc.ServiceContext) []service.Service {
	l := logx.WithContext(ctx)
	l.Infof("order trade : %+v", c.AlgoPlatformOrderTradeConf)
	// l.Infof("order result : %+v", c.AlgoPlatformOrderResultConf)
	l.Infof("market info %+v", c.AlgoPlatformMarketConf)
	return []service.Service{
		// 交易订单信息接收
		kq.MustNewQueue(c.AlgoPlatformOrderTradeConf, logic.NewAlgoPlatformOrderTrade(ctx, svcContext)),
		// 交易订单成交结果接收
		// kq.MustNewQueue(c.AlgoPlatformOrderResultConf, logic.NewAlgoPlatformOrderResult(ctx, svcContext)),
		// 行情信息接收
		kq.MustNewQueue(c.AlgoPlatformMarketConf, logic.NewAlgoPlatformMarketInfo(ctx, svcContext)),
	}
}
