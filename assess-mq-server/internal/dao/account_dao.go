// Package dao
/*
 Author: hawrkchen
 Date: 2022/7/6 15:23
 Desc:
*/
package dao

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/pkg/tools"
	"context"
	"errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"sync"
)

var GAccountMap = AccountMap{
	Account:    make(map[string]AccountInfo),
	AccountKey: make(map[int]string),
}

type AccountMap struct {
	Account    map[string]AccountInfo // key-> user_id,
	AccountKey map[int]string         // 辅助map, key -> account_id , value -> user_id , 用来根据par_user_id 找到管理员的user_id
	sync.RWMutex
}

type AccountInfo struct {
	UserName  string
	UserType  int
	ParUserId []string // 这个就直接映射成管理员的user_id了 ---这里需要扩展成支持多个管理员
}

type AccountDb struct{}

func (d *AccountDb) LoadAccountInfo(svcContext *svc.ServiceContext) error {
	infos, err := svcContext.AccountRepo.GetAccountInfos(context.Background())
	if err != nil {
		return err
	}
	GAccountMap.RWMutex.Lock()
	defer GAccountMap.RWMutex.Unlock()

	for _, v := range infos {
		GAccountMap.AccountKey[v.AccountId] = strings.TrimSpace(v.UserId)
	}

	for _, v := range infos {
		//GAccountMap.AccountKey[v.AccountId] = strings.TrimSpace(v.UserId)
		info := AccountInfo{
			UserName:  tools.RMu0000(v.UserName),
			UserType:  v.UserType,
			ParUserId: GetParentUserId(GAccountMap.AccountKey, v.ParUserId),
		}
		GAccountMap.Account[v.UserId] = info
		//logx.Infof("get userID:%s, %+v", v.UserId, info)
	}
	if len(GAccountMap.Account) <= 0 {
		return errors.New("account info no data")
	}

	logx.Info("load account info success, len:", len(GAccountMap.Account))
	return nil
}

func loadAccountInfo(svcContext *svc.ServiceContext) error {
	infos, err := svcContext.AccountRepo.GetAccountInfos(context.Background())
	if err != nil {
		return err
	}
	GAccountMap.RWMutex.Lock()
	defer GAccountMap.RWMutex.Unlock()

	for _, v := range infos {
		GAccountMap.AccountKey[v.AccountId] = strings.TrimSpace(v.UserId)

		info := AccountInfo{
			UserName:  tools.RMu0000(v.UserName),
			UserType:  v.UserType,
			ParUserId: GetParentUserId(GAccountMap.AccountKey, v.ParUserId),
		}
		GAccountMap.Account[v.UserId] = info
		//logx.Infof("get userID:%s, %+v", v.UserId, info)
	}
	if len(GAccountMap.Account) <= 0 {
		return errors.New("account info no data")
	}

	logx.Info("load account info success, len:", len(GAccountMap.Account))
	return nil
}

// GetParentUserId 根据传入的管理员ID列表找到其管理员账号
func GetParentUserId(m map[int]string, parId string) []string {
	ids := strings.Split(parId, ",")
	if len(ids) <= 0 {
		return []string{}
	}
	var ret []string
	for _, v := range ids {
		if val, exist := m[cast.ToInt(v)]; exist {
			ret = append(ret, val)
		}
	}
	return ret
}
