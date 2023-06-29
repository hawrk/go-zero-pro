package svc

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/market-mq-server/internal/config"
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
	Config            config.Config
	RedisCluster     *oriredis.ClusterClient
	RedisClient       *redis.Redis
	AssessMQClient    mqservice.AssessMqService
	MarketLevelRepo   repo.MarketLevelRepo
	ShMarketLevelRepo repo.SHMarketLevelRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
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
	// 2. 初始化redis
	// 注意：这里redis使用单机版的话，直接使用框架自带的功能就可以了，如果需要集群，需要初始化原生的go-redis库进行初始化
	//ips := strings.Split(c.Redis.Host, ",")
	//rdb := oriredis.NewClusterClient(&oriredis.ClusterOptions{
	//	Addrs: ips,
	//})

	return &ServiceContext{
		Config:            c,
		//RedisCluster:     rdb,   // 需要集群时，把这个替换成RedisClient
		RedisClient:       redis.New(c.Redis.Host),
		AssessMQClient:    mqservice.NewAssessMqService(zrpc.MustNewClient(c.AssessMQRPC)),
		MarketLevelRepo:   repo.NewMarketLevelRepo(db),
		ShMarketLevelRepo: repo.NewSHMarketLevelRepo(db),
	}
}
