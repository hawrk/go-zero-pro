// Package repo
/*
 Author: hawrkchen
 Date: 2022/7/14 16:09
 Desc:
*/
package repo

import (
	"algo_assess/models"
	"gorm.io/gorm"
)

type AlgoOptimizeRepo interface {
	GetAllAlgoOptimize(algoType int32, page, limit int) (int64, []*models.TbAlgoOptimize, error)
	GetAlgoOptimizeBySecurityId(list []string) (int64, []*models.TbAlgoOptimize, error)
	GetSingleAlgoOptimize(securityId string) (*models.TbAlgoOptimize, error)
	GetAlgoOptimize(securityId string, algoId int) (*models.TbAlgoOptimize, error)
	AddAlgoOptimize(in *models.TbAlgoOptimize) error
	UpdateAlgoOptimize(in *models.TbAlgoOptimize) error
	InitAlgoOptimize(in *models.TbAlgoOptimize) error
}

type DefaultAlgoOptimize struct {
	DB *gorm.DB
}

func NewAlgoOptimizeRepo(conn *gorm.DB) AlgoOptimizeRepo {
	return &DefaultAlgoOptimize{
		DB: conn,
	}
}

func (d *DefaultAlgoOptimize) GetAllAlgoOptimize(algoType int32, page, limit int) (int64, []*models.TbAlgoOptimize, error) {
	var infos []*models.TbAlgoOptimize
	var count int64
	result := d.DB.Model(&models.TbAlgoOptimize{}).Select("id, sec_id, sec_name, algo_name").
		Where("algo_type=?", algoType).Count(&count).Limit(limit).Offset((page - 1) * limit).Find(&infos)
	if result.Error != nil {
		return count, nil, result.Error
	}
	return count, infos, nil
}

func (d *DefaultAlgoOptimize) GetAlgoOptimizeBySecurityId(list []string) (int64, []*models.TbAlgoOptimize, error) {
	var infos []*models.TbAlgoOptimize
	var count int64
	result := d.DB.Model(&models.TbAlgoOptimize{}).Select("id,provider_id,provider_name, sec_id, sec_name,algo_id ,algo_name").
		Where("sec_id in ?", list).Count(&count). /*.Limit(limit).Offset((page - 1) * limit)*/ Find(&infos)
	if result.Error != nil {
		return count, nil, result.Error
	}
	return count, infos, nil
}

func (d *DefaultAlgoOptimize) GetSingleAlgoOptimize(securityId string) (*models.TbAlgoOptimize, error) {
	var info *models.TbAlgoOptimize
	var count int64
	result := d.DB.Model(&models.TbAlgoOptimize{}).Select("id, score").Where("sec_id = ?", securityId).Count(&count).Find(&info)
	if result.Error != nil {
		return nil, result.Error
	}
	if count == 0 {
		return nil, nil
	}
	return info, nil
}

func (d *DefaultAlgoOptimize) AddAlgoOptimize(in *models.TbAlgoOptimize) error {
	return d.DB.Create(in).Error
}

func (d *DefaultAlgoOptimize) UpdateAlgoOptimize(in *models.TbAlgoOptimize) error {
	return d.DB.Updates(in).Error
}

func (d *DefaultAlgoOptimize) InitAlgoOptimize(in *models.TbAlgoOptimize) error {
	return d.DB.Model(in).Select("algo_id", "algo_type", "algo_name", "score", "update_time").Updates(in).Error
}

func (d *DefaultAlgoOptimize) GetAlgoOptimize(securityId string, algoId int) (*models.TbAlgoOptimize, error) {
	var info *models.TbAlgoOptimize
	var count int64
	result := d.DB.Model(&models.TbAlgoOptimize{}).Where("sec_id = ?", securityId).Where("algo_id", algoId).Count(&count).Find(&info)
	if result.Error != nil {
		return nil, result.Error
	}
	if count == 0 {
		return nil, nil
	}
	return info, nil
}
