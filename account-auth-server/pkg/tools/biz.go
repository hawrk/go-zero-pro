// Package tools
/*
 Author: hawrkchen
 Date: 2022/10/13 14:20
 Desc:
*/
package tools

import (
	"account-auth/account-auth-server/global"
	"time"
)

// Time2String time 格式转换成字符串
func Time2String(time2 time.Time) string {
	if time2.IsZero() {
		return ""
	}
	return time2.Format(global.TimeFormat)
}

// QuarterExpireTime 每季度的登陆过期时间点
func QuarterExpireTime() int64 {
	year := time.Now().Format("2006")
	month := int(time.Now().Month())
	var expireDay string
	if month >= 1 && month <= 3 {
		expireDay = year + "-03-31 23:59:59"
	} else if month >= 4 && month <= 6 {
		expireDay = year + "-06-30 23:59:59"
	} else if month >= 7 && month <= 9 {
		expireDay = year + "-09-30 23:59:59"
	} else {
		expireDay = year + "-12-31 23:59:59"
	}
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	tl, _ := time.ParseInLocation("2006-01-02 15:04:05", expireDay, Loc)
	return tl.Unix()
}
