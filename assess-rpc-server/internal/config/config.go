package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource     string
		IdleConn       int
		MaxOpenConn    int
		EnablePrintSQL int
	}

	Kafka struct {
		Addrs         []string
		SHMarketTopic string
		SZMarketTopic string
		AlgoTopic     string
		ChildTopic    string
	}

}
