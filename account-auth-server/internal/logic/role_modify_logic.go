package logic

import (
	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleModifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleModifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleModifyLogic {
	return &RoleModifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleModifyLogic) RoleModify(req *types.RoleModifyReq) (resp *types.RoleModifyRsp, err error) {
	l.Logger.Infof("in RoleModify, get req:%+v", *req)
	if req.OperType == 1 { // 新增
		// 新增时先检查role_id 和role_name是否已存在
		if req.RoleId == 0 || req.RoleName == "" {
			l.Logger.Error("field role_id/role_name not set")
			return &types.RoleModifyRsp{
				Code:   320,
				Msg:    "角色ID或角色名称为空",
				Result: 2,
			}, nil
		}
		id, _, err := l.svcCtx.AuthRoleRepo.CheckExistRole(l.ctx, req.ChanType, int(req.RoleId), "")
		if err != nil {
			l.Logger.Error("CheckExistRole error:", err)
		}
		if id > 0 {
			l.Logger.Error("add role : role_id exist.")
			return &types.RoleModifyRsp{
				Code:   320,
				Msg:    "新增角色失败,原因:[角色ID已存在]",
				Result: 2,
			}, nil
		}
		_, name, err := l.svcCtx.AuthRoleRepo.CheckExistRole(l.ctx, req.ChanType, 0, req.RoleName)
		if err != nil {
			l.Logger.Error("CheckExistRole error:", err)
		}
		if name != "" {
			l.Logger.Error("add role : role_name exist.")
			return &types.RoleModifyRsp{
				Code:   320,
				Msg:    "新增角色失败,原因:[角色名称已存在]",
				Result: 2,
			}, nil
		}
		if err := l.svcCtx.AuthRoleRepo.CreateAuthRole(l.ctx, req); err != nil {
			l.Logger.Error("create auth role error:", err)
			return &types.RoleModifyRsp{
				Code:   320,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}
	} else if req.OperType == 2 { // 修改
		if err := l.svcCtx.AuthRoleRepo.UpdateAuthRole(l.ctx, req); err != nil {
			l.Logger.Error("update auth role error:", err)
			return &types.RoleModifyRsp{
				Code:   320,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}

	} else if req.OperType == 3 { // 删除
		// 如果该角色有绑定了用户，则不可删除
		out, err := l.svcCtx.AuthUserRepo.GetUserByRoleId(l.ctx, int(req.RoleId), req.ChanType)
		if err != nil {
			l.Logger.Error("GetUserByRoleId error:", err)
			return &types.RoleModifyRsp{
				Code:   360,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}
		if len(out) > 0 {
			l.Logger.Info("get user_id by roleId count:", len(out))
			return &types.RoleModifyRsp{
				Code:   360,
				Msg:    "已有用户绑定了该角色，不可删除",
				Result: 2,
			}, nil
		}
		// 可以正常删除了
		if err := l.svcCtx.AuthRoleRepo.DelAuthRole(l.ctx, int(req.RoleId)); err != nil {
			l.Logger.Error("del auth role error:", err)
			return &types.RoleModifyRsp{
				Code:   320,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}
	} else {
		l.Logger.Error("unsupported oper_type:", req.OperType)
		return &types.RoleModifyRsp{
			Code:   320,
			Msg:    "unsupported oper_type",
			Result: 2,
		}, nil
	}

	return &types.RoleModifyRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}
