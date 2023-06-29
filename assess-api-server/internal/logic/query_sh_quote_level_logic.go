package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryShQuoteLevelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryShQuoteLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryShQuoteLevelLogic {
	return &QueryShQuoteLevelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryShQuoteLevelLogic) QueryShQuoteLevel(req *types.ReqShQuoteLevel) (resp *types.RespShQuoteLevel, err error) {
	l.Logger.Infof("in QueryShQuoteLevel get req:%+v", *req)
	in := &assessservice.ReqQueryQuoteLevel{
		SecurityId: req.SecurityId,
		PageId:     req.PageId,
		PageNum:    req.PageNum,
		MaxId:      req.MaxId,
	}
	o, err := l.svcCtx.AssessClient.QueryShQuoteLevel(l.ctx, in)
	if err != nil {
		return nil, err
	}

	resp = &types.RespShQuoteLevel{
		Code:  int(o.GetCode()),
		Msg:   o.GetMsg(),
		Total: o.Total,
		MaxId: o.GetMaxId(),
		Data: func(data []*assessservice.QuoteLevel) (ret []types.QuoteLevel) {
			for _, d := range data {
				a := types.QuoteLevel{
					Id:            d.Id,
					SeculityId:    d.SeculityId[3:],
					OrgiTime:      d.OrgiTime,
					LastPrice:     float64(d.LastPrice) / 100,
					AskPrice:      d.AskPrice,
					AskVol:        d.AskVol,
					BidPrice:      d.BidPrice,
					BidVol:        d.BidVol,
					TotalTradeVol: d.TotalTradeVol,
					TotalAskVol:   d.TotalAskVol,
					TotalBidVol:   d.TotalBidVol,
					MkVwap:        float64(d.MkVwap),
					FixFlag:       int(d.FixFlag),
				}
				ret = append(ret, a)
			}
			return ret
		}(o.GetParts()),
	}
	return
}
