// Package repo
/*
 Author: hawrkchen
 Date: 2022/4/24 18:08
 Desc:
*/
package repo

import (
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type MarketLevelRepo interface {
	CreateMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	CreateMarketLevelBatch(ctx context.Context, data []*global.QuoteLevel2Data) error
	UpdateMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	GetMarketLevelByTime(ctx context.Context, secId string, start, end int64) ([]*models.TbSzQuoteLevel, *gorm.DB)
	GetSzMarketLevelByPageWithCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) (int64, []*models.TbSzQuoteLevel, error)
	GetSzMarketLevelByPageWithoutCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) ([]*models.TbSzQuoteLevel, error)
	GetSzMarketLevelCount(ctx context.Context) int64
	GetSzMarketLevelBySecId(ctx context.Context, secId, date string) ([]*models.TbSzQuoteLevel, error)
	GetSzClosingPrice(ctx context.Context, secId, date string) (int64, error)
	// 数据修复用
	QuerySzMarketLevel(ctx context.Context, secId string, orgiTime int64) ([]*models.TbSzQuoteLevel, error)
	QueryAllDataBySecId(ctx context.Context, secId string, orgiTime int64) ([]*models.TbSzQuoteLevel, error)
}

type DefaultMarketLevel struct {
	DB *gorm.DB
}

func NewMarketLevelRepo(conn *gorm.DB) MarketLevelRepo {
	return &DefaultMarketLevel{
		DB: conn,
	}
}

// CreateMarketLevel 数据修复调用，正常数据走批量
func (d *DefaultMarketLevel) CreateMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error {
	marketLevel := models.TbSzQuoteLevel{
		SeculityId:    data.SecID,
		OrgiTime:      data.OrigTime,
		LastPrice:     data.LastPrice,
		AskPrice:      data.AskPrice,
		AskVol:        data.AskVol,
		BidPrice:      data.BidPrice,
		BidVol:        data.BidVol,
		TotalTradeVol: data.TotalTradeVol,
		TotalAskVol:   data.TotalAskVol,
		TotalBidVol:   data.TotalBidVol,
		MkVwap:        data.Vwap,
		FixFlag:       1,
	}
	results := d.DB.Create(&marketLevel)
	return results.Error
}

func (d *DefaultMarketLevel) CreateMarketLevelBatch(ctx context.Context, data []*global.QuoteLevel2Data) error {
	batch := make([]models.TbSzQuoteLevel, 0, len(data))
	for _, v := range data {
		marketLevel := models.TbSzQuoteLevel{
			SeculityId:    v.SecID,
			OrgiTime:      v.OrigTime,
			LastPrice:     v.LastPrice,
			AskPrice:      v.AskPrice,
			AskVol:        v.AskVol,
			BidPrice:      v.BidPrice,
			BidVol:        v.BidVol,
			TotalTradeVol: v.TotalTradeVol,
			TotalAskVol:   v.TotalAskVol,
			TotalBidVol:   v.TotalBidVol,
			MkVwap:        v.Vwap,
		}
		batch = append(batch, marketLevel)
	}
	results := d.DB.CreateInBatches(batch, len(batch))
	return results.Error
}

func (d *DefaultMarketLevel) UpdateMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error {
	marketLevel := &models.TbSzQuoteLevel{
		LastPrice:     data.LastPrice,
		AskPrice:      data.AskPrice,
		AskVol:        data.AskVol,
		BidPrice:      data.BidPrice,
		BidVol:        data.BidVol,
		TotalTradeVol: data.TotalTradeVol,
		MkVwap:        data.Vwap,
	}
	result := d.DB.Model(models.TbSzQuoteLevel{}).Where("seculity_id=? and orgi_time=?", data.SecID, data.OrigTime).
		Updates(marketLevel)
	return result.Error
}

// GetMarketLevelByTime 客户端行情数据
func (d *DefaultMarketLevel) GetMarketLevelByTime(ctx context.Context, secId string, start, end int64) ([]*models.TbSzQuoteLevel, *gorm.DB) {
	var level []*models.TbSzQuoteLevel
	result := d.DB.Select("orgi_time, last_price, ask_price, ask_vol, bid_price, bid_vol, total_trade_vol, mk_vwap").
		Where("seculity_id = ? and orgi_time between ? and ?",
			secId, start, end).
		Find(&level)
	if result.Error != nil {
		return nil, result
	}
	return level, result
}

// GetSzMarketLevelByPageWithCount 带统计数量，seculity_id 必填
func (d *DefaultMarketLevel) GetSzMarketLevelByPageWithCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) (int64, []*models.TbSzQuoteLevel, error) {
	var infos []*models.TbSzQuoteLevel
	var count int64
	model := d.DB.Model(&models.TbSzQuoteLevel{}).Where("seculity_id", "sz:"+in.GetSecurityId())
	err := model.Count(&count).Limit(int(in.GetPageNum())).Offset(int((in.GetPageId() - 1) * in.GetPageNum())).Find(&infos).Error
	if err != nil {
		return 0, nil, err
	}
	return count, infos, nil
}

// GetSzMarketLevelByPageWithoutCount 不带统计数量，只取数据
func (d *DefaultMarketLevel) GetSzMarketLevelByPageWithoutCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) ([]*models.TbSzQuoteLevel, error) {
	var infos []*models.TbSzQuoteLevel
	// 不能直接偏移，会有性能问题，需要根据前一刷的最后一个ID偏移位置取值
	//result := d.DB.Model(&models.TbSzQuoteLevel{}).Limit(int(in.GetPageNum())).Offset(int((in.GetPageId() - 1) * in.GetPageNum())).Find(&infos)
	//offsetId := (in.GetPageId() - 1) * in.GetPageNum()
	result := d.DB.Model(&models.TbSzQuoteLevel{}).Where("id > ?", in.GetMaxId()).Limit(int(in.GetPageNum())).Find(&infos)
	if result.Error != nil {
		return nil, result.Error
	}
	return infos, nil
}

func (d *DefaultMarketLevel) GetSzMarketLevelCount(ctx context.Context) int64 {
	type tableStatus struct {
		Name      string
		Engine    string
		Version   int64
		RowFormat string
		Rows      int64
	}
	var t tableStatus
	d.DB.Raw("show table status like 'tb_sz_quote_level';").Scan(&t)
	return t.Rows
}

// GetSzMarketLevelBySecId 绩效计算需要取到达价，如果总线不传，需要从数据库里取,取单支证券一天数据
func (d *DefaultMarketLevel) GetSzMarketLevelBySecId(ctx context.Context, secId, date string) ([]*models.TbSzQuoteLevel, error) {
	var out []*models.TbSzQuoteLevel
	result := d.DB.Select("orgi_time,last_price").Where("seculity_id=? and orgi_time like ?", "sz:"+secId, date+"%").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetSzClosingPrice 取收盘价，用作行情兜底用
func (d *DefaultMarketLevel) GetSzClosingPrice(ctx context.Context, secId, date string) (int64, error) {
	var closingPrice []int64
	result := d.DB.Model(models.TbSzQuoteLevel{}).Select("last_price").Where("seculity_id=? and orgi_time like ?", "sz:"+secId, date+"%").
		Order("orgi_time desc").Find(&closingPrice)
	if result.Error != nil {
		return 0, result.Error
	}
	if len(closingPrice) >= 1 {
		return closingPrice[0], nil
	}
	return 0, nil
}

func (d *DefaultMarketLevel) QuerySzMarketLevel(ctx context.Context, secId string, orgiTime int64) ([]*models.TbSzQuoteLevel, error) {
	var level []*models.TbSzQuoteLevel
	result := d.DB.Select("seculity_id,orgi_time").Where("orgi_time=? and seculity_id = ?", orgiTime, secId).Find(&level)
	if result.Error != nil {
		return nil, result.Error
	}
	return level, nil
}

// QueryAllDataBySecId 查单个证券ID当天的所有记录
func (d *DefaultMarketLevel) QueryAllDataBySecId(ctx context.Context, secId string, orgiTime int64) ([]*models.TbSzQuoteLevel, error) {
	t := cast.ToString(orgiTime)[:8] // 取当天时间
	var out []*models.TbSzQuoteLevel
	result := d.DB.Select("orgi_time, total_trade_vol").Where("orgi_time like ? and seculity_id=?", t+"%", secId).
		Order("orgi_time").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}
