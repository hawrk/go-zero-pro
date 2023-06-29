package logic

import (
	"algo_assess/pkg/tools"
	"context"
	"time"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type SecurityListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSecurityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SecurityListLogic {
	return &SecurityListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SecurityList  配置:证券列表
func (l *SecurityListLogic) SecurityList(in *proto.SecurityListReq) (*proto.SecurityListRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("SecurityList, get req:", in)
	var infos []*proto.SecurityInfo
	out, count, err := l.svcCtx.SecurityRepo.GetSecurityList(l.ctx, in.GetSecId(), int(in.GetPage()), int(in.GetLimit()))
	if err != nil {
		l.Logger.Error("GetSecurityList error :", err)
		return &proto.SecurityListRsp{
			Code:  210,
			Msg:   err.Error(),
			Total: 0,
			Infos: nil,
		}, nil
	}

	for _, v := range out {
		s := &proto.SecurityInfo{
			Id:         v.Id,
			SecId:      v.SecurityId,
			SecName:    tools.RMu0000(v.SecurityName),
			Status:     int32(v.Status),
			FundType:   int32(v.FundType),
			StockType:  int32(v.StockType),
			Liquidity:  int32(v.Liquidity),
			Industry:   v.Industry,
			UpdateTime: GetUpdateTime(v.UpdateTime, v.CreateTime),
		}
		infos = append(infos, s)
	}

	return &proto.SecurityListRsp{
		Code:  200,
		Msg:   "success",
		Total: count,
		Infos: infos,
	}, nil
}

// GetUpdateTime 如果更新时间为空，则用创建时间填充
func GetUpdateTime(ut, ct time.Time) string {
	if ut.IsZero() {
		return tools.Time2String(ct)
	}
	return tools.Time2String(ut)
}
