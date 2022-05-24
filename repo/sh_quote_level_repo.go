// Package repo
/*
 Author: hawrkchen
 Date: 2022/5/16 19:00
 Desc:
*/
package repo

import (
	"algo_assess/global"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
)

type SHMarketLevelRepo interface {
	CreateShMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	UpdateShMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	GetShMarketLevelByTime(ctx context.Context, secId string, start, end int64) ([]*models.TbShQuoteLevel, *gorm.DB)
}

type DefaultSHMarketLevelRepo struct {
	DB *gorm.DB
}

func NewSHMarketLevelRepo(conn *gorm.DB) SHMarketLevelRepo {
	return &DefaultSHMarketLevelRepo{
		DB: conn,
	}
}

func (d *DefaultSHMarketLevelRepo) CreateShMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error {
	marketLevel := models.TbShQuoteLevel{
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

func (d *DefaultSHMarketLevelRepo) UpdateShMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error {
	marketLevel := &models.TbShQuoteLevel{
		LastPrice:     data.LastPrice,
		AskPrice:      data.AskPrice,
		AskVol:        data.AskVol,
		BidPrice:      data.BidPrice,
		BidVol:        data.BidVol,
		TotalTradeVol: data.TotalTradeVol,
		MkVwap:        data.Vwap,
	}
	result := d.DB.Model(models.TbShQuoteLevel{}).Where("seculity_id=? and orgi_time=?", data.SecID, data.OrigTime).
		Updates(marketLevel)
	return result.Error
}

func (d *DefaultSHMarketLevelRepo) GetShMarketLevelByTime(ctx context.Context, secId string, start, end int64) ([]*models.TbShQuoteLevel, *gorm.DB) {
	var level []*models.TbShQuoteLevel
	result := d.DB.Select("orgi_time, last_price, ask_price, ask_vol, bid_price, bid_vol, total_trade_vol, mk_vwap").
		Where("seculity_id = ? and orgi_time between ? and ?",
			secId, start, end).
		Find(&level)
	if result.Error != nil {
		return nil, result
	}
	return level, result
}
