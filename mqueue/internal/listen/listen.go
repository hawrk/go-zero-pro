// Package listen
/*
 Author: hawrkchen
 Date: 2022/3/24 10:08
 Desc:
*/
package listen

import (
	"algo_assess/mqueue/internal/config"
	"algo_assess/mqueue/internal/svc"
	"context"
	"github.com/zeromicro/go-zero/core/service"
)

//返回所有消费者
func Mqs(c config.Config) []service.Service {
	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	var services []service.Service
	// kq ：消息队列.
	services = append(services, KqMqs(ctx, c, svcContext)...)

	return services
}
