// Package listen
/*
 Author: hawrkchen
 Date: 2022/3/24 10:09
 Desc:
*/
package listen

import (
	"algo_assess/mq-router-server/internal/config"
	"algo_assess/mq-router-server/internal/consumer"
	"algo_assess/mq-router-server/internal/svc"
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func KafkaQueues(ctx context.Context, c config.Config, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		// 交易订单信息接收
		kq.MustNewQueue(c.AlgoOrderTradeConf, consumer.NewAlgoOrderConsume(ctx, svcContext)),
		// 母单信息接收
		kq.MustNewQueue(c.ChildOrderTradeConf, consumer.NewChildOrderConsume(ctx, svcContext)),
	}
}
