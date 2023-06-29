// Package repo
/*
 Author: hawrkchen
 Date: 2022/5/16 19:00
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

type SHMarketLevelRepo interface {
	CreateShMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	CreateMarketLevelBatch(ctx context.Context, data []*global.QuoteLevel2Data) error
	UpdateShMarketLevel(ctx context.Context, data *global.QuoteLevel2Data) error
	GetShMarketLevelByTime(ctx context.Context, secId string, start, end int64) ([]*models.TbShQuoteLevel, *gorm.DB)
	GetShMarketLevelByPageWithCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) (int64, []*models.TbShQuoteLevel, error)
	GetShMarketLevelByPageWithoutCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) ([]*models.TbShQuoteLevel, error)
	GetShMarketLevelCount(ctx context.Context) int64
	GetShMarketLevelBySecId(ctx context.Context, secId, date string) ([]*models.TbShQuoteLevel, error)
	GetShClosingPrice(ctx context.Context, secId, date string) (int64, error)

	// 数据修复用
	QueryShMarketLevel(ctx context.Context, secId string, orgiTime int64) ([]*models.TbShQuoteLevel, error)
	QueryAllDataBySecId(ctx context.Context, secId string, orgiTime int64) ([]*models.TbShQuoteLevel, error)
}

type DefaultSHMarketLevelRepo struct {
	DB *gorm.DB
}

func NewSHMarketLevelRepo(conn *gorm.DB) SHMarketLevelRepo {
	return &DefaultSHMarketLevelRepo{
		DB: conn,
	}
}

// CreateShMarketLevel 数据修复调用，正常数据走批量
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
		FixFlag:       1,
	}
	results := d.DB.Create(&marketLevel)
	return results.Error
}

func (d *DefaultSHMarketLevelRepo) CreateMarketLevelBatch(ctx context.Context, data []*global.QuoteLevel2Data) error {
	batch := make([]models.TbShQuoteLevel, 0, len(data))
	for _, v := range data {
		marketLevel := models.TbShQuoteLevel{
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

func (d *DefaultSHMarketLevelRepo) GetShMarketLevelByPageWithCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) (int64, []*models.TbShQuoteLevel, error) {
	var infos []*models.TbShQuoteLevel
	var count int64
	model := d.DB.Model(&models.TbShQuoteLevel{}).Where("seculity_id=?", "sh:"+in.GetSecurityId())
	err := model.Count(&count).Limit(int(in.GetPageNum())).Offset(int((in.GetPageId() - 1) * in.GetPageNum())).Find(&infos).Error
	if err != nil {
		return 0, nil, err
	}
	return count, infos, nil
}

func (d *DefaultSHMarketLevelRepo) GetShMarketLevelByPageWithoutCount(ctx context.Context, in *proto.ReqQueryQuoteLevel) ([]*models.TbShQuoteLevel, error) {
	var infos []*models.TbShQuoteLevel
	//offsetId := (in.GetPageId() - 1) * in.GetPageNum()
	result := d.DB.Model(&models.TbShQuoteLevel{}).Where("id > ?", in.GetMaxId()).Limit(int(in.GetPageNum())).Find(&infos)
	if result.Error != nil {
		return nil, result.Error
	}
	return infos, nil
}

func (d *DefaultSHMarketLevelRepo) GetShMarketLevelCount(ctx context.Context) int64 {
	type tableStatus struct {
		Name      string
		Engine    string
		Version   int64
		RowFormat string
		Rows      int64
	}
	var t tableStatus
	d.DB.Raw("show table status like 'tb_sh_quote_level';").Scan(&t)
	return t.Rows
}

// GetShMarketLevelBySecId 绩效计算需要取到达价，如果总线不传，需要从数据库里取,取单支证券一天数据
func (d *DefaultSHMarketLevelRepo) GetShMarketLevelBySecId(ctx context.Context, secId, date string) ([]*models.TbShQuoteLevel, error) {
	var out []*models.TbShQuoteLevel
	result := d.DB.Select("orgi_time,last_price").Where("seculity_id=? and orgi_time like ?", "sh:"+secId, date+"%").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultSHMarketLevelRepo) GetShClosingPrice(ctx context.Context, secId, date string) (int64, error) {
	var closingPrice []int64
	result := d.DB.Model(models.TbShQuoteLevel{}).Select("last_price").Where("seculity_id=? and orgi_time like ?", "sh:"+secId, date+"%").
		Order("orgi_time desc").Find(&closingPrice)
	if result.Error != nil {
		return 0, result.Error
	}
	if len(closingPrice) >= 1 {
		return closingPrice[0], nil
	}
	return 0, nil
}

func (d *DefaultSHMarketLevelRepo) QueryShMarketLevel(ctx context.Context, secId string, orgiTime int64) ([]*models.TbShQuoteLevel, error) {
	var level []*models.TbShQuoteLevel
	result := d.DB.Select("seculity_id,orgi_time").Where("orgi_time=? and seculity_id = ?", orgiTime, secId).Find(&level)
	if result.Error != nil {
		return nil, result.Error
	}
	return level, nil
}

// QueryAllDataBySecId 查单个证券ID当天的所有记录
func (d *DefaultSHMarketLevelRepo) QueryAllDataBySecId(ctx context.Context, secId string, orgiTime int64) ([]*models.TbShQuoteLevel, error) {
	t := cast.ToString(orgiTime)[:8] // 取当天时间
	var out []*models.TbShQuoteLevel
	result := d.DB.Select("orgi_time, total_trade_vol").Where("orgi_time like ? and seculity_id=?", t+"%", secId).
		Order("orgi_time").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}
