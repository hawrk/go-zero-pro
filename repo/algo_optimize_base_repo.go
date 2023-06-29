// Package repo
/*
 Author: hawrkchen
 Date: 2022/7/14 16:09
 Desc:
*/
package repo

import (
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/models"
	"bytes"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AlgoOptimizeBaseRepo interface {
	AddAlgoOptimizeBase(in *proto.AddOptimizeBaseReq) error
	UpdateAlgoOptimizeBase(in *proto.UpdateOptimizeBaseReq) error
	UpdateAlgoOptimizeBaseBySecIdAndAlgoId(in *proto.UpdateOptimizeBaseReq) error
	DeleteAlgoOptimizeBase(id int64) error
	SelectAlgoOptimizeBase(in *proto.SelectOptimizeBaseReq) (int64, []*models.TbAlgoOptimizeBase, error)
	SelectOptimizeBaseBySecId(secId string) ([]*models.TbAlgoOptimizeBase, error)
	GetOptimize(secIds []string, algoIds []int32) (int64, []*models.TbAlgoOptimizeBase, error)
	SelectAllOptimizeBase() (int64, []*models.TbAlgoOptimizeBase, error)
	SelectOptimizeBaseById(id int64) (*models.TbAlgoOptimizeBase, error)
	CountOptimize(secId string, algoId int) (int64, error)
	BatchUpload(in *proto.UploadOptimizeBaseReq) error
}

type DefaultAlgoBaseOptimize struct {
	DB *gorm.DB
}

func NewAlgoOptimizeBaseRepo(conn *gorm.DB) AlgoOptimizeBaseRepo {
	return &DefaultAlgoBaseOptimize{
		DB: conn,
	}
}

func (d *DefaultAlgoBaseOptimize) AddAlgoOptimizeBase(in *proto.AddOptimizeBaseReq) error {
	algoOptimizeBase := &models.TbAlgoOptimizeBase{
		ProviderId:   int(in.GetProviderId()),
		ProviderName: in.GetProviderName(),
		SecId:        in.GetSecId(),
		SecName:      in.GetSecName(),
		AlgoId:       int(in.GetAlgoId()),
		AlgoType:     int(in.GetAlgoType()),
		AlgoName:     in.GetAlgoName(),
		OpenRate:     in.GetOpenRate(),
		IncomeRate:   in.GetIncomeRate(),
		BasisPoint:   in.GetBasisPoint(),
		CreateTime:   time.Now(),
		UpdateTime:   time.Now(),
	}
	return d.DB.Create(algoOptimizeBase).Error
}

func (d *DefaultAlgoBaseOptimize) UpdateAlgoOptimizeBaseBySecIdAndAlgoId(in *proto.UpdateOptimizeBaseReq) error {
	algoOptimizeBase := &models.TbAlgoOptimizeBase{
		ProviderId:   int(in.GetProviderId()),
		ProviderName: in.GetProviderName(),
		SecName:      in.GetSecName(),
		AlgoType:     int(in.GetAlgoType()),
		AlgoName:     in.GetAlgoName(),
		OpenRate:     in.GetOpenRate(),
		IncomeRate:   in.GetIncomeRate(),
		BasisPoint:   in.GetBasisPoint(),
		UpdateTime:   time.Now(),
	}
	return d.DB.Model(algoOptimizeBase).Where("sec_id=?", in.GetSecId()).Where("algo_id=?", in.GetAlgoId()).Updates(algoOptimizeBase).Error
}

func (d *DefaultAlgoBaseOptimize) UpdateAlgoOptimizeBase(in *proto.UpdateOptimizeBaseReq) error {
	algoOptimizeBase := &models.TbAlgoOptimizeBase{
		Id:           in.GetId(),
		ProviderId:   int(in.GetProviderId()),
		ProviderName: in.GetProviderName(),
		SecId:        in.GetSecId(),
		SecName:      in.GetSecName(),
		AlgoId:       int(in.GetAlgoId()),
		AlgoType:     int(in.GetAlgoType()),
		AlgoName:     in.GetAlgoName(),
		OpenRate:     in.GetOpenRate(),
		IncomeRate:   in.GetIncomeRate(),
		BasisPoint:   in.GetBasisPoint(),
		UpdateTime:   time.Now(),
	}
	return d.DB.Updates(algoOptimizeBase).Error
}

func (d *DefaultAlgoBaseOptimize) DeleteAlgoOptimizeBase(id int64) error {
	return d.DB.Where("id=?", id).Delete(&models.TbAlgoOptimizeBase{}).Error
}

func (d *DefaultAlgoBaseOptimize) SelectAlgoOptimizeBase(in *proto.SelectOptimizeBaseReq) (int64, []*models.TbAlgoOptimizeBase, error) {
	var infos []*models.TbAlgoOptimizeBase
	var count int64
	model := d.DB.Model(&models.TbAlgoOptimizeBase{})
	providerId := in.GetProviderId()
	if providerId != 0 {
		model.Where("provider_id=?", providerId)
	}
	algoId := in.GetAlgoId()
	if algoId != 0 {
		model.Where("algo_id", algoId)
	}
	secId := in.GetSecId()
	if secId != "" {
		model.Where("sec_id", secId)
	}
	err := model.Count(&count).Limit(int(in.GetLimit())).Offset(int((in.GetPage() - 1) * in.GetLimit())).Find(&infos).Error
	return count, infos, err
}

func (d *DefaultAlgoBaseOptimize) SelectOptimizeBaseBySecId(secId string) ([]*models.TbAlgoOptimizeBase, error) {
	var infos []*models.TbAlgoOptimizeBase
	err := d.DB.Model(&models.TbAlgoOptimizeBase{}).Where("sec_id", secId).Find(&infos).Error
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (d *DefaultAlgoBaseOptimize) SelectAllOptimizeBase() (int64, []*models.TbAlgoOptimizeBase, error) {
	var infos []*models.TbAlgoOptimizeBase
	var count int64
	err := d.DB.Model(&models.TbAlgoOptimizeBase{}).Count(&count).Find(&infos).Error
	if err != nil {
		return 0, nil, err
	}
	return count, infos, nil
}

func (d *DefaultAlgoBaseOptimize) SelectOptimizeBaseById(id int64) (*models.TbAlgoOptimizeBase, error) {
	var infos *models.TbAlgoOptimizeBase
	err := d.DB.Model(&models.TbAlgoOptimizeBase{}).Where("id", id).Find(&infos).Error
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func (d *DefaultAlgoBaseOptimize) CountOptimize(secId string, algoId int) (int64, error) {
	var count int64
	err := d.DB.Model(&models.TbAlgoOptimizeBase{}).Where("sec_id", secId).Where("algo_id", algoId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *DefaultAlgoBaseOptimize) GetOptimize(secIds []string, algoIds []int32) (int64, []*models.TbAlgoOptimizeBase, error) {
	var infos []*models.TbAlgoOptimizeBase
	var count int64
	//result := d.DB.Raw(`select a.* from tb_algo_optimize_base a inner join (select id, sec_id,algo_id, max(open_rate + income_rate + basis_point) scope
	//                 from tb_algo_optimize_base
	//                 WHERE sec_id in (?)
	//                   AND algo_id in (?)
	//                 group by sec_id) b on a.id = b.id and a.sec_id=b.sec_id and a.algo_id = b.algo_id
	//				group by a.sec_id`, secIds, algoIds).Scan(&infos)
	result := d.DB.Raw(`select *
							from tb_algo_optimize_base a,
								 (
									 select sec_id, max(open_rate + income_rate + basis_point) scope
									 from tb_algo_optimize_base b
									 where sec_id in (?)
									   AND algo_id in (?)
									 group by sec_id
								 ) b
							where a.sec_id = b.sec_id
							  and (a.basis_point + a.income_rate + a.open_rate) = b.scope group by a.sec_id;`, secIds, algoIds).Scan(&infos)
	if result.Error != nil {
		return count, nil, result.Error
	}
	return int64(len(infos)), infos, nil
}

func (d *DefaultAlgoBaseOptimize) BatchUpload(in *proto.UploadOptimizeBaseReq) error {
	list := in.GetList()
	if len(list) == 0 {
		return nil
	}
	var buffer bytes.Buffer
	sql := "replace into tb_algo_optimize_base(id, provider_id, provider_name, sec_id, sec_name, algo_id, algo_type, algo_name, open_rate, income_rate, basis_point, create_time, update_time) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	var execErr error
	for i, o := range list {
		if len(list)-i < 1000 {
			// 后半部的拼接
			if i == len(list)-1 {
				buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s','%s','%d','%d','%s',%f,%f,%f,'%s','%s')", o.Id, o.ProviderId, o.ProviderName, o.SecId, o.SecName, o.AlgoId, o.AlgoType, o.AlgoName, o.OpenRate, o.IncomeRate, o.BasisPoint, o.CreateTime, o.UpdateTime))
				execErr = d.DB.Exec(buffer.String()).Error
				if execErr != nil {
					return execErr
				}
				buffer.Reset()
				buffer.WriteString(sql)
			} else {
				buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s','%s','%d','%d','%s',%f,%f,%f,'%s','%s'),", o.Id, o.ProviderId, o.ProviderName, o.SecId, o.SecName, o.AlgoId, o.AlgoType, o.AlgoName, o.OpenRate, o.IncomeRate, o.BasisPoint, o.CreateTime, o.UpdateTime))
			}
		} else {
			if i%1000 == 0 {
				buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s','%s','%d','%d','%s',%f,%f,%f,'%s','%s')", o.Id, o.ProviderId, o.ProviderName, o.SecId, o.SecName, o.AlgoId, o.AlgoType, o.AlgoName, o.OpenRate, o.IncomeRate, o.BasisPoint, o.CreateTime, o.UpdateTime))
				execErr = d.DB.Exec(buffer.String()).Error
				if execErr != nil {
					return execErr
				}
				buffer.Reset()
				buffer.WriteString(sql)
			} else {
				buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s','%s','%d','%d','%s',%f,%f,%f,'%s','%s'),", o.Id, o.ProviderId, o.ProviderName, o.SecId, o.SecName, o.AlgoId, o.AlgoType, o.AlgoName, o.OpenRate, o.IncomeRate, o.BasisPoint, o.CreateTime, o.UpdateTime))
			}
		}
	}
	return execErr
}
