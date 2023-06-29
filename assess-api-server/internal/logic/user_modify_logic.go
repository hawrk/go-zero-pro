package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"
	"errors"

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

func (l *UserModifyLogic) UserModify(req *types.ModifyUserReq) (resp *types.ModifyUserRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in UserModify, get req:%+v", req)
	if err := l.CheckParams(req); err != nil {
		return &types.ModifyUserRsp{
			Code:   250,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	var ls []*mqservice.UserUpdate
	for _, v := range req.Lists {
		l := &mqservice.UserUpdate{
			UserId:   v.UserId,
			UserName: v.UserName,
			Grade:    v.Grade,
		}
		ls = append(ls, l)
	}
	rsp, err := l.svcCtx.AssessMQClient.UserUpdate(l.ctx, &mqservice.UserModifyReq{
		OperType: req.OperType,
		List:     ls,
	})
	if err != nil {
		l.Logger.Error("rpc call UserUpdate error:", err)
		return &types.ModifyUserRsp{
			Code:   220,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	return &types.ModifyUserRsp{
		Code:   200,
		Msg:    rsp.GetMsg(),
		Result: rsp.GetResult(),
	}, nil
}

func (l *UserModifyLogic) CheckParams(req *types.ModifyUserReq) error {
	if req.OperType == 1 {
		for _, v := range req.Lists {
			if v.UserId == "" || v.UserName == "" || v.Grade == "" {
				return errors.New("field user_id, user_name, grade not set")
			}
		}
	} else if req.OperType == 2 {
		for _, v := range req.Lists {
			if v.Grade == "" {
				return errors.New("field grade not set")
			}
		}
	}
	return nil
}
