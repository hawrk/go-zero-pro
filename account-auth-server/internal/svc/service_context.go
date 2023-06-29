package svc

import (
	"account-auth/account-auth-server/internal/config"
	"account-auth/account-auth-server/internal/middleware"
	"account-auth/account-auth-server/repo"
	"fmt"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type ServiceContext struct {
	Config config.Config
	//HRedisClient *redis.Client
	Cors         rest.Middleware
	AuthUserRepo repo.AuthUserRepo
	AuthRoleRepo repo.AuthRoleRepo
	AuthMenuRepo repo.AuthMenuRepo
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化redis
	//rdb := redis.NewClient(&redis.Options{
	//	Addr: c.HRedis.Host,
	//	DB:   c.HRedis.DB,
	//})

	// 初始化gorm
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
		Config: c,
		//HRedisClient: rdb,
		Cors:         middleware.NewCorsMiddleware().Handle,
		AuthUserRepo: repo.NewAuthUserRepo(db),
		AuthRoleRepo: repo.NewAuthRole(db),
		AuthMenuRepo: repo.NewAuthMenu(db),
	}
}
