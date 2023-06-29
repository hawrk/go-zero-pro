// Package busrepo
/*
 Author: hawrkchen
 Date: 2022/9/27 10:35
 Desc:
*/
package busrepo

import (
	"algo_assess/busmodels"
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"
	"gorm.io/gorm"
)

type AlgoInfoRepo interface {
	GetAlgoNameByIds(ctx context.Context, ids []int, provider string, algoType int) ([]string, error)
	GetAlgoProviders(ctx context.Context, ids []int) ([]string, error)
	GetAlgoTypeNames(ctx context.Context, ids []int, provider string) ([]string, error)
	GetBusAlgoBase(ctx context.Context) ([]*busmodels.TbAlgoInfo, error)
	GetAlgoInfoByName(ctx context.Context, algoName string) (string, string, error)
}

type DefaultAlgoInfo struct {
	DB *gorm.DB
}

func NewAlgoInfo(conn *gorm.DB) AlgoInfoRepo {
	return &DefaultAlgoInfo{
		DB: conn,
	}
}

// GetAlgoNameByIds 根据算法ID列表查询算法名称列表
func (d *DefaultAlgoInfo) GetAlgoNameByIds(ctx context.Context, ids []int, provider string, algoType int) ([]string, error) {
	var name []string
	result := d.DB.Model(busmodels.TbAlgoInfo{}).Select("distinct(algo_name)").
		Joins("inner join  tb_user_info b on tb_algo_info.id  in ? and tb_algo_info.uuser_id= b.id", ids)
	if provider != "" {
		result = result.Where("b.user_name like ?", provider+"%")
	}
	if algoType != 0 {
		result = result.Where("algorithm_type=?", algoType)
	}
	result = result.Find(&name)
	if result.Error != nil {
		return nil, result.Error
	}
	var out []string
	for _, v := range name {
		out = append(out, tools.RMu0000(v))
	}
	return out, nil
}

// GetAlgoProviders 根据算法ID列表查询厂商列表（去重)
func (d *DefaultAlgoInfo) GetAlgoProviders(ctx context.Context, ids []int) ([]string, error) {
	var providerName []string
	// 需要联表
	result := d.DB.Model(busmodels.TbAlgoInfo{}).Select("distinct(b.user_name) as provider_name").
		Joins("inner join  tb_user_info b on  tb_algo_info.id  in ? and tb_algo_info.uuser_id= b.id", ids).
		Find(&providerName)
	//result := d.DB.Model(busmodels.TbAlgoInfo{}).Select("distinct(provider_name)").Where("id in ?", ids).Find(&providerName)
	if result.Error != nil {
		return nil, result.Error
	}
	var out []string
	for _, v := range providerName {
		out = append(out, tools.RMu0000(v))
	}
	return out, nil
}

func (d *DefaultAlgoInfo) GetAlgoTypeNames(ctx context.Context, ids []int, provider string) ([]string, error) {
	var algoType []int
	result := d.DB.Model(busmodels.TbAlgoInfo{}).Select("distinct(algorithm_type)").Where("id in ?", ids).Find(&algoType)
	//Where("provider_name like ?", provider+"%").Find(&algoType) // 总线的数据每个字段后面都会补空格，需要加上Like
	if result.Error != nil {
		return nil, result.Error
	}
	// algo_type 转换为algo_type_name
	var name []string
	for _, v := range algoType {
		if v == 1 {
			name = append(name, global.AlgoTypeNameT0)
		} else if v == 2 {
			name = append(name, global.AlgoTypeNameSplit)
		}
	}
	return name, nil
}

// GetBusAlgoBase 取所有算法基本数据
func (d *DefaultAlgoInfo) GetBusAlgoBase(ctx context.Context) ([]*busmodels.TbAlgoInfo, error) {
	var infos []*busmodels.TbAlgoInfo
	// 总线取消provider_name 字段，需要进行联表查询
	//result := d.DB.Model(&busmodels.TbAlgoInfo{}).Select("id,algo_name,algorithm_type,provider_name").
	//	Where("algorithm_status=3").Find(&infos)
	result := d.DB.Model(&busmodels.TbAlgoInfo{}).Select("tb_algo_info.id, algo_name, algorithm_type, tb_algo_info.uuser_id, b.user_name as provider_name").
		Joins("inner join  tb_user_info b on algorithm_status = 3 and tb_algo_info.uuser_id= b.id").Scan(&infos)
	if result.Error != nil && result.RowsAffected != 0 {
		return nil, result.Error
	}
	return infos, nil
}

// GetAlgoInfoByName 根据算法名称反查厂商，算法类型名称信息
func (d *DefaultAlgoInfo) GetAlgoInfoByName(ctx context.Context, algoName string) (string, string, error) {
	var out busmodels.TbAlgoInfo
	//result := d.DB.Model(&busmodels.TbAlgoInfo{}).Select("uuser_id, algorithm_type").
	//Where("algo_name like ?", algoName+"%").Find(&out)
	result := d.DB.Model(&busmodels.TbAlgoInfo{}).Select("b.user_name as provider_name, algorithm_type").
		Joins("inner join  tb_user_info b on algo_name like ? and tb_algo_info.uuser_id= b.id", algoName+"%").
		Find(&out)
	if result.Error != nil {
		return "", "", result.Error
	}
	var algoTypeName string
	provider := tools.RMu0000(out.ProviderName)
	if out.AlgorithmType == 1 {
		algoTypeName = global.AlgoTypeNameT0
	} else if out.AlgorithmType == 2 {
		algoTypeName = global.AlgoTypeNameSplit
	}
	return provider, algoTypeName, nil
}
