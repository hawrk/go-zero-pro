// Package tools
/*
 Author: hawrkchen
 Date: 2023/4/13 10:12
 Desc:
*/
package tools

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"strconv"
	"time"
)

var node *snowflake.Node

func init() {
	var st time.Time
	st, err := time.Parse("2006-01-02", "2021-03-01")
	if err != nil {
		fmt.Println("parse time err:", err)
		return
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(1)
	if err != nil {
		fmt.Println("general new node err:", err)
	}
}

// GeneralID 生成雪花算法ID， 取前9位
func GeneralID() int64 {
	// 正常雪花算法会生成18位长度的自增ID，但我们没有那么高的并发量，在秒级内生成的ID不重复即可
	// 取前9位即能保证不会重复
	id := node.Generate().String()[:9]
	r, _ := strconv.ParseInt(id, 10, 64)
	return r
}
