package logic

import (
	"account-auth/account-auth-server/global"
	"context"
	"time"

	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 登陆接口，校验用户权限新增的用户
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginRsp, err error) {
	l.Logger.Infof("in Login, get req:%+v", *req)
	if req.UserId == "" {
		return &types.LoginRsp{
			Code:  360,
			Msg:   "user_id not found",
			Allow: 0,
		}, nil
	}
	out, err := l.svcCtx.AuthUserRepo.GetAuthUserById(l.ctx, req.UserId, req.ChanType, 0)
	if err != nil {
		l.Logger.Error("GetAuthUserById error:", err)
		return &types.LoginRsp{
			Code:  360,
			Msg:   err.Error(),
			Allow: 0,
		}, nil
	}
	// 如果总线有记录， 并且UserID为空，则表示该用户在绩效这边没有记录，需要插入---只针对绩效平台
	if req.BusPasswd != "" && req.BusAllow == 1 && out.UserId == "" {
		r := &types.UserModfiyReq{
			UserId:   req.UserId,
			Password: req.Passwd,
			UserName: req.BusUserName,
			ChanType: req.ChanType,
		}
		if err := l.svcCtx.AuthUserRepo.CreateAuthUser(l.ctx, r); err != nil {
			l.Logger.Error("create Auth user error:", err)
		}
		// 绩效同步数据后，直接校验总线的密码
		if req.BusPasswd != req.Passwd {
			l.Logger.Error("bus user info  password not match, buspasswd:", req.BusPasswd, ", req:", req.Passwd)
			return &types.LoginRsp{
				Code:  360,
				Msg:   "用户名或密码错误",
				Allow: 0,
			}, nil
		}
		return &types.LoginRsp{
			Code:       200,
			Msg:        "success",
			Allow:      1,
			UserType:   2, // 默认普通用户
			FirstLogin: 0, // 默认首次登陆
		}, nil
	} else if req.BusPasswd == "" && out.UserId == "" { // 两边都无数据，密码校验失败
		l.Logger.Info("both bus and assess no data")
		return &types.LoginRsp{
			Code:  360,
			Msg:   "用户名或密码错误",
			Allow: 0,
		}, nil
	}
	// 绩效有数据
	if out.UserPasswd != req.Passwd {
		l.Logger.Error("password not match, req:", req.Passwd, ", db:", out.UserPasswd)
		return &types.LoginRsp{
			Code:  360,
			Msg:   "用户名或密码错误",
			Allow: 0,
		}, nil
	}
	// 增加校验状态
	if out.Status == global.UserStatusDel {
		l.Logger.Error("status has delete")
		return &types.LoginRsp{
			Code:  360,
			Msg:   "该用户已删除",
			Allow: 0,
		}, nil
	} else if out.Status == global.UserStatusForbid {
		l.Logger.Error("status has forbidden")
		return &types.LoginRsp{
			Code:  360,
			Msg:   "该用户已被禁用",
			Allow: 0,
		}, nil
	}
	// 增加密码过期时间校验
	// 如果当前过期时间为空，则从当前时间算当前时间
	if out.ExpireTime == 0 {
		if err := l.svcCtx.AuthUserRepo.UpdateExpireTime(l.ctx, req.UserId, req.ChanType); err != nil {
			l.Logger.Error("UpdateExpireTime error:", err)
		}
	} else {
		curTime := time.Now().Unix()
		if curTime >= out.ExpireTime { // 密码已过期
			return &types.LoginRsp{
				Code:       200,
				Msg:        "密码已过期,请重新设置密码",
				Allow:      1,
				UserType:   out.UserType,
				FirstLogin: 2,
			}, nil
		}
	}

	return &types.LoginRsp{
		Code:       200,
		Msg:        "success",
		Allow:      1,
		UserType:   out.UserType,
		FirstLogin: out.FirstLogin,
	}, nil
}
