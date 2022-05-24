// Package listen
/*
 Author: hawrkchen
 Date: 2022/3/24 10:09
 Desc:
*/
package listen

import (
	"algo_assess/market-mq-server/internal/config"
	"algo_assess/market-mq-server/internal/consumer"
	"algo_assess/market-mq-server/internal/svc"
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

func KafkaQueues(ctx context.Context, c config.Config, svcContext *svc.ServiceContext) []service.Service {
	l := logx.WithContext(ctx)
	l.Infof("market info %+v", c.AlgoPlatformMarketConf)
	l.Infof("market sh info %+v", c.AlgoPlatFormSHMarketConf)
	return []service.Service{
		// 行情信息接收---深圳
		kq.MustNewQueue(c.AlgoPlatformMarketConf, consumer.NewAlgoPlatformMarketInfo(ctx, svcContext)),
		// 行情信息接收---上海
		kq.MustNewQueue(c.AlgoPlatFormSHMarketConf, consumer.NewAlgoPlatformSHMarketInfo(ctx, svcContext)),
	}
}
