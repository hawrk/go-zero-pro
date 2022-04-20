package svc

import (
	"algo_assess/assess-rpc-server/internal/config"
	"algo_assess/repo"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type ServiceContext struct {
	Config config.Config
	// DB sqlx.SqlConn   // 不用框架自带的，比较难用
	//DB          *gorm.DB
	OrderAssessRepo repo.OrderAssessRepo
	RedisClient     *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化Gorm
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

	return &ServiceContext{
		Config:          c,
		OrderAssessRepo: repo.NewAlgoAssessRepo(db),
		RedisClient:     redis.New(c.Redis.Host),
	}
}
