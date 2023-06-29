// Package consumer
/*
 Author: hawrkchen
 Date: 2022/11/28 10:48
 Desc:
*/
package consumer

import (
	"algo_assess/market-mq-server/internal/svc"
	pb "algo_assess/market-mq-server/proto"
	"algo_assess/models"
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type PerfFixSHMarketInfo struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPerfFixSHMarketInfo(ctx context.Context, svcCtx *svc.ServiceContext) *PerfFixSHMarketInfo {
	return &PerfFixSHMarketInfo{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (s *PerfFixSHMarketInfo) Consume(_, val string) error {
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
	if pre == "1" || pre == "2" || pre == "0" { // 沪市一般只保留6开头
		//s.Logger.Info("unnecessary secID.... ")
		return nil
	}
	shQuote := TransQuoteWithPB(data, 2)
	s.Logger.Infof("get shQuote:%+v", shQuote)
	// 行情数据修复正常redis的数据都已经落到DB了，只能查DB的数据
	out, err := s.svcCtx.ShMarketLevelRepo.QueryShMarketLevel(s.ctx, shQuote.SecID, shQuote.OrigTime)
	if err != nil {
		s.Logger.Error("QueryShMarketLevel error:", err)
		return nil
	}
	if len(out) > 0 { // 记录已存在，不再处理
		s.Logger.Info("record has exist...")
		return nil
	}
	// 记录不存在，开始计算
	// 计算净成交量时，需要从DB中取出上一个时间点的成交量记录
	quote, err := s.svcCtx.ShMarketLevelRepo.QueryAllDataBySecId(s.ctx, shQuote.SecID, shQuote.OrigTime)
	if err != nil {
		s.Logger.Error("QueryAllDataBySecId error:", err)
		return nil
	}
	preTradeVol := GetShPreTradeVol(quote, shQuote.OrigTime)
	netTotalVol := CalculateFixMarketVwap(&shQuote, preTradeVol)
	//保存到内存中供绩效计算
	if err := BuildLevel2Data(s.svcCtx.RedisClient, netTotalVol, &shQuote); err != nil {
		s.Logger.Error("build level2 data fail:", err)
		return nil
	}
	// 回写一次DB表  --直接插入
	if err := s.svcCtx.ShMarketLevelRepo.CreateShMarketLevel(s.ctx, &shQuote); err != nil {
		s.Logger.Error("insert db market level fail:", err)
	}

	return nil
}

// GetShPreTradeVol 从已排序的列表找到上一条记录
// originTime为当前修复的数据时间点， out: 输出上一个时间点的成交数量
func GetShPreTradeVol(in []*models.TbShQuoteLevel, originTime int64) int64 {
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
