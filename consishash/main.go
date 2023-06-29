// Package consishash
/*
 Author: hawrkchen
 Date: 2023/6/19 16:14
 Desc:
*/
package main

import (
	"fmt"
	"github.com/spf13/cast"
	"time"
)

func main() {
	t := time.Now()
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	fmt.Println("addTime:", addTime)
	return
	fmt.Println("get consistent hash...")
	h := NewHashRing(4, nil)
	h.AddNodes("192.168.1.1", "192.168.1.2", "192.168.1.3", "192.168.1.4")

	for i := 0; i < 50; i++ {
		v := h.GetNode(cast.ToString(i))
		fmt.Println("get key:", i, ",v:", v)
	}
}
