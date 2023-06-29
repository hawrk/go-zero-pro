package logic

import (
	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAuthLogic {
	return &UserAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserAuth 取用户登陆菜单权限
func (l *UserAuthLogic) UserAuth(req *types.UserAuthReq) (resp *types.UserAuthRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in UserAuth, get req:%+v", *req)
	// 查询用户权限表
	out, err := l.svcCtx.AuthUserRepo.GetAuthUserById(l.ctx, req.UserId, req.ChanType, 1)
	if err != nil {
		l.Logger.Error("GetAuthUserById error:", err)
		return &types.UserAuthRsp{
			Code: 370,
			Msg:  err.Error(),
			Auth: "",
		}, nil
	}
	//if req.UserType == 1 || out.RoleId == 0 { // 如果是超管，则返回所有权限; 如果是普通用户，则返回默认权限
	//	menu, err := l.svcCtx.AuthMenuRepo.GetDefaultAuthMenu(l.ctx, req.ChanType)
	//	if err != nil {
	//		l.Logger.Error("GetDefaultAuthMenu error:", err)
	//		return &types.UserAuthRsp{
	//			Code: 370,
	//			Msg:  err.Error(),
	//			Auth: "",
	//		}, nil
	//	}
	//	authJson, err := BuildMenuList(req.UserType, menu)
	//	if err != nil {
	//		l.Logger.Error("BuildMenuList error:", err)
	//	}
	//	return &types.UserAuthRsp{
	//		Code: 200,
	//		Msg:  "success",
	//		Auth: authJson,
	//	}, nil
	//}
	// 其他case,直接返回角色权限的列表
	auth, err := l.svcCtx.AuthRoleRepo.GetAuthByRoleId(l.ctx, out.RoleId)
	if err != nil {
		l.Logger.Error("GetAuthByRoleId, error:", err)
		return &types.UserAuthRsp{
			Code: 370,
			Msg:  err.Error(),
			Auth: "",
		}, nil
	}
	// 兼容逻辑，如果角色权限表是空表的话，至少要保证超级管理员是有数据返回的
	if auth.RoleAuth == "" {
		menu, err := l.svcCtx.AuthMenuRepo.GetDefaultAuthMenu(l.ctx, req.ChanType)
		if err != nil {
			l.Logger.Error("GetDefaultAuthMenu error:", err)
			return &types.UserAuthRsp{
				Code: 370,
				Msg:  err.Error(),
				Auth: "",
			}, nil
		}
		// 如果是超管，则返回所有权限; 如果是普通用户，则返回默认权限
		authJson, err := BuildMenuList(req.UserType, menu)
		if err != nil {
			l.Logger.Error("BuildMenuList error:", err)
		}
		//l.Logger.Info("get json:", authJson)
		return &types.UserAuthRsp{
			Code: 200,
			Msg:  "success",
			Auth: authJson,
		}, nil
	}

	//var ret MenuList
	////var unma FirstMenu
	//if err := json.Unmarshal([]byte(auth.RoleAuth), &ret.List); err != nil {
	//	l.Logger.Error("json unmarshal error:", err)
	//}
	//l.Logger.Info("get auth.RoleAuth:", auth.RoleAuth)
	////ret.List = append(ret.List, unma)
	//o, err := json.Marshal(ret)
	return &types.UserAuthRsp{
		Code: 200,
		Msg:  "success",
		Auth: auth.RoleAuth,
	}, nil
}
