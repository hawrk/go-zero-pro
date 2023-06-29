// Package busrepo
/*
 Author: hawrkchen
 Date: 2023/3/23 15:54
 Desc:
*/
package busrepo

import (
	"algo_assess/busmodels"
	"context"
	"gorm.io/gorm"
)

type AssetInfoRepo interface {
	GetAssetInfoById(ctx context.Context, id int64) (busmodels.TbAssetInfo, error)
}

type DefaultAssetInfo struct {
	DB *gorm.DB
}

func NewAssetInfo(conn *gorm.DB) AssetInfoRepo {
	return &DefaultAssetInfo{
		DB: conn,
	}
}

func (d *DefaultAssetInfo) GetAssetInfoById(ctx context.Context, id int64) (busmodels.TbAssetInfo, error) {
	var out busmodels.TbAssetInfo
	result := d.DB.Model(busmodels.TbAssetInfo{}).Select("balance,frozen,currency_type,account_type,cust_orgid,cust_branchid").
		Where("id=?", id).Find(&out)
	if result.Error != nil {
		return busmodels.TbAssetInfo{}, result.Error
	}

	return out, nil
}
