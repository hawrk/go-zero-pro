package logic

import (
	"algo_assess/pkg/tools"
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExportSecurityInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExportSecurityInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportSecurityInfoLogic {
	return &ExportSecurityInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ExportSecurityInfo 配置：证券信息导出
func (l *ExportSecurityInfoLogic) ExportSecurityInfo(in *proto.ExportSecurityReq) (*proto.ExportSecurityRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in ExportSecurityInfo, get req:", in)
	out, err := l.svcCtx.SecurityRepo.GetSecurityInfos(l.ctx)
	if err != nil {
		l.Logger.Error("GetSecurityInfos error:", err)
		return &proto.ExportSecurityRsp{
			Infos: nil,
		}, nil
	}
	var lists []*proto.SecurityInfo
	for _, v := range out {
		i := &proto.SecurityInfo{
			Id:         v.Id,
			SecId:      v.SecurityId,
			SecName:    v.SecurityName,
			Status:     int32(v.Status),
			FundType:   int32(v.FundType),
			StockType:  int32(v.StockType),
			Liquidity:  int32(v.Liquidity),
			Industry:   v.Industry,
			UpdateTime: tools.Time2String(v.UpdateTime),
		}
		lists = append(lists, i)
	}
	return &proto.ExportSecurityRsp{
		Infos: lists,
	}, nil
}
