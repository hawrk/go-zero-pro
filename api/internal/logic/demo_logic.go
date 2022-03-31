package logic

import (
	"algo_assess/rpc/assess/assess"
	"context"

	"algo_assess/api/internal/svc"
	"algo_assess/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DemoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDemoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DemoLogic {
	return &DemoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DemoLogic) Demo(req *types.DemoReq) (resp *types.DemoRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Info(" into Demo, req:", req.Hello)
	//p := &types.DemoRsp{Reply: "hello from api"}
	demoReq := &assess.DemoReq{
		Hello: req.Hello,
	}
	rsp, err := l.svcCtx.AssessClient.GetDemo(l.ctx, demoReq)
	if err != nil {
		l.Logger.Error("rpc assess error:", err)
		return nil, err
	}
	p := &types.DemoRsp{Reply: rsp.Reply}
	l.Logger.Info("get resp:", p)

	return p, nil
}
