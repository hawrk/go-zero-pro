package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	Mysql struct {
		DataSource     string
		IdleConn       int
		MaxOpenConn    int
		EnablePrintSQL int
	}
	// redis 配置
	//HRedis struct {
	//	Host string
	//	DB   int
	//}
}
