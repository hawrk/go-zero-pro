package logic

import (
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGeneralLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGeneralLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGeneralLogic {
	return &GetGeneralLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取绩效概况
func (l *GetGeneralLogic) GetGeneral(in *proto.GeneralReq) (*proto.GeneralRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in GetGeneral, get Req:%+v", in)
	data, result := l.svcCtx.OrderAssessRepo.GetAlgoAssess(l.ctx, in)
	if result.Error != nil {
		l.Logger.Error("get assess error :", result.Error)
		return nil, result.Error
	}
	l.Logger.Info("get assess result num:", result.RowsAffected)
	infos := make([]*proto.AssessInfo, 0, len(data))
	for _, v := range data {
		info := &proto.AssessInfo{
			TransactTime:          v.TransactTime,
			OrderQty:              v.OrderQty,
			LastQty:               v.LastQty,
			CancelledQty:          v.CancelQty,
			RejectedQty:           v.RejectedQty,
			Vwap:                  v.Vwap,
			VwapDeviation:         v.VwapDeviation,
			LastPrice:             v.LastPrice,
			ArrivedPrice:          (v.ArrivedPrice) / 100,
			ArrivedPriceDeviation: v.ArrivedPriceDeviation,
			MarketRate:            v.MarketRate,
			DealRate:              v.DealRate,
			DealProgress:          v.DealProgress,
		}
		infos = append(infos, info)
	}
	rsp := &proto.GeneralRsp{
		Code: 200,
		Msg:  "SUCCESS",
		Info: infos,
	}
	return rsp, nil
}
