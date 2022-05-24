package svc

import (
	"algo_assess/assess-mq-server/internal/config"
	"algo_assess/repo"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type ServiceContext struct {
	Config          config.Config
	RedisClient     *redis.Redis
	OrderAssessRepo repo.OrderAssessRepo
	OrderDetailRepo repo.OrderDetailRepo
	AlgoOrderRepo   repo.AlgoOrderRepo
	BigCache        *bigcache.BigCache
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

	config := bigcache.Config{
		Shards:      1024,
		CleanWindow: 24 * time.Hour,
	}
	cache, initErr := bigcache.NewBigCache(config)
	if initErr != nil {
		fmt.Println("init big cache error")
	}

	return &ServiceContext{
		Config:          c,
		RedisClient:     redis.New(c.Redis.Host),
		OrderAssessRepo: repo.NewAlgoAssessRepo(db),
		OrderDetailRepo: repo.NewOrderDetailRepo(db),
		AlgoOrderRepo:   repo.NewDefaultAlgoOrder(db),
		BigCache:        cache,
	}
}
