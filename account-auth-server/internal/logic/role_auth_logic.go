package logic

import (
	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"
	"context"
	//"encoding/json"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleAuthLogic {
	return &RoleAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// RoleAuth 拉取角色权限，在权限新增或修改时拉取
func (l *RoleAuthLogic) RoleAuth(req *types.RoleAuthReq) (resp *types.RoleAuthRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in RoleAuth, get req:%+v", *req)
	// 普通用户
	out, err := l.svcCtx.AuthMenuRepo.GetDefaultAuthMenu(l.ctx, req.ChanType)
	if err != nil {
		l.Logger.Error("GetDefaultAuthMenu error:", err)
		return &types.RoleAuthRsp{
			Code:     330,
			Msg:      err.Error(),
			RoleAuth: "",
		}, nil
	}
	//l.Logger.Infof("get out:%+v", out)
	authJson, err := BuildMenuList(0, out)
	if err != nil {
		l.Logger.Error("build auth list error:", err)
	}
	return &types.RoleAuthRsp{
		Code:     200,
		Msg:      "success",
		RoleAuth: authJson,
	}, nil
}
