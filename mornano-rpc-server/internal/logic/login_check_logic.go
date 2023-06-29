package logic

import (
	"context"
	"errors"

	"algo_assess/mornano-rpc-server/internal/svc"
	"algo_assess/mornano-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginCheckLogic {
	return &LoginCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LoginCheck 登陆校验
func (l *LoginCheckLogic) LoginCheck(in *proto.LoginReq) (*proto.LoginRsp, error) {
	l.Logger.Info("LoginCheck, get req:", in)
	var userType uint
	var userPasswd, userName string
	var allow int32
	out, err := l.svcCtx.UserInfoRepo.CheckLogin(l.ctx, in.GetLoginName())
	if err != nil {
		l.Logger.Error("CheckLogin error:", err)
		return &proto.LoginRsp{
			Code:  209,
			Msg:   err.Error(),
			Allow: 0,
		}, nil
	}
	if len(out) <= 0 {
		l.Logger.Error("CheckLogin, no record")
		return &proto.LoginRsp{
			Code:  209,
			Msg:   errors.New("account_id not found").Error(),
			Allow: 0,
		}, nil
	}
	for _, v := range out {
		userType = v.UserType
		userPasswd = v.UserPasswd
		userName = v.UserName
		if in.GetPassword() == v.UserPasswd {
			allow = 1 // 把密码和总线的校验结果返回
		}
		break
	}
	return &proto.LoginRsp{
		Code:     200,
		Msg:      "success",
		Allow:    allow,
		Role:     int32(userType),
		Passwd:   userPasswd,
		UserName: userName,
	}, nil
}
