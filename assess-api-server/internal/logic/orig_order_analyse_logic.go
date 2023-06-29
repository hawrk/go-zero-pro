package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrigOrderAnalyseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrigOrderAnalyseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrigOrderAnalyseLogic {
	return &OrigOrderAnalyseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrigOrderAnalyseLogic) OrigOrderAnalyse(req *types.OrigAnalyseReq) (resp *types.OrigAnalyseResp, err error) {
	l.Logger.Infof("OrigAlgoAnalyseReq:%+v", req)
	var lists []*assessservice.Analyse
	for _, v := range req.Orders {
		l := &assessservice.Analyse{
			Date: v.Date,
			Id:   v.Id,
		}
		lists = append(lists, l)
	}
	a, err := l.svcCtx.AssessClient.OrigOrderAnalyse(l.ctx, &assessservice.OrigAnalyseReq{
		Orders:    lists,
		OrderType: req.OrderType,
	})
	if err != nil {
		return
	}
	return &types.OrigAnalyseResp{
		Code:      int(a.Code),
		Msg:       a.Msg,
		BatchNo:   a.BatchNo,
		StartTime: a.GetStartTime(),
		EndTime:   a.GetEndTime(),
	}, nil
}
