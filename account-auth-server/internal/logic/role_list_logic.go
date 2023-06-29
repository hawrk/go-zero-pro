package logic

import (
	"account-auth/account-auth-server/models"
	"account-auth/account-auth-server/pkg/tools"
	"context"
	"encoding/json"

	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoleListLogic {
	return &RoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RoleListLogic) RoleList(req *types.RoleListReq) (resp *types.RoleListRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in RoleList get req:%+v", *req)
	var list []*models.TbAuthRole
	id := (req.Page - 1) * req.Limit
	var count int64
	if req.Scene == 1 {
		list, count, err = l.svcCtx.AuthRoleRepo.GetRoleList(l.ctx, req.RoleId, req.RoleName, req.ChanType, int(req.Page), int(req.Limit))
		if err != nil {
			l.Logger.Error("GetRoleList err:", err)
			return &types.RoleListRsp{
				Code:  310,
				Msg:   err.Error(),
				Total: 0,
				List:  nil,
			}, nil
		}
	} else if req.Scene == 2 { // 只拉取生效的，并且不用分页
		list, err = l.svcCtx.AuthRoleRepo.GetEffectRoleList(l.ctx, req.ChanType)
		if err != nil {
			l.Logger.Error("GetRoleList err:", err)
			return &types.RoleListRsp{
				Code:  310,
				Msg:   err.Error(),
				Total: 0,
				List:  nil,
			}, nil
		}
	} else {
		l.Logger.Error(" invalid, scene:", req.Scene)
	}

	var ls []types.RoleInfos
	for _, v := range list {
		id++
		r := types.RoleInfos{
			Id:         int64(id),
			RoleId:     int(v.RoleId),
			RoleName:   v.RoleName,
			RoleAuth:   v.RoleAuth,
			RoleDesc:   l.ParseRoleAuth(v.RoleAuth),
			Stutus:     v.Status,
			CreateTime: tools.Time2String(v.CreateTime),
		}
		ls = append(ls, r)
	}

	return &types.RoleListRsp{
		Code:  200,
		Msg:   "success",
		Total: count,
		List:  ls,
	}, nil
}

// ParseRoleAuth 解析json串，把节点都提取出来
func (l *RoleListLogic) ParseRoleAuth(js string) string {
	if len(js) <= 0 {
		return ""
	}
	// 反序列化
	var roleAuth []FirstMenu
	if err := json.Unmarshal([]byte(js), &roleAuth); err != nil {
		l.Logger.Error("json Unmarshal fail:", err)
		return ""
	}
	// 提取菜单节点
	//l.Logger.Infof("get roleAuth:%+v", roleAuth)
	var listDesc string
	for _, v := range roleAuth {
		if v.Auth == 1 {
			listDesc += v.OneLevelName + "/"
		}
		//for _, v2 := range v.Cmpt {
		//	listDesc += v2.Name + "/"
		//}
		for _, v3 := range v.SecondM {
			if v3.Auth == 1 {
				listDesc += v3.TwoLevelName + "/"
			}
		}
		if len(listDesc) > 0 {
			listDesc = listDesc[:len(listDesc)-1] + ";"
		}
	}
	//l.Logger.Info("get list Desc:", listDesc)

	return listDesc
}
