package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource     string
		IdleConn       int
		MaxOpenConn    int
		EnablePrintSQL int
	}
	WorkProcesser struct {
		GorontineNum int
		ElapseMin    int
	}
	// kq
	AlgoPlatformOrderTradeConf kq.KqConf
	AlgoOrderTradeConf         kq.KqConf
}
