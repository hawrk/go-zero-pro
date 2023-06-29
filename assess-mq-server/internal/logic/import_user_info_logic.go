package logic

import (
	"algo_assess/assess-mq-server/internal/dao"
	"context"
	"errors"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewImportUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportUserInfoLogic {
	return &ImportUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ImportUserInfo 配置： 用户信息导入
func (l *ImportUserInfoLogic) ImportUserInfo(in *proto.ImportUserReq) (*proto.ImportUserRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in ImportUserInfo, process")
	if len(in.GetInfos()) == 0 {
		l.Logger.Error("empty list...")
		return &proto.ImportUserRsp{
			Code:   300,
			Msg:    errors.New("empty list").Error(),
			Result: 2,
		}, nil
	}
	// 缓存一份数据，用来判断导入的数据是插入还是更新
	m := make(map[string]struct{})
	dao.GAccountMap.RWMutex.RLock()
	for k, _ := range dao.GAccountMap.Account {
		m[k] = struct{}{}
	}
	dao.GAccountMap.RWMutex.RUnlock()

	for _, v := range in.GetInfos() {
		if _, exist := m[v.UserId]; exist {
			if err := l.svcCtx.AccountRepo.ImportUserUpdate(l.ctx, v); err != nil {
				l.Logger.Error("ImportUserUpdate error :", err)
				continue
			}
		} else {
			if err := l.svcCtx.AccountRepo.ImportUserCreate(l.ctx, v); err != nil {
				l.Logger.Error("ImportUserCreate error:", err)
				continue
			}
		}
	}

	return &proto.ImportUserRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}
