package logic

import (
	"account-auth/account-auth-server/pkg/tools"
	"context"

	"account-auth/account-auth-server/internal/svc"
	"account-auth/account-auth-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthUserListLogic {
	return &AuthUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthUserListLogic) AuthUserList(req *types.UserListReq) (resp *types.UserListRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in AuthUserList, get req:%+v", *req)
	id := (req.Page - 1) * req.Limit
	list, count, err := l.svcCtx.AuthUserRepo.GetUserList(l.ctx, req.UserName, req.ChanType, int(req.Page), int(req.Limit))
	if err != nil {
		l.Logger.Error("GetUserList error:", err)
		return &types.UserListRsp{
			Code:  310,
			Msg:   err.Error(),
			Total: 0,
			List:  nil,
		}, nil
	}
	var ls []types.UserInfos
	for _, v := range list {
		id++
		l := types.UserInfos{
			Id:         int64(id),
			UserId:     v.UserId,
			UserName:   v.UserName,
			RoleId:     int32(v.RoleId),
			RoleName:   v.RoleName,
			UserType:   int32(v.UserType),
			Status:     v.Status,
			CreateTime: tools.Time2String(v.CreateTime),
		}
		ls = append(ls, l)
	}
	return &types.UserListRsp{
		Code:  200,
		Msg:   "success",
		Total: count,
		List:  ls,
	}, nil
}
