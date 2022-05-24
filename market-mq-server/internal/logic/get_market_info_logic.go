package logic

import (
	"algo_assess/global"
	"algo_assess/market-mq-server/internal/svc"
	"algo_assess/market-mq-server/proto"
	"algo_assess/pkg/tools"
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMarketInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMarketInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMarketInfoLogic {
	return &GetMarketInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  获取行情价格数量信息
func (l *GetMarketInfoLogic) GetMarketInfo(in *proto.MarkReq) (*proto.MarkRsp, error) {
	l.Infof("get req:%+v", in)
	// 行情的数据源不一样，所以这里要区分一下
	t := cast.ToString(in.GetStartTime())
	var hashKey string
	if in.GetSecSource() == 0 || in.GetSecSource() == 102 {
		hashKey = fmt.Sprintf("sz:%s:%s", in.GetSecId(), t[:8])
	} else if in.GetSecSource() == 103 {
		hashKey = fmt.Sprintf("sh:%s:%s", in.GetSecId(), t[:8])
	}
	infos := make(map[int64]*proto.LevelInfo)
	if true {

		l.Info("get hashKey:", hashKey)
		m, err := l.svcCtx.RedisClient.Hgetall(hashKey)
		if err != nil {
			l.Error("hget all fail:", err)
		}
		// 对 value 反序列化
		for k, v := range m {
			var out global.QuoteLevel2Data
			if err := json.Unmarshal(tools.String2Bytes(v), &out); err != nil {
				l.Error("unnarshal error:", err)
			}
			info := &proto.LevelInfo{
				LastPrice:  out.LastPrice,
				TradeVol:   out.TotalTradeVol,
				AskPrice:   out.AskPrice,
				AskVol:     out.AskVol,
				BidPrice:   out.BidPrice,
				BidVol:     out.BidVol,
				MarketVwap: out.Vwap,
			}
			orgiTime := cast.ToInt64(k)
			infos[orgiTime] = info
		}

	} else {
		data, result := l.svcCtx.MarketLevelRepo.GetMarketLevelByTime(l.ctx, hashKey, in.GetStartTime(), in.GetEndTime())
		if result.Error != nil {
			l.Logger.Error("get market level error :", result.Error)
			return nil, result.Error
		}
		l.Logger.Info("get market level result num:", result.RowsAffected)

		for _, v := range data {
			info := &proto.LevelInfo{
				LastPrice:  v.LastPrice,
				TradeVol:   v.TotalTradeVol,
				AskPrice:   v.AskPrice,
				AskVol:     v.AskVol,
				BidPrice:   v.BidPrice,
				BidVol:     v.BidVol,
				MarketVwap: v.MkVwap,
			}
			infos[v.OrgiTime] = info
		}
	}
	l.Info("get market rsp len:", len(infos))
	rsp := &proto.MarkRsp{
		Code:  0,
		Msg:   "success",
		Attrs: infos,
	}
	//l.Logger.Infof("get rsp:%+v", rsp)
	return rsp, nil
}
