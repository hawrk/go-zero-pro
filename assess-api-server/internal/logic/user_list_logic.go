package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.ListUserReq) (resp *types.ListUserRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in UserList, get req:%+v", req)
	uReq := &mqservice.UserListReq{
		UserId: req.UserId,
		Page:   req.Page,
		Limit:  req.Limit,
	}
	rsp, err := l.svcCtx.AssessMQClient.UserList(l.ctx, uReq)
	if err != nil {
		l.Logger.Error("rpc call UserList error:", err)
		return &types.ListUserRsp{
			Code:  210,
			Msg:   err.Error(),
			Total: 0,
			Infos: nil,
		}, nil
	}

	var list []types.UserInfo
	for _, v := range rsp.GetInfos() {
		l := types.UserInfo{
			Id:         v.GetId(),
			UserId:     v.GetUserId(),
			UserName:   v.GetUserName(),
			UserType:   v.GetUserType(),
			UserGrade:  v.GetUserGrade(),
			UpdateTime: v.GetUpdateTime(),
		}
		list = append(list, l)
	}

	return &types.ListUserRsp{
		Code:  200,
		Msg:   "success",
		Total: rsp.GetTotal(),
		Infos: list,
	}, nil
}
