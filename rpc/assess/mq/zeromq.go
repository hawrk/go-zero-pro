/*
@Time : 2022/3/17 6:28 下午
@Author : hawrkchen
@File : zeromq.go
@Desc:
*/

package mq

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
)

var (
	MQServer *zmq.Socket // 接收消息
	MQClient *zmq.Socket // 发送消息
)

func init() {
	initServer()
	initClient()
}

func initServer() {
	zctx, err := zmq.NewContext()
	if err != nil {
		fmt.Println("zmq init server new context fail:", err)
		return
	}
	MQServer, err := zctx.NewSocket(zmq.REQ)
	if err != nil {
		fmt.Println("zmq init server new socket fail:", err)
		return
	}
	if err := MQServer.Bind("tcp://127.0.0.1:5555"); err != nil {
		fmt.Println("zmq init server bind fail:", err)
		return
	}

	fmt.Println("zmq server init success...")
}

func initClient() {
	zctx, err := zmq.NewContext()
	if err != nil {
		fmt.Println("zmq init client new context fail:", err)
		return
	}
	MQClient, err := zctx.NewSocket(zmq.REQ)
	if err != nil {
		fmt.Println("zmq init client new socket fail:", err)
		return
	}
	if err := MQClient.Connect("tcp://127.0.0.1:5555"); err != nil {
		fmt.Println("zmq init client connect fail:", err)
		return
	}
	fmt.Println("zmq client init success...")
}
