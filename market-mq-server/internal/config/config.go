package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// mq rpc
	AssessMQRPC zrpc.RpcClientConf
	//kq
	AlgoPlatformMarketConf kq.KqConf
}
