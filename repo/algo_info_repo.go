// Package repo
/*
 Author: hawrkchen
 Date: 2022/6/23 9:49
 Desc:
*/
package repo

import (
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"
	"gorm.io/gorm"
	"time"
)

type AlgoInfoRepo interface {
	GetAllAlgoInfo(ctx context.Context) ([]*AllAlgoInfo, error)
	GetAlgoBase(ctx context.Context) ([]*models.TbAlgoInfo, error)
	GetAlgoBaseById(ctx context.Context, id uint32) ([]*models.TbAlgoInfo, error)
	UpdateAlgoBaseById(ctx context.Context, id uint32, perf *pb.AlgoInfoPerf) error
	CreateAlgoBaseInfo(ctx context.Context, perf *pb.AlgoInfoPerf) error
	GetAlgoSummary(ctx context.Context) (algoCnt, providerCnt int32, err error) // dashboard 展示算法数量和厂商数量
	GetAlgoProvider(ctx context.Context) ([]string, error)
	GetAlgoTypeName(ctx context.Context, provider string) ([]string, error)
	GetAlgoName(ctx context.Context, provider string, algoTypeName string) ([]string, error)
	GetAlgoIdByAlgoName(ctx context.Context, algoName string) (int32, error)
	GetAlgoIdsByAlgoNames(ctx context.Context, algoName []string) ([]int32, error) // 批量
	GetAllAlgoTypeName(ctx context.Context) ([]string, error)
	GetAlgoTypeIdByAlgoTypeName(ctx context.Context, algoTypeName string) (int32, error)
	GetAlgoInfoByAlgoName(ctx context.Context, algoName string) (models.TbAlgoInfo, error)
}

type DefaultAlgoInfo struct {
	DB *gorm.DB
}

func NewAlgoInfoRepo(conn *gorm.DB) AlgoInfoRepo {
	return &DefaultAlgoInfo{
		DB: conn,
	}
}

// GetAllAlgoInfo 根据算法类型获取所有算法的数量信息
func (d *DefaultAlgoInfo) GetAllAlgoInfo(ctx context.Context) ([]*AllAlgoInfo, error) {
	var infos []*AllAlgoInfo
	result := d.DB.Model(&models.TbAlgoInfo{}).Select(" algo_type, any_value(algo_type_name) as algo_type_name, count(algo_type) as count").
		Group("algo_type").Find(&infos)
	if result.Error != nil {
		return nil, result.Error
	}
	return infos, nil
}

// GetAlgoBase 取所有算法基本数据
func (d *DefaultAlgoInfo) GetAlgoBase(ctx context.Context) ([]*models.TbAlgoInfo, error) {
	var infos []*models.TbAlgoInfo
	result := d.DB.Model(&models.TbAlgoInfo{}).Select("algo_id,algo_name,algo_type,algo_type_name,provider,algo_status").
		Where("1=1").Find(&infos)
	if result.Error != nil && result.RowsAffected != 0 {
		return nil, result.Error
	}
	return infos, nil
}

func (d *DefaultAlgoInfo) GetAlgoBaseById(ctx context.Context, id uint32) ([]*models.TbAlgoInfo, error) {
	var infos []*models.TbAlgoInfo
	result := d.DB.Select("algo_id,algo_name,algo_type,provider").Where("algo_id = ?", id).Find(&infos)
	if result.Error != nil && result.RowsAffected != 0 {
		return nil, result.Error
	}
	return infos, nil
}

func (d *DefaultAlgoInfo) UpdateAlgoBaseById(ctx context.Context, id uint32, perf *pb.AlgoInfoPerf) error {
	var status int
	if perf.GetAlgorithmStatus() == 3 {
		status = 1
	}
	result := d.DB.Model(models.TbAlgoInfo{}).Where("algo_id = ?", id).
		Updates(models.TbAlgoInfo{
			AlgoName: tools.RMu0000(perf.GetAlgoName()),
			AlgoType: int(perf.GetAlgorithmType()),
			//AlgoTypeName: tools.RMu0000(perf.GetAlgoName()),
			Provider:   tools.RMu0000(perf.GetProviderName()),
			AlgoStatus: status,
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (d *DefaultAlgoInfo) CreateAlgoBaseInfo(ctx context.Context, perf *pb.AlgoInfoPerf) error {
	var status int
	if perf.GetAlgorithmStatus() == 3 {
		status = 1
	}
	var algoTypeName string
	if perf.GetAlgorithmType() == 1 {
		algoTypeName = "日内回转"
	} else if perf.GetAlgorithmType() == 2 {
		algoTypeName = "智能委托"
	}
	info := models.TbAlgoInfo{
		AlgoId:       int(perf.GetId()),
		AlgoName:     tools.RMu0000(perf.GetAlgoName()),
		AlgoType:     int(perf.GetAlgorithmType()),
		AlgoTypeName: algoTypeName,
		Provider:     tools.RMu0000(perf.GetProviderName()),
		AlgoStatus:   status,
	}
	result := d.DB.Create(&info)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultAlgoInfo) GetAlgoSummary(ctx context.Context) (algoCnt, providerCnt int32, err error) {
	type Count struct {
		AlgoCnt     int32
		ProviderCnt int32
	}
	var cnt Count
	result := d.DB.Model(&models.TbAlgoInfo{}).Select("count(distinct(algo_id)) as algo_cnt, count(distinct(provider)) as provider_cnt ").
		Where("algo_status=1").First(&cnt)
	if result.Error != nil {
		return 0, 0, result.Error
	}
	return cnt.AlgoCnt, cnt.ProviderCnt, nil
}

// GetAlgoProvider 算法条件筛选，拉取所有厂商
func (d *DefaultAlgoInfo) GetAlgoProvider(ctx context.Context) ([]string, error) {
	var out []string
	result := d.DB.Model(&models.TbAlgoInfo{}).Select("distinct(provider)").Where("algo_status=1").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	var ret []string
	for _, v := range out {
		if len(v) == 0 {
			continue
		}
		ret = append(ret, tools.RMu0000(v))
	}
	return ret, nil
}

// GetAlgoTypeName 算法条件筛选， 根据算法厂商拉取所有算法类型
func (d *DefaultAlgoInfo) GetAlgoTypeName(ctx context.Context, provider string) ([]string, error) {
	var out []string
	result := d.DB.Model(&models.TbAlgoInfo{}).Select("distinct(algo_type_name)")
	if provider != "" {
		result.Where("provider=?", provider)
	}
	result.Where("algo_status=1").Find(&out)

	if result.Error != nil {
		return nil, result.Error
	}
	var ret []string
	for _, v := range out {
		ret = append(ret, tools.RMu0000(v))
	}
	return ret, nil
}

// GetAlgoName 算法条件筛选, 根据算法厂商， 算法类型拉取所有算法
func (d *DefaultAlgoInfo) GetAlgoName(ctx context.Context, provider string, algoTypeName string) ([]string, error) {
	var out []string
	result := d.DB.Model(&models.TbAlgoInfo{}).Select("algo_name")
	if provider != "" {
		result.Where("provider=?", provider)
	}
	if algoTypeName != "" {
		result.Where("algo_type_name=?", algoTypeName)
	}
	result.Where("algo_status=1").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	var ret []string
	for _, v := range out {
		ret = append(ret, tools.RMu0000(v))
	}
	return ret, nil
}

// GetAlgoIdByAlgoName 根据算法名称反查算法ID
func (d *DefaultAlgoInfo) GetAlgoIdByAlgoName(ctx context.Context, algoName string) (int32, error) {
	var out int32
	result := d.DB.Model(models.TbAlgoInfo{}).Select("algo_id").Where("algo_name=?", algoName).Find(&out)
	if result.Error != nil {
		return 0, result.Error
	}
	return out, nil
}

// GetAlgoIdsByAlgoNames 批量根据算法名称反查算法ID     对比分析场景用
func (d *DefaultAlgoInfo) GetAlgoIdsByAlgoNames(ctx context.Context, algoName []string) ([]int32, error) { // 批量
	var out []int32
	result := d.DB.Model(models.TbAlgoInfo{}).Select("algo_id").Where("algo_name in ?", algoName).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

// GetAllAlgoTypeName 无条件查询所有的算法类型名称
func (d *DefaultAlgoInfo) GetAllAlgoTypeName(ctx context.Context) ([]string, error) {
	var out []models.TbAlgoInfo
	result := d.DB.Model(&models.TbAlgoInfo{}).Select("distinct(algo_type) as algo_type").Where(" algo_status=1").
		Order("algo_type").Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	var name []string
	for _, v := range out {
		if v.AlgoType == 1 {
			name = append(name, "日内回转")
		} else if v.AlgoType == 2 {
			name = append(name, "智能委托")
		}
		//name = append(name, tools.RMu0000(v.AlgoTypeName))
	}
	return name, nil
}

// GetAlgoTypeIdByAlgoTypeName 根据算法类型名称 反查算法类型ID
func (d *DefaultAlgoInfo) GetAlgoTypeIdByAlgoTypeName(ctx context.Context, algoTypeName string) (int32, error) {
	var out int32
	result := d.DB.Model(models.TbAlgoInfo{}).Select("distinct(algo_type)").Where("algo_type_name=?", algoTypeName).Find(&out)
	if result.Error != nil {
		return 0, result.Error
	}
	return out, nil
}

// GetAlgoInfoByAlgoName 根据算法名称查询返回其对应的厂商，算法类型
func (d *DefaultAlgoInfo) GetAlgoInfoByAlgoName(ctx context.Context, algoName string) (models.TbAlgoInfo, error) {
	var out models.TbAlgoInfo
	result := d.DB.Model(models.TbAlgoInfo{}).Select("algo_id,algo_name,provider,algo_type_name").Where("algo_name=?", algoName).
		Find(&out)
	if result.Error != nil {
		return models.TbAlgoInfo{}, result.Error
	}
	return out, nil
}
