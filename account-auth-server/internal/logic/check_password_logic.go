package logic

import (
	"context"

	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckPasswordLogic {
	return &CheckPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckPasswordLogic) CheckPassword(req *types.CheckPasswordReq) (resp *types.CheckPasswordRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in CheckPassword, get req:%+v", *req)
	passwd, err := l.svcCtx.AuthUserRepo.CheckUserPasswd(l.ctx, req.UserId, req.ChanType)
	if err != nil {
		l.Logger.Error("CheckUserPasswd error:", err)
		return &types.CheckPasswordRsp{
			Code:   360,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	if passwd != req.OriginPassword {
		l.Logger.Info("password not match, origin:", req.OriginPassword, ",db:", passwd)
		return &types.CheckPasswordRsp{
			Code:   360,
			Msg:    "password not match",
			Result: 2,
		}, nil
	}

	return &types.CheckPasswordRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}
