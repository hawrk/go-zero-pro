package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	AssessRPC   zrpc.RpcClientConf
	AssessMQRPC zrpc.RpcClientConf
	MarketMQRPC zrpc.RpcClientConf

	WebSocket struct {
		DurationTime int
	}
}
