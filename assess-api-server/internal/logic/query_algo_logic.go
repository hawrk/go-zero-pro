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

type QueryAlgoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryAlgoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryAlgoLogic {
	return &QueryAlgoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// QueryAlgo 查询数据修复母单数据
func (l *QueryAlgoLogic) QueryAlgo(req *types.ReqQueryAlgoOrder) (resp *types.RespQueryAlgoOrder, err error) {
	l.Logger.Infof("in QueryAlgo, get req:%+v", *req)
	var start, end int32
	if req.StartTime == 0 || req.EndTime == 0 {
		start = 0
		end = 0
	} else {
		start = cast.ToInt32(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
		end = cast.ToInt32(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))
	}
	in := &assessservice.ReqQueryAlgoOrder{
		AlgoId:    req.AlgoId,
		SecId:     req.SecId,
		PageId:    req.PageId,
		PageNum:   req.PageNum,
		AlgoName:  req.AlgoName,
		Scene:     int32(req.Scene),
		UserId:    req.UserId,
		UserType:  int32(req.UserType),
		StartTime: start,
		EndTime:   end,
	}
	o, err := l.svcCtx.AssessClient.QueryAlgoOrder(l.ctx, in)
	if err != nil {
		return nil, err
	}
	resp = &types.RespQueryAlgoOrder{
		Code:  int(o.GetCode()),
		Msg:   o.GetMsg(),
		Total: o.Total,
		Data: func(data []*assessservice.AlgoOrder) (ret []types.QueryAlgoOrder) {
			for _, d := range data {
				a := types.QueryAlgoOrder{
					Id:           d.GetId(),
					Date:         d.GetDate(),
					BatchNo:      d.GetBatchNo(),
					BatchName:    d.GetBatchName(),
					BasketId:     d.GetBasketId(),
					AlgoId:       d.GetAlgoId(),
					AlgorithmId:  d.GetAlgorithmId(),
					UserId:       d.GetUserId(),
					SecId:        d.GetSecId(),
					AlgoOrderQty: d.GetAlgoOrderQty(),
					TransTime:    d.GetTransTime(),
					StartTime:    d.GetStartTime(),
					EndTime:      d.GetEndTime(),
					FixFlag:      d.FixFlag,
					CreateTime:   d.GetCreateTime(),
					AlgoName:     d.GetAlgoName(),
				}
				ret = append(ret, a)
			}
			return ret
		}(o.GetParts()),
	}
	return
}
