package logic

import (
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"
	"algo_assess/global"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMqGeneralLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMqGeneralLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMqGeneralLogic {
	return &GetMqGeneralLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取绩效概况
func (l *GetMqGeneralLogic) GetMqGeneral(in *proto.GeneralReq) (*proto.GeneralRsp, error) {
	l.Logger.Info("into mq GetGeneral:req:", in)
	// 读取本地缓存数据
	global.GlobalAssess.RWMutex.RLock()

	infos := make([]*proto.AssessInfo, 0, len(global.GlobalAssess.CalAlgo))
	for _, val := range global.GlobalAssess.CalAlgo {
		if val.SecurityId != in.GetSecId() { // 非该证券的不取
			continue
		}
		if val.AlgorithmId != int(in.GetAlgoId()) { // 非该算法的不取
			continue
		}
		if val.UserId != in.GetUserId() { // 非该用户的不取
			continue
		}
		info := &proto.AssessInfo{
			TransactTime:          val.TransactAt,
			OrderQty:              val.OrderQty,
			LastQty:               val.LastQty,
			CancelledQty:          val.CancelQty,
			RejectedQty:           val.RejectedQty,
			Vwap:                  val.Vwap,
			VwapDeviation:         val.VwapDeviation,
			LastPrice:             val.LastPrice,
			ArrivedPrice:          val.ArrivedPrice,
			ArrivedPriceDeviation: val.ArrivedPriceDeviation,
			MarketRate:            val.MarketRate,
			DealRate:              val.DealRate,
			DealProgress:          val.DealProgress,
		}
		infos = append(infos, info)
	}
	global.GlobalAssess.RWMutex.RUnlock()
	rsp := &proto.GeneralRsp{
		Code: 0,
		Msg:  "success",
		Info: infos,
	}
	l.Logger.Infof("get cache assess, rsp:%+v", infos)
	return rsp, nil
}
