package svc

import (
	"algo_assess/assess-mq-server/internal/config"
	"algo_assess/mornano-rpc-server/mornanoservice"
	"algo_assess/repo"
	"fmt"
	oriredis "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type ServiceContext struct {
	Config          config.Config
	RedisCluster    *oriredis.ClusterClient
	RedisClient     *redis.Redis
	OrderAssessRepo repo.OrderAssessRepo  // 一期绩效表
	OrderDetailRepo repo.OrderDetailRepo  // 子单落地明细表
	AlgoOrderRepo   repo.AlgoOrderRepo    // 母单落地表
	AccountRepo     repo.AccountInfoRepo  // 账户信息表
	AlgoBaseRepo    repo.AlgoInfoRepo     // 算法基础信息表
	SecurityRepo    repo.SecurityInfoRepo // 证券基础信息表
	ProfileRepo     repo.AlgoProfileRepo  // 算法画像明细表
	BusiConfigRepo  repo.BusiConfigRepo   // 基础配置表
	TimeLineRepo    repo.AlgoTimeLineRepo // 时间线图表
	SummaryRepo     repo.AlgoSummaryRepo  // 汇总表
	//ProfitRepo      repo.AlgoProfitRepo   // 收益表
	//BigCache        *bigcache.BigCache
	MornanoClient mornanoservice.MornanoService // 总线 rpc
	SzMarketRepo  repo.MarketLevelRepo          // 深市行情
	ShMarketRepo  repo.SHMarketLevelRepo        // 沪市行情
	// fix ---
	//FixAlgoOrderRepo repo.FixAlgoOrderRepo
	//FixOrderRepo     repo.FixOrderDetailRepo
	ProfileOrigRepo  repo.AlgoProfileOrigRepo  // 算法画像明细表
	TimeLineOrigRepo repo.AlgoTimeLineOrigRepo // 时间线图表
	SummaryOrigRepo  repo.AlgoSummaryOrigRepo  // 汇总表
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 1. 初始化 DB
	level := logger.Default.LogMode(logger.Warn)
	if c.Mysql.EnablePrintSQL == 1 {
		level = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{
		SkipDefaultTransaction: false,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "tb_",
		},
		Logger: level,
	})
	if err != nil {
		fmt.Println("init gorm fail:", err)
		//panic(err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(c.Mysql.IdleConn)    // 空闲连接池数量
	sqlDB.SetMaxOpenConns(c.Mysql.MaxOpenConn) // 最大连接数量
	sqlDB.SetConnMaxLifetime(time.Hour)        // 设置可复用最大时间

	//config := bigcache.Config{
	//	Shards:      1024,
	//	CleanWindow: 24 * time.Hour,
	//}
	//cache, initErr := bigcache.NewBigCache(config)
	//if initErr != nil {
	//	fmt.Println("init big cache error")
	//}
	// 2. 初始化redis
	// 注意：这里redis使用单机版的话，直接使用框架自带的功能就可以了，如果需要集群，需要初始化原生的go-redis库进行初始化
	//ips := strings.Split(c.Redis.Host, ",")
	//rdb := oriredis.NewClusterClient(&oriredis.ClusterOptions{
	//	Addrs: ips,
	//})
	return &ServiceContext{
		Config: c,
		//RedisCluster:     rdb,   // 需要集群时，把这个替换成RedisClient
		RedisClient:     redis.New(c.Redis.Host),
		OrderAssessRepo: repo.NewAlgoAssessRepo(db),
		OrderDetailRepo: repo.NewOrderDetailRepo(db),
		AlgoOrderRepo:   repo.NewDefaultAlgoOrder(db),
		AccountRepo:     repo.NewAccountInfo(db),
		AlgoBaseRepo:    repo.NewAlgoInfoRepo(db),
		SecurityRepo:    repo.NewSecurityInfo(db),
		ProfileRepo:     repo.NewAlgoProfile(db),
		BusiConfigRepo:  repo.NewDefaultBusiConfig(db),
		TimeLineRepo:    repo.NewAlgoTimeLine(db),
		SummaryRepo:     repo.NewAlgoSummary(db),
		//ProfitRepo:      repo.NewAlgoProfit(db),
		//BigCache:        cache,
		MornanoClient: mornanoservice.NewMornanoService(zrpc.MustNewClient(c.MornanoRPC)),
		SzMarketRepo:  repo.NewMarketLevelRepo(db),
		ShMarketRepo:  repo.NewSHMarketLevelRepo(db),
		// fix
		//FixAlgoOrderRepo: repo.NewPerfFixAlgoOrder(db),
		//FixOrderRepo:     repo.NewFixOrderDetailRepo(db),
		ProfileOrigRepo:  repo.NewAlgoProfileOrig(db),
		TimeLineOrigRepo: repo.NewAlgoTimeLineOrig(db),
		SummaryOrigRepo:  repo.NewAlgoSummaryOrig(db),
	}
}
