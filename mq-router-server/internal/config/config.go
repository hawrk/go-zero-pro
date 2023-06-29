package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	// mq  master rpc
	AssessMQMRPC zrpc.RpcClientConf

	// mq slave rpc
	AssessMQSRPC zrpc.RpcClientConf

	// kq
	AlgoOrderTradeConf  kq.KqConf
	ChildOrderTradeConf kq.KqConf
}
