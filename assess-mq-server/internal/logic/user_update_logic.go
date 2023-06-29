package logic

import (
	"context"
	"errors"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserUpdate 配置：用户级别修改
func (l *UserUpdateLogic) UserUpdate(in *proto.UserModifyReq) (*proto.UserModifyRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in UserUpdate, get req:", in)
	if len(in.GetList()) <= 0 {
		l.Logger.Error("no data update")
		return &proto.UserModifyRsp{
			Code:   220,
			Msg:    errors.New("no data").Error(),
			Result: 2,
		}, nil
	}
	if in.GetOperType() == 1 {
		for _, v := range in.GetList() {
			if v.GetUserId() == "" {
				l.Logger.Error("user_id empty...")
				continue
			}
			if err := l.svcCtx.AccountRepo.AddUser(l.ctx, v.GetUserId(), v.GetUserName(), v.GetGrade()); err != nil {
				l.Logger.Error("error add user :", err)
				continue
			}
		}
	} else if in.GetOperType() == 2 {
		for _, v := range in.GetList() {
			if v.GetUserId() == "" {
				l.Logger.Error("user_id empty...")
				continue
			}
			if err := l.svcCtx.AccountRepo.ModifyUserProperty(l.ctx, v.GetUserId(), v.GetGrade()); err != nil {
				l.Logger.Error("error modify user property:", err)
				continue
			}
			//先暂时不用更新全局数据，因为用不到
		}
	} else if in.GetOperType() == 3 {
		for _, v := range in.GetList() {
			if v.GetUserId() == "" {
				l.Logger.Error("user_id empty...")
				continue
			}
			if err := l.svcCtx.AccountRepo.DelUser(l.ctx, v.GetUserId()); err != nil {
				l.Logger.Error("error del user:", err)
				continue
			}
		}
	} else {
		l.Logger.Info("unsupported oper_type....")
	}
	return &proto.UserModifyRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}
