// Package repo
/*
 Author: hawrkchen
 Date: 2022/7/6 16:07
 Desc:
*/
package repo

import (
	"algo_assess/assess-mq-server/proto"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/models"
	"context"
	"gorm.io/gorm"
	"strings"
	"time"
)

type SecurityInfoRepo interface {
	GetSecurityInfos(ctx context.Context) ([]*models.TbSecurityInfo, error)
	GetSecurityInfoById(ctx context.Context, secId string) ([]*models.TbSecurityInfo, error)
	UpdateSecurityById(ctx context.Context, secId string, perf *pb.SecurityInfoPerf) error
	CreateSecurityInfo(ctx context.Context, perf *pb.SecurityInfoPerf) error
	GetRanking(ctx context.Context, page, limit int) ([]*models.TbSecurityInfo, int64, error)
	GetSecurityList(ctx context.Context, secId string, page, limit int) ([]*models.TbSecurityInfo, int64, error)
	ModifySecurityProperty(ctx context.Context, secId string, fundType, stockType int32, industry string, liq int32) error
	AddSecurity(ctx context.Context, in *proto.SecurityUpdate) error
	DelSecurityProperty(ctx context.Context, secId string) error
	ImportSecurityUpdate(ctx context.Context, info *proto.SecurityInfo) error
	ImportSecurityCreate(ctx context.Context, info *proto.SecurityInfo) error
}

type DefaultSecurityInfo struct {
	DB *gorm.DB
}

func NewSecurityInfo(conn *gorm.DB) SecurityInfoRepo {
	return &DefaultSecurityInfo{
		DB: conn,
	}
}

// GetSecurityInfos   server start load data
func (s *DefaultSecurityInfo) GetSecurityInfos(ctx context.Context) ([]*models.TbSecurityInfo, error) {
	var infos []*models.TbSecurityInfo
	result := s.DB.Model(models.TbSecurityInfo{}).Select("security_id,security_source,security_name,pre_close_px,status," +
		"upper_limit_price,lower_limit_price, fund_type, stock_type,liquidity,industry").Where("status=48").Find(&infos)
	if result.Error != nil && result.RowsAffected != 0 {
		return nil, result.Error
	}
	return infos, nil
}

// GetSecurityInfoById 只是用来判断该记录是否存在
func (s *DefaultSecurityInfo) GetSecurityInfoById(ctx context.Context, secId string) ([]*models.TbSecurityInfo, error) {
	var info []*models.TbSecurityInfo
	result := s.DB.Select("security_id").Where("security_id = ?", secId).Find(&info)
	if result.Error != nil {
		return nil, result.Error
	}
	return info, nil
}

func (s *DefaultSecurityInfo) UpdateSecurityById(ctx context.Context, secId string, perf *pb.SecurityInfoPerf) error {
	result := s.DB.Model(models.TbSecurityInfo{}).Where("security_id", secId).
		Updates(models.TbSecurityInfo{
			SecuritySource:  perf.GetSecurityIdSource(),
			SecurityName:    perf.GetSecurityName(),
			PreClosePx:      perf.GetPrevClosePx(),
			Status:          int(perf.GetSecurityStatus()),
			IsPriceLimit:    int(perf.GetHasPriceLimit()),
			LimtType:        int(perf.GetLimitType()),
			Property:        int(perf.GetProperty()),
			UpperLimitPrice: int64(perf.GetUpperLimitPrice()),
			LowerLimitPrice: int64(perf.GetLowerLimitPrice()),
			UpdateTime:      time.Now(),
		})
	return result.Error
}

func (s *DefaultSecurityInfo) CreateSecurityInfo(ctx context.Context, perf *pb.SecurityInfoPerf) error {
	info := models.TbSecurityInfo{
		SecurityId:      perf.GetSecurityId(),
		SecuritySource:  perf.GetSecurityIdSource(),
		SecurityName:    perf.GetSecurityName(),
		PreClosePx:      perf.GetPrevClosePx(),
		Status:          int(perf.GetSecurityStatus()),
		IsPriceLimit:    int(perf.GetHasPriceLimit()),
		LimtType:        int(perf.GetLimitType()),
		Property:        int(perf.GetProperty()),
		UpperLimitPrice: int64(perf.GetUpperLimitPrice()),
		LowerLimitPrice: int64(perf.GetLowerLimitPrice()),
	}
	result := s.DB.Create(&info)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (s *DefaultSecurityInfo) GetRanking(ctx context.Context, page, limit int) ([]*models.TbSecurityInfo, int64, error) {
	var info []*models.TbSecurityInfo
	var count int64
	result := s.DB.Select("security_id,security_name").Where("1=1").Order("pre_close_px desc").
		Offset((page - 1) * limit).Limit(limit).Find(&info).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return info, count, nil
}

// GetSecurityList 配置列表 展示证券信息
func (s *DefaultSecurityInfo) GetSecurityList(ctx context.Context, secId string, page, limit int) ([]*models.TbSecurityInfo, int64, error) {
	var infos []*models.TbSecurityInfo
	var count int64
	result := s.DB.Model(models.TbSecurityInfo{}).Select("id,security_id,security_name,status," +
		"fund_type, stock_type,liquidity,industry,create_time,update_time")
	if secId != "" {
		result = result.Where("security_id like ?", "%"+secId+"%")
	}
	result = result.Offset((page - 1) * limit).Limit(limit).Find(&infos).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil && result.RowsAffected != 0 {
		return nil, 0, result.Error
	}
	return infos, count, nil
}

// ModifySecurityProperty 配置菜单，修改证券属性
func (s *DefaultSecurityInfo) ModifySecurityProperty(ctx context.Context, secId string, fundType, stockType int32, industry string, liq int32) error {
	err := s.DB.Model(models.TbSecurityInfo{}).Where("security_id=?", secId).
		Updates(map[string]interface{}{"fund_type": fundType, "stock_type": stockType, "industry": industry, "liquidity": liq, "update_time": time.Now()}).Error
	if err != nil {
		return err
	}
	return nil
}

// AddSecurity 绩效平台的新增
func (s *DefaultSecurityInfo) AddSecurity(ctx context.Context, in *proto.SecurityUpdate) error {
	info := models.TbSecurityInfo{
		SecurityId:   in.SecId,
		SecurityName: in.SecName,
		Status:       48,
		FundType:     int(in.FundType),
		StockType:    int(in.StockType),
		Liquidity:    int(in.Liquidity),
		Industry:     in.Industry,
	}
	result := s.DB.Create(&info)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

// DelSecurityProperty 绩效平台删除
func (s *DefaultSecurityInfo) DelSecurityProperty(ctx context.Context, secId string) error {
	result := s.DB.Model(models.TbSecurityInfo{}).Where("security_id=?", secId).
		Updates(models.TbSecurityInfo{
			Status:     2,
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (s *DefaultSecurityInfo) ImportSecurityUpdate(ctx context.Context, info *proto.SecurityInfo) error {
	result := s.DB.Model(models.TbSecurityInfo{}).Where("security_id=?", strings.TrimSpace(info.GetSecId())).
		Updates(models.TbSecurityInfo{
			SecurityName: strings.TrimSpace(info.GetSecName()),
			Status:       48,
			FundType:     int(info.GetFundType()),
			StockType:    int(info.GetStockType()),
			Liquidity:    int(info.GetLiquidity()),
			Industry:     info.GetIndustry(),
			UpdateTime:   time.Now(),
		})
	return result.Error
}

func (s *DefaultSecurityInfo) ImportSecurityCreate(ctx context.Context, info *proto.SecurityInfo) error {
	e := models.TbSecurityInfo{
		SecurityId:   strings.TrimSpace(info.GetSecId()),
		SecurityName: strings.TrimSpace(info.GetSecName()),
		Status:       48,
		FundType:     int(info.GetFundType()),
		StockType:    int(info.GetStockType()),
		Liquidity:    int(info.GetLiquidity()),
		Industry:     info.GetIndustry(),
	}
	result := s.DB.Create(&e)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}
