// Package dao
/*
 Author: hawrkchen
 Date: 2022/12/1 16:28
 Desc:
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"fmt"
	"time"
)

func SetHeartBeat(svcContext *svc.ServiceContext) error {
	// 启动时，如果已经有了，则跳过
	v, err := svcContext.RedisClient.Get(global.HeartBeatKey)
	if err != nil {
		return err
	}
	if v != "" {
		return nil
	}
	// 写入心跳的key, 格式： IP:时间戳(精确到秒)
	now := time.Now().Unix()
	ip := tools.GetLocalIP()
	val := fmt.Sprintf("%s:%d", ip, now)
	if err := svcContext.RedisClient.Setex(global.HeartBeatKey, val, 60*60); err != nil {
		return err
	}
	return nil
}
