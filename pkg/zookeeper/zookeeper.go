// Package zookeeper
/*
 Author: hawrkchen
 Date: 2022/6/8 16:09
 Desc:
*/
package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

type Zookeeper struct {
	ZKHost  []string      // zookeeper 服务器节点
	TimeOut time.Duration // 超时时间
	Conn    *zk.Conn
}

func NewZookeeper(host string, t int) *Zookeeper {
	return &Zookeeper{
		ZKHost:  []string{host},
		TimeOut: time.Duration(t),
		Conn:    nil,
	}
}

// Connect 连接
func (z *Zookeeper) Connect() error {
	conn, _, err := zk.Connect(z.ZKHost, time.Second*z.TimeOut)
	if err != nil {
		return err
	}
	z.Conn = conn
	return nil
}

func (z *Zookeeper) RegisterServer(path, serverHost string) error {
	exist, _, err := z.Conn.Exists(path)
	if !exist {
		_, err = z.Conn.Create(path, nil, 0, zk.WorldACL(zk.PermAll)) // 根节点持久化
		if err != nil {
			return err
		}
	}
	_, err = z.Conn.Create(path+"/"+serverHost, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	if err != nil {
		return err
	}
	return nil
}
