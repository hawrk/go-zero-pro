// Package consumer
/*
 Author: hawrkchen
 Date: 2022/4/24 16:20
 Desc:
*/
package consumer

import (
	"sync"
	"time"
)

var Location *time.Location
var once sync.Once

func init() {
	once.Do(func() {
		Location, _ = time.LoadLocation("Asia/Shanghai")
	})
}
