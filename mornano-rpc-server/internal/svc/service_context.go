package svc

import (
	"algo_assess/busrepo"
	"algo_assess/mornano-rpc-server/internal/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type ServiceContext struct {
	Config           config.Config
	UserInfoRepo     busrepo.UserInfoRepo
	AlgoInfoRepo     busrepo.AlgoInfoRepo
	AlgoGroupRepo    busrepo.AlgoGroupRepo
	AssetRepo        busrepo.AssetInfoRepo
	UserPositionRepo busrepo.UserPositionRepo
	SecurityRepo     busrepo.SecurityInfoRepo
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
		Config:           c,
		UserInfoRepo:     busrepo.NewUserInfo(db),
		AlgoInfoRepo:     busrepo.NewAlgoInfo(db),
		AlgoGroupRepo:    busrepo.NewAlgoGroup(db),
		AssetRepo:        busrepo.NewAssetInfo(db),
		UserPositionRepo: busrepo.NewUserPosition(db),
		SecurityRepo:     busrepo.NewSecurityInfo(db),
	}
}
