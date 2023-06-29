package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryChildLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryChildLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryChildLogic {
	return &QueryChildLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// QueryChild 数据修复查询子单数据
func (l *QueryChildLogic) QueryChild(req *types.ReqQueryChildOrder) (resp *types.RespQueryChildOrder, err error) {
	l.Logger.Infof("in QueryChild, get req:%+v", *req)
	var start, end int32
	if req.StartTime == 0 || req.EndTime == 0 {
		start = 0
		end = 0
	} else {
		start = cast.ToInt32(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
		end = cast.ToInt32(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))
	}
	in := &assessservice.ReqQueryChildOrder{
		UserId:       req.UserId,
		SecurityId:   req.SecurityId,
		ChildOrderId: req.ChildOrderId,
		AlgoName:     req.AlgoName,
		PageId:       req.PageId,
		PageNum:      req.PageNum,
		Scene:        int32(req.Scene),
		UserType:     int32(req.UserType),
		StartTime:    start,
		EndTime:      end,
		AlgoOrderId:  req.AlgoOrderId,
	}
	o, err := l.svcCtx.AssessClient.QueryChildOrder(l.ctx, in)
	if err != nil {
		return nil, err
	}
	resp = &types.RespQueryChildOrder{
		Code:  int(o.GetCode()),
		Msg:   o.GetMsg(),
		Total: o.Total,
		Data: func(data []*assessservice.ChildOrder) (ret []types.ChildOrderInfo) {
			for _, d := range data {
				a := types.ChildOrderInfo{
					Id:            d.Id,
					Date:          d.GetDate(),
					BatchNo:       d.GetBatchNo(),
					BatchName:     d.GetBatchName(),
					ChildOrderId:  d.ChildOrderId,
					AlgoOrderId:   d.AlgoOrderId,
					AlgorithmType: d.AlgorithmType,
					AlgorithmId:   d.AlgorithmId,
					UserId:        d.UserId,
					UsecurityId:   d.UsecurityId,
					SecurityId:    d.SecurityId,
					TradeSide:     int8(d.TradeSide),
					OrderQty:      d.OrderQty,
					Price:         float64(d.Price) / 10000,
					OrderType:     uint16(d.OrderType),
					LastPx:        float64(d.LastPx) / 10000,
					LastQty:       d.LastQty,
					ComQty:        d.ComQty,
					ArrivedPrice:  float64(d.ArrivedPrice) / 10000,
					TotalFee:      d.TotalFee,
					OrdStatus:     uint16(d.OrdStatus),
					TransactTime:  d.TransactTime,
					TransactAt:    d.TransactAt,
					ProcStatus:    uint16(d.ProcStatus),
					FixFlag:       d.FixFlag,
					CreateTime:    d.CreateTime,
					AlgoName:      d.GetAlgoName(),
				}
				ret = append(ret, a)
			}
			return ret
		}(o.GetParts()),
	}
	return
}
