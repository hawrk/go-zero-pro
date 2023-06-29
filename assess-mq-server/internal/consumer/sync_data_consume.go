// Package consumer
/*
 Author: hawrkchen
 Date: 2022/6/24 15:41
 Desc: 交易账户消息推送
*/
package consumer

import (
	"algo_assess/assess-mq-server/internal/dao"
	"algo_assess/assess-mq-server/internal/svc"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/pkg/tools"
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
	"strings"
)

type SyncDataInfo struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAccountInfo(ctx context.Context, svcCtx *svc.ServiceContext) *SyncDataInfo {
	return &SyncDataInfo{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *SyncDataInfo) Consume(key string, val string) error {
	s.Logger.Info("-----------------account info start------------------")
	data := pb.DataSyncPerf{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return err
	}
	s.Logger.Infof("get sync data:%+v", data)
	if data.GetMsgType() == 1 {
		s.SyncAccountInfo(data.GetUserInfo())
	} else if data.GetMsgType() == 2 {
		s.SyncAlgoInfo(data.GetAlgoInfo())
	} else if data.GetMsgType() == 3 {
		s.SyncSecurityInfo(data.GetSecInfo())
	} else {
		s.Logger.Error(" sync data type not support")
		return nil
	}
	return nil
}

func (s *SyncDataInfo) SyncAccountInfo(data []*pb.UserInfoPerf) error {
	for _, v := range data {
		//s.Logger.Infof("get account info:%+v", data)
		if v.GetUserId() == "" {
			s.Logger.Error("account key invalid")
			return nil
		}
		// 1. 更新本地缓存
		dao.GAccountMap.RWMutex.Lock()
		userId := strings.TrimSpace(tools.RMu0000(v.GetUserId()))
		dao.GAccountMap.AccountKey[int(v.GetId())] = userId // 先更新key
		info := dao.AccountInfo{
			UserName:  tools.RMu0000(v.GetUserName()),
			UserType:  int(v.GetUserType()),
			ParUserId: dao.GetParentUserId(dao.GAccountMap.AccountKey, v.GetUuserId()),
		}
		dao.GAccountMap.Account[userId] = info
		dao.GAccountMap.RWMutex.Unlock()

		// 2.更新本地DB
		infos, err := s.svcCtx.AccountRepo.GetAccountInfoById(s.ctx, int(v.Id))
		if err != nil {
			s.Logger.Error("query account info error:", err)
			return nil
		}
		if len(infos) > 0 { // 已有记录，更新值
			if err := s.svcCtx.AccountRepo.UpdateAccountInfoById(s.ctx, int(v.Id), v); err != nil {
				s.Logger.Error("update account info error:", err)
			}
		} else { //插入
			if err := s.svcCtx.AccountRepo.CreateAccountInfo(s.ctx, v); err != nil {
				s.Logger.Error(" insert account info error :", err)
			}
		}
	}
	return nil
}

func (s *SyncDataInfo) SyncAlgoInfo(data []*pb.AlgoInfoPerf) error {
	for _, v := range data {
		//s.Logger.Infof("get sync algo data:%+v", v)
		if v.GetId() == 0 {
			s.Logger.Error("sync, invalid algoid")
			return errors.New("sync, invalid algoid")
		}
		info := dao.AlgoInfo{
			AlgoName:     tools.RMu0000(v.GetAlgoName()),
			AlgoType:     int(v.GetAlgorithmType()),
			AlgoTypeName: tools.RMu0000(v.GetAlgorithmTypeName()),
			Provider:     tools.RMu0000(v.GetProviderName()),
		}
		dao.GAlgoBaseMap.RWMutex.Lock()
		dao.GAlgoBaseMap.AlgoBase[int(v.GetId())] = info
		dao.GAlgoBaseMap.RWMutex.Unlock()

		// 2.更新本地DB
		infos, err := s.svcCtx.AlgoBaseRepo.GetAlgoBaseById(s.ctx, v.GetId())
		if err != nil {
			s.Logger.Error("query algo info error:", err)
			return nil
		}
		if len(infos) > 0 { // 已有记录，更新值
			if err := s.svcCtx.AlgoBaseRepo.UpdateAlgoBaseById(s.ctx, v.GetId(), v); err != nil {
				s.Logger.Error("update algo base info error:", err)
			}
		} else { //插入
			if err := s.svcCtx.AlgoBaseRepo.CreateAlgoBaseInfo(s.ctx, v); err != nil {
				s.Logger.Error(" insert algo base info error :", err)
			}
		}
	}
	return nil
}

func (s *SyncDataInfo) SyncSecurityInfo(data []*pb.SecurityInfoPerf) error {
	for _, v := range data {
		//s.Logger.Infof("get security info :%+v", v)
		if len(v.GetSecurityId()) == 0 {
			s.Logger.Error("sync invalid security id")
			return errors.New("sync invalid security id")
		}
		info := dao.SecurityInfo{
			SecuritySource:  v.GetSecurityIdSource(),
			SecurityName:    tools.RMu0000(v.GetSecurityName()),
			PreClosePrice:   v.GetPrevClosePx(),
			Status:          int(v.GetSecurityStatus()),
			UpperLimitPrice: int64(v.GetUpperLimitPrice()),
			LowerLimitPrice: int64(v.GetLowerLimitPrice()),
		}
		dao.GSecurityMap.RWMutex.Lock()
		dao.GSecurityMap.SecurityBase[v.GetSecurityId()] = info
		dao.GSecurityMap.RWMutex.Unlock()

		// 2.更新本地DB
		infos, err := s.svcCtx.SecurityRepo.GetSecurityInfoById(s.ctx, v.GetSecurityId())
		if err != nil {
			s.Logger.Error("query security info error:", err)
			return nil
		}
		if len(infos) > 0 { // 已有记录，更新值
			if err := s.svcCtx.SecurityRepo.UpdateSecurityById(s.ctx, v.GetSecurityId(), v); err != nil {
				s.Logger.Error("update security base info error:", err)
			}
		} else { //插入
			if err := s.svcCtx.SecurityRepo.CreateSecurityInfo(s.ctx, v); err != nil {
				s.Logger.Error(" insert security base info error :", err)
			}
		}
	}
	return nil
}
