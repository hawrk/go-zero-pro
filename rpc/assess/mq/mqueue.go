// Package mq
/*
 Author: hawrkchen
 Date: 2022/3/21 9:51
 Desc:
*/
package mq

import "github.com/zeromicro/go-zero/core/logx"

func StartZeroMQ() {
	logx.Logger.Info("start to listen zero mq ....")

	for {
		msg, _ := MQServer.Recv(0)
		logx.Logger.Info("get mq message :", msg)
		// TODO: MQ处理
	}
}
