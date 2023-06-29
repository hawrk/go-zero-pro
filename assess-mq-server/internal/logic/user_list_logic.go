package logic

import (
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserList 配置：用户列表
func (l *UserListLogic) UserList(in *proto.UserListReq) (*proto.UserListRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in SecurityList, get req:", in)
	var lists []*proto.UserInfo
	out, count, err := l.svcCtx.AccountRepo.GetAccountLists(l.ctx, in.GetUserId(), int(in.GetPage()), int(in.GetLimit()))
	if err != nil {
		l.Logger.Error("GetAccountLists error:", err)
		return &proto.UserListRsp{
			Code:  210,
			Msg:   err.Error(),
			Total: 0,
			Infos: nil,
		}, nil
	}
	for _, v := range out {
		s := &proto.UserInfo{
			Id:         int64(v.AccountId),
			UserId:     v.UserId,
			UserName:   v.UserName,
			UserType:   int32(v.UserType),
			UserGrade:  v.UserGrade,
			UpdateTime: GetUpdateTime(v.UpdateTime, v.CreateTime),
		}
		lists = append(lists, s)
	}
	return &proto.UserListRsp{
		Code:  200,
		Msg:   "success",
		Total: count,
		Infos: lists,
	}, nil
}
