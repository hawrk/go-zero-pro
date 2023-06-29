package svc

import (
	"algo_assess/assess-api-server/internal/config"
	"algo_assess/assess-api-server/internal/middleware"
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/assess-rpc-server/assessservice"
	mkservice "algo_assess/market-mq-server/marketservice"
	"algo_assess/mornano-rpc-server/mornanoservice"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	AssessClient   assessservice.AssessService
	MornanoClient  mornanoservice.MornanoService
	AssessMQClient mqservice.AssessMqService
	MarketMQClient mkservice.MarketService
	Interceptor    rest.Middleware
	//HRedisClient   *redis.ClusterClient
	HRedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化redis, 官方redis不支持选择数据库
	// 支持集群，但不能再指定DB
	//ips := strings.Split(c.HRedis.Host, ",")
	//rdb := redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs: ips,
	//})
	rdb := redis.NewClient(&redis.Options{
		Addr: c.HRedis.Host,
		DB:   c.HRedis.DB,
	})

	//result, err := rdb.Ping(context.Background()).Result()
	//fmt.Println("result:", result, "err:", err)

	return &ServiceContext{
		Config:         c,
		AssessClient:   assessservice.NewAssessService(zrpc.MustNewClient(c.AssessRPC)),
		MornanoClient:  mornanoservice.NewMornanoService(zrpc.MustNewClient(c.MornanoRPC)),
		AssessMQClient: mqservice.NewAssessMqService(zrpc.MustNewClient(c.AssessMQRPC)),
		MarketMQClient: mkservice.NewMarketService(zrpc.MustNewClient(c.MarketMQRPC)),
		Interceptor:    middleware.NewInterceptorMiddleware().Handle,
		HRedisClient:   rdb,
	}
}
