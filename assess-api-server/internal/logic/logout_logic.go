package logic

import (
	"algo_assess/global"
	"context"
	"fmt"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in Logout, get req:%+v", req)
	// 直接把该key删除
	tokenKey := fmt.Sprintf("%s:%s", global.UserTokenKey, req.UserName)
	_, err = l.svcCtx.HRedisClient.Del(l.ctx, tokenKey).Result()
	if err != nil {
		l.Logger.Error("del token key error:", err)
	}

	return &types.LogoutRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}
