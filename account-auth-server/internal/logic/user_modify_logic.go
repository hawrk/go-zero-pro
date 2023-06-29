package logic

import (
	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserModifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserModifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserModifyLogic {
	return &UserModifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserModifyLogic) UserModify(req *types.UserModfiyReq) (resp *types.UserModifyRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in UserModify, get req:%+v", *req)
	if req.OperType == 1 { // 新增
		// 新增需要校验入参
		if req.UserId == "" || req.UserName == "" || req.RoleId == 0 || req.Password == "" {
			l.Logger.Error("param not enough")
			return &types.UserModifyRsp{
				Code:   320,
				Msg:    "user_id/user_name/role_id/password param missed",
				Result: 2,
			}, nil
		}
		if err := l.svcCtx.AuthUserRepo.CreateAuthUser(l.ctx, req); err != nil {
			l.Logger.Error("create Auth user error:", err)
			return &types.UserModifyRsp{
				Code:   322,
				Msg:    "该用户已存在",
				Result: 2,
			}, nil
		}
	} else if req.OperType == 2 { // 修改
		// 用户信息修改，如果是首次登陆修改密码的请求，则修改首次登陆标识
		// 有修改密码的请求，都要更新密码过期时间
		// 修改密码时，如果是总线的用户，在登陆检查的时候已经同步到绩效这边了，这时直接更新密码就行
		if req.UserId == "" {
			l.Logger.Error("field user_id not set")
			return &types.UserModifyRsp{
				Code:   320,
				Msg:    "用户ID必填",
				Result: 2,
			}, nil
		}
		if req.Password != "" {
			if err := l.svcCtx.AuthUserRepo.UpdateAuthUserPasswd(l.ctx, req); err != nil {
				l.Logger.Error("UpdateAuthUserPasswd error:", err)
				return &types.UserModifyRsp{
					Code:   320,
					Msg:    err.Error(),
					Result: 2,
				}, nil
			}
		} else {
			if err := l.svcCtx.AuthUserRepo.UpdateAuthUser(l.ctx, req); err != nil {
				l.Logger.Error("update Auth user error:", err)
				return &types.UserModifyRsp{
					Code:   320,
					Msg:    err.Error(),
					Result: 2,
				}, nil
			}
		}

	} else if req.OperType == 3 { // 删除,直接改状态即可
		if err := l.svcCtx.AuthUserRepo.DelAuthUser(l.ctx, req.UserId, req.ChanType); err != nil {
			l.Logger.Error("del Auth user error:", err)
			return &types.UserModifyRsp{
				Code:   320,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}
	} else {
		l.Logger.Error("unsupported user oper_type", req.OperType)
		return &types.UserModifyRsp{
			Code:   320,
			Msg:    "unsupported oper_type",
			Result: 2,
		}, nil
	}
	return &types.UserModifyRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}
