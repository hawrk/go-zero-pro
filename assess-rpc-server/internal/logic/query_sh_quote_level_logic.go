package logic

import (
	"algo_assess/models"
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryShQuoteLevelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryShQuoteLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryShQuoteLevelLogic {
	return &QueryShQuoteLevelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  沪市行情信息
func (l *QueryShQuoteLevelLogic) QueryShQuoteLevel(in *proto.ReqQueryQuoteLevel) (*proto.RespQueryQuoteLevel, error) {
	l.Logger.Info("in QueryShQuoteLevel, get req:", in)
	var count int64
	var markets []*models.TbShQuoteLevel
	var err error
	if in.GetSecurityId() != "" {
		count, markets, err = l.svcCtx.SHMarketLevelRepo.GetShMarketLevelByPageWithCount(l.ctx, in)
		if err != nil {
			l.Logger.Error("GetShMarketLevelByPageWithCount error:", err)
			return &proto.RespQueryQuoteLevel{
				Code:  10000,
				Msg:   "查询失败",
				Total: count,
			}, err
		}
	} else { // 全量查询要小心性能问题
		count = l.svcCtx.SHMarketLevelRepo.GetShMarketLevelCount(l.ctx)
		markets, err = l.svcCtx.SHMarketLevelRepo.GetShMarketLevelByPageWithoutCount(l.ctx, in)
		if err != nil {
			l.Logger.Error("GetShMarketLevelByPageWithoutCount error:", err)
			return &proto.RespQueryQuoteLevel{
				Code:  10000,
				Msg:   "查询失败",
				Total: count,
			}, err
		}
	}

	id := (in.GetPageId() - 1) * in.GetPageNum()
	var lastSeqId int64
	var ret []*proto.QuoteLevel
	for _, m := range markets {
		id++
		lastSeqId = int64(m.Id)
		//l.Logger.Info("get m.id:", lastSeqId)
		a := &proto.QuoteLevel{
			Id:            int64(id),
			SeculityId:    m.SeculityId,
			OrgiTime:      m.OrgiTime,
			LastPrice:     m.LastPrice,
			AskPrice:      m.AskPrice,
			AskVol:        m.AskVol,
			BidPrice:      m.BidPrice,
			BidVol:        m.BidVol,
			TotalTradeVol: m.TotalTradeVol,
			TotalAskVol:   m.TotalAskVol,
			TotalBidVol:   m.TotalBidVol,
			MkVwap:        float32(m.MkVwap),
			FixFlag:       int32(m.FixFlag),
		}
		ret = append(ret, a)
	}
	return &proto.RespQueryQuoteLevel{
		Code:  200,
		Msg:   "查询成功",
		Total: count,
		Parts: ret,
		MaxId: lastSeqId,
	}, nil
}
