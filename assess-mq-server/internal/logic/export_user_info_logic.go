package logic

import (
	"algo_assess/pkg/tools"
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExportUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportUserInfoLogic {
	return &ExportUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ExportUserInfo 配置： 用户信息导出
func (l *ExportUserInfoLogic) ExportUserInfo(in *proto.ExportUserReq) (*proto.ExportUserRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in ExportUserInfo, get req:", in)
	out, err := l.svcCtx.AccountRepo.GetAccountInfos(l.ctx)
	if err != nil {
		l.Logger.Error("GetAccountInfos error :", err)
		return &proto.ExportUserRsp{
			Infos: nil,
		}, nil
	}

	var lists []*proto.UserInfo
	for _, v := range out {
		l := &proto.UserInfo{
			UserId:     v.UserId,
			UserName:   v.UserName,
			UserType:   int32(v.UserType),
			UserGrade:  v.UserGrade,
			UpdateTime: tools.Time2String(v.UpdateTime),
		}
		lists = append(lists, l)
	}

	return &proto.ExportUserRsp{
		Infos: lists,
	}, nil
}
