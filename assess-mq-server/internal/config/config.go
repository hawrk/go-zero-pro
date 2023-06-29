package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MornanoRPC zrpc.RpcClientConf
	Mysql      struct {
		DataSource     string
		IdleConn       int
		MaxOpenConn    int
		EnablePrintSQL int
	}
	WorkProcesser struct {
		GorontineNum         int
		ElapseMin            int
		EnableFirstPhase     bool
		EnableSecondPhase    bool
		EnableReloadCache    bool
		EnableCheckHeartBeat bool
		EnableConcurrency    bool
	}
	// kq
	AlgoPlatformOrderTradeConf kq.KqConf
	AlgoOrderTradeConf         kq.KqConf
	AlgoAccountInfoConf        kq.KqConf
	// perf fix
	//PerfFixOrderTradeConf     kq.KqConf
	//PerfFixAlgoOrderTradeConf kq.KqConf
}
