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
	//Deployment struct {
	//	Env string
	//}
	// mq rpc
	AssessMQRPC zrpc.RpcClientConf
	//kq-sz
	AlgoPlatformMarketConf kq.KqConf
	//kq-sh
	AlgoPlatFormSHMarketConf kq.KqConf
	// 数据修复  -sz
	PerfFixSZMarketConf kq.KqConf
	// 数据修复 -sh
	PerfFixSHMarketConf kq.KqConf
}
