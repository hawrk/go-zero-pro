// Package consumer
/*
 Author: hawrkchen
 Date: 2022/11/28 10:47
 Desc: 深市行情修复
*/
package consumer

import (
	"algo_assess/global"
	"algo_assess/market-mq-server/internal/svc"
	pb "algo_assess/market-mq-server/proto"
	"algo_assess/models"
	"context"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type PerfFixSZMarketInfo struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPerfFixSZMarketInfo(ctx context.Context, svcCtx *svc.ServiceContext) *PerfFixSZMarketInfo {
	return &PerfFixSZMarketInfo{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *PerfFixSZMarketInfo) Consume(_, val string) error {

	// 行情数据修复，在原来行情的基础上进行修复，不独立一份数据
	// 与正常流程的二进制不同，修复数据直接解pb格式数据
	data := pb.QuoteLevel{}
	if err := proto.Unmarshal([]byte(val), &data); err != nil {
		s.Logger.Error("Unmarshal data fail:", err)
		return nil
	}
	s.Logger.Infof("get unmarshal: %+v", data)
	// 过滤掉期货和债券
	pre := data.SeculityId[:1]
	if pre == "1" || pre == "2" {
		//s.Logger.Info("unnecessary secID.... ")
		return nil
	}
	szQuote := TransQuoteWithPB(data, 1)
	s.Logger.Infof("get szQuote:%+v", szQuote)
	// 行情数据修复正常redis的数据都已经落到DB了，只能查DB的数据
	out, err := s.svcCtx.MarketLevelRepo.QuerySzMarketLevel(s.ctx, szQuote.SecID, szQuote.OrigTime)
	if err != nil {
		s.Logger.Error("QuerySzMarketLevel error:", err)
		return nil
	}
	if len(out) > 0 { // 记录已存在，不再处理
		s.Logger.Info("record has exist...")
		return nil
	}
	// 记录不存在，开始计算
	// 计算净成交量时，需要从DB中取出上一个时间点的成交量记录
	quote, err := s.svcCtx.MarketLevelRepo.QueryAllDataBySecId(s.ctx, szQuote.SecID, szQuote.OrigTime)
	if err != nil {
		s.Logger.Error("QueryAllDataBySecId error:", err)
		return nil
	}
	preTradeVol := GetSzPreTradeVol(quote, szQuote.OrigTime)
	netTotalVol := CalculateFixMarketVwap(&szQuote, preTradeVol)
	//保存到内存中供绩效计算
	if err := BuildLevel2Data(s.svcCtx.RedisClient, netTotalVol, &szQuote); err != nil {
		s.Logger.Error("build level2 data fail:", err)
		return nil
	}
	// 回写一次DB表  --直接插入
	if err := s.svcCtx.MarketLevelRepo.CreateMarketLevel(s.ctx, &szQuote); err != nil {
		s.Logger.Error("insert db market level fail:", err)
	}

	return nil
}

// GetSzPreTradeVol 从已排序的列表找到上一条记录
// originTime为当前修复的数据时间点， out: 输出上一个时间点的成交数量
func GetSzPreTradeVol(in []*models.TbSzQuoteLevel, originTime int64) int64 {
	for k, v := range in {
		if v.OrgiTime < originTime {
			continue
		}
		if k == 0 { // 没有上一个时间点的数据，直接置0
			return 0
		} else {
			return in[k-1].TotalTradeVol
		}
	}
	return 0
}

func CalculateFixMarketVwap(data *global.QuoteLevel2Data, preVol int64) int64 {
	// 累计数量计算vwap,就是算一天的总量,不是一分钟的总量
	//logx.Infof("get data:%+v", data)
	dayKey := "fix:" + data.SecID + ":" + cast.ToString(data.OrigTime)[:8]
	//logx.Info("get dayKey:", dayKey)
	// 计算市场vwap
	global.GlobalMarketVwap.RWMutex.Lock()
	defer global.GlobalMarketVwap.RWMutex.Unlock()

	//logx.Info("debug core dump: lastVol:", lastVol)
	var lastVol int64
	if global.GlobalMarketVwap.LastVol[dayKey] == 0 { // 该笔是修复数据的第一笔时，取已取得的上一笔的成交总量
		lastVol = preVol
	} else { // 非第一笔时，直接取修复数据的上一笔
		lastVol = global.GlobalMarketVwap.LastVol[dayKey] // 取上一个成交总量
	}
	curVol := data.TotalTradeVol - lastVol // 当前时间点净成交量 =  当前成交总量 - 上一个成交总量
	//logx.Info("curVol:", curVol)

	global.GlobalMarketVwap.TotalPrxCal[dayKey] += curVol * data.LastPrice // ∑(订单成交数量 *成交价格)
	global.GlobalMarketVwap.TotalVol[dayKey] += curVol                     // ∑订单成交数量
	if global.GlobalMarketVwap.TotalPrxCal[dayKey] <= 0 || global.GlobalMarketVwap.TotalVol[dayKey] <= 0 {
		data.Vwap = 0.0000
	} else {
		fVwap := (float64(global.GlobalMarketVwap.TotalPrxCal[dayKey]) / float64(global.GlobalMarketVwap.TotalVol[dayKey])) / 100
		data.Vwap, _ = decimal.NewFromFloat(fVwap).Truncate(4).Float64()
	}
	global.GlobalMarketVwap.MVwap[dayKey] = data.Vwap
	global.GlobalMarketVwap.LastVol[dayKey] = data.TotalTradeVol // 填充当前成交总量
	//logx.Infof("after count get market vwap :%+v", global.GlobalMarketVwap)
	return curVol // 返回当前时间点的净增量，推给下游
}
