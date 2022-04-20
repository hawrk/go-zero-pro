package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/assess-rpc-server/assessservice"
	"context"
	"github.com/jinzhu/copier"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GeneralLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeneralLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeneralLogic {
	return &GeneralLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeneralLogic) General(req *types.GeneralReq) (resp *types.GeneralRsp, err error) {
	l.Logger.Info("into General,req:", req)

	gRspCh := make(chan *assessservice.GeneralRsp)
	mRspCh := make(chan *mqservice.GeneralRsp)
	go func() {
		var gReq assessservice.GeneralReq
		copier.Copy(&gReq, req)
		rsp, err := l.svcCtx.AssessClient.GetGeneral(l.ctx, &gReq)
		if err != nil {
			l.Logger.Error("assess rpc call error :", err)
			gRspCh <- &assessservice.GeneralRsp{}
			return
		}
		gRspCh <- rsp
	}()

	go func() {
		var mReq mqservice.GeneralReq
		copier.Copy(&mReq, req)
		rsp, err := l.svcCtx.AssessMQClient.GetMqGeneral(l.ctx, &mReq)
		if err != nil {
			l.Logger.Error("assess mq rpc call error:", err)
			mRspCh <- &mqservice.GeneralRsp{}
			return
		}
		mRspCh <- rsp
	}()

	gRsp, mRsp := <-gRspCh, <-mRspCh

	var data []types.GeneralData

	for _, v := range mRsp.Info {
		var detail types.GeneralData
		_ = copier.Copy(&detail, v)
		data = append(data, detail)
	}
	for _, v := range gRsp.Info {
		var detail types.GeneralData
		_ = copier.Copy(&detail, v)
		data = append(data, detail)
	}
	l.Logger.Infof("get rsp: %+v", data)
	p := &types.GeneralRsp{
		Code: 200,
		Msg:  "success",
		Data: data,
	}

	return p, nil
}
