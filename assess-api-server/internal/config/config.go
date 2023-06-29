package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	AssessRPC   zrpc.RpcClientConf
	MornanoRPC  zrpc.RpcClientConf
	AssessMQRPC zrpc.RpcClientConf
	MarketMQRPC zrpc.RpcClientConf
	//PerformanceMQRPC zrpc.RpcClientConf

	WorkControl struct {
		EnableFakeMsg bool
		EnableFakeDay bool
	}
	WebSocket struct {
		DurationTime int
	}
	// redis 只作鉴权用
	HRedis struct {
		Host      string
		DB        int
		NeedCheck bool
	}
	// zk
	Zookeeper struct {
		Host           string
		TimeOut        int
		EnableRegister bool
	}
	// JWT 鉴权
	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
		TokenCheck   bool
	}

	Domain struct {
		DomainName string
		UseDomain  bool
	}

	AccountAuth struct {
		UrlPrefix string
	}
}
