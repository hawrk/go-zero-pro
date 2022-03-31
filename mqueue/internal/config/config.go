// Package config
/*
 Author: hawrkchen
 Date: 2022/3/24 10:01
 Desc:
*/
package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	service.ServiceConf
	Mysql struct {
		DataSource     string
		IdleConn       int
		MaxOpenConn    int
		EnablePrintSQL int
	}
	Redis redis.RedisConf
	// kq
	AlgoPlatformOrderTradeConf  kq.KqConf
	AlgoPlatformOrderResultConf kq.KqConf
	AlgoPlatformMarketConf      kq.KqConf
}
