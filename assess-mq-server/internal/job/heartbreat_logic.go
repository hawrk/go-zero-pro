// Package job
/*
 Author: hawrkchen
 Date: 2022/12/1 15:34
 Desc:
*/
package job

import (
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"fmt"
	"strings"
	"time"
)

func (o *AssessJob) KeepHeartBeat() {
	if !o.s.Config.WorkProcesser.EnableCheckHeartBeat {
		return
	}
	// 1. 查redis 中有没有数据
	v, err := o.s.RedisClient.Get(global.HeartBeatKey)
	if err != nil {
		o.Logger.Error("heart beat key get error:", err)
		return
	}
	// 2. 没有,则写入
	if v == "" {
		o.SyncHeartBeat()
		return
	}
	// 3.有数据了，先判断一下是不是自已的
	s := strings.Split(v, ":")
	if len(s) < 2 {
		o.Logger.Error("heart beat value invalid")
		return
	}
	if tools.GetLocalIP() == s[0] { // 是自己的，则更新时间戳
		o.SyncHeartBeat()
	} else {
		// 不是自已的话，判断该时间戳是否超过时限了，如果超过了，则表示另外一个主机可能挂了，这里更新成本机的信息
		//now := time.Now().Unix()
		//if now-cast.ToInt64(s[1]) > 90 {    // 注： 如果出现另外一台主机挂掉没有更新到时间戳的话，这里还不能直接更新成自己的，
		// 因为直接更新的话，在consume主流程就有可能reload不到快照数据
		//	o.SyncHeartBeat()
		//}
		// 没超就不用管他了
	}
}

func (o *AssessJob) SyncHeartBeat() {
	// 写入心跳的key, 格式： IP:时间戳(精确到秒)
	now := time.Now().Unix()
	ip := tools.GetLocalIP()
	v := fmt.Sprintf("%s:%d", ip, now)
	//o.Logger.Info("get heartbeat value:", v)
	if err := o.s.RedisClient.Setex(global.HeartBeatKey, v, 60*60); err != nil {
		o.Logger.Error("set heart beat key error:", err)
	}
}
