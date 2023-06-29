// Package listen
/*
 Author: hawrkchen
 Date: 2022/3/24 10:09
 Desc:
*/
package listen

import (
	"algo_assess/assess-mq-server/internal/config"
	"algo_assess/assess-mq-server/internal/consumer"
	"algo_assess/assess-mq-server/internal/svc"
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func KafkaQueues(ctx context.Context, c config.Config, svcContext *svc.ServiceContext) []service.Service {
	// l := logx.WithContext(ctx)
	// l.Infof("order trade : %+v", c.AlgoPlatformOrderTradeConf)
	// l.Infof("market info %+v", c.AlgoPlatformMarketConf)
	return []service.Service{
		// 交易订单信息接收
		kq.MustNewQueue(c.AlgoPlatformOrderTradeConf, consumer.NewAlgoPlatformOrderTrade(ctx, svcContext)),
		// 母单信息接收
		kq.MustNewQueue(c.AlgoOrderTradeConf, consumer.NewAlgoOrderTrade(ctx, svcContext)),
		// 交易账户信息推送
		kq.MustNewQueue(c.AlgoAccountInfoConf, consumer.NewAccountInfo(ctx, svcContext)),

		// 子单修复
		//kq.MustNewQueue(c.PerfFixOrderTradeConf, consumer.NewPerfFixOrderTrade(ctx, svcContext)),
		// 母单修复
		//kq.MustNewQueue(c.PerfFixAlgoOrderTradeConf, consumer.NewPerfFixAlgoOrderTrade(ctx, svcContext)),
	}
}
