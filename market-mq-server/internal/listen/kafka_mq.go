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
	"github.com/zeromicro/go-zero/core/service"
)

func KafkaQueues(ctx context.Context, c config.Config, svcContext *svc.ServiceContext) []service.Service {
	//l := logx.WithContext(ctx)
	//l.Infof("market sz info %+v, market sh info %+v ", c.AlgoPlatformMarketConf, c.AlgoPlatFormSHMarketConf)
	//l.Infof("market sz fix info %+v, market sh fix info:%+v", c.PerfFixSZMarketConf, c.PerfFixSHMarketConf)
	return []service.Service{
		// 行情信息接收---深圳
		kq.MustNewQueue(c.AlgoPlatformMarketConf, consumer.NewAlgoPlatformMarketInfo(ctx, svcContext)),
		// 行情信息接收---上海
		kq.MustNewQueue(c.AlgoPlatFormSHMarketConf, consumer.NewAlgoPlatformSHMarketInfo(ctx, svcContext)),
		// 数据修复 --sz
		kq.MustNewQueue(c.PerfFixSZMarketConf, consumer.NewPerfFixSZMarketInfo(ctx, svcContext)),
		// 数据修复 --sh
		kq.MustNewQueue(c.PerfFixSHMarketConf, consumer.NewPerfFixSHMarketInfo(ctx, svcContext)),
	}
}
