// Package repo
/*
 Author: hawrkchen
 Date: 2022/4/24 18:08
 Desc:
*/
package repo

import (
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
)

type MarketLevelRepo interface {
	CreateMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	UpdateMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	GetMarketLevelByTime(ctx context.Context, secId string, start, end int64) ([]*models.TbSzQuoteLevel, *gorm.DB)
}

type DefaultMarketLevel struct {
	DB *gorm.DB
}

func NewMarketLevelRepo(conn *gorm.DB) MarketLevelRepo {
	return &DefaultMarketLevel{
		DB: conn,
	}
}

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
	}
	results := d.DB.Create(&marketLevel)
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
