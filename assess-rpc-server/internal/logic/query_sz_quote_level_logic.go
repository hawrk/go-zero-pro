package logic

import (
	"algo_assess/models"
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuerySzQuoteLevelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQuerySzQuoteLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerySzQuoteLevelLogic {
	return &QuerySzQuoteLevelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  深市行情信息推送
func (l *QuerySzQuoteLevelLogic) QuerySzQuoteLevel(in *proto.ReqQueryQuoteLevel) (*proto.RespQueryQuoteLevel, error) {
	// 指定具体证券代码时，直接查询数据和总条数即可
	// 如果没有指定证券代码，查询数据和总条数需要分开查询，否则会有性能问题
	var count int64
	var markets []*models.TbSzQuoteLevel
	var err error
	if in.GetSecurityId() != "" {
		count, markets, err = l.svcCtx.MarketLevelRepo.GetSzMarketLevelByPageWithCount(l.ctx, in)
		if err != nil {
			l.Logger.Error("GetSzMarketLevelByPageWithCount error:", err)
			return &proto.RespQueryQuoteLevel{
				Code:  10000,
				Msg:   "查询失败",
				Total: count,
			}, err
		}
	} else { // 全量查询要小心性能问题
		count = l.svcCtx.MarketLevelRepo.GetSzMarketLevelCount(l.ctx)
		markets, err = l.svcCtx.MarketLevelRepo.GetSzMarketLevelByPageWithoutCount(l.ctx, in)
		if err != nil {
			l.Logger.Error("GetSzMarketLevelByPageWithoutCount error:", err)
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
