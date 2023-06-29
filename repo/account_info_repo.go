// Package repo
/*
 Author: hawrkchen
 Date: 2022/6/24 13:50
 Desc:
*/
package repo

import (
	"algo_assess/assess-mq-server/proto"
	pb "algo_assess/assess-mq-server/proto/order"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"
	"gorm.io/gorm"
	"time"
)

type AccountInfoRepo interface {
	GetAccountInfos(ctx context.Context) ([]*models.TbAccountInfo, error)
	GetAccountInfoCounts(ctx context.Context) (int64, error)
	GetAccountInfoById(ctx context.Context, id int) ([]*models.TbAccountInfo, error)
	GetAccountInfoByUserId(ctx context.Context, userId string) (models.TbAccountInfo, error)
	UpdateAccountInfoById(ctx context.Context, id int, perf *pb.UserInfoPerf) error
	CreateAccountInfo(ctx context.Context, perf *pb.UserInfoPerf) error
	CheckLogin(ctx context.Context, userId string) ([]*models.TbAccountInfo, error)
	GetRanking(ctx context.Context, page, limit int) ([]*models.TbAccountInfo, int64, error)
	GetProviderAccountInfo(ctx context.Context) ([]*ProviderAccountInfo, error) // 这个需要联表查询
	GetAccountLists(ctx context.Context, userId string, page, limit int) ([]*models.TbAccountInfo, int64, error)
	ModifyUserProperty(ctx context.Context, userId string, userGrade string) error
	AddUser(ctx context.Context, userId, userName, grade string) error
	DelUser(ctx context.Context, userId string) error
	ImportUserUpdate(ctx context.Context, info *proto.UserInfo) error
	ImportUserCreate(ctx context.Context, info *proto.UserInfo) error
}

type DefaultAccountInfo struct {
	DB *gorm.DB
}

func NewAccountInfo(conn *gorm.DB) AccountInfoRepo {
	return &DefaultAccountInfo{
		DB: conn,
	}
}

func (d *DefaultAccountInfo) GetAccountInfos(ctx context.Context) ([]*models.TbAccountInfo, error) {
	var account []*models.TbAccountInfo
	result := d.DB.Select("account_id, user_id, user_name, user_type, par_user_id").
		Where("1=1").Find(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

// GetAccountInfoCounts 根据用户类型统计其总数量
func (d *DefaultAccountInfo) GetAccountInfoCounts(ctx context.Context) (int64, error) {
	var count int64
	result := d.DB.Model(models.TbAccountInfo{}).Select("count(*) as count ").Where("user_type=1").Find(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (d *DefaultAccountInfo) GetAccountInfoById(ctx context.Context, id int) ([]*models.TbAccountInfo, error) {
	var account []*models.TbAccountInfo
	result := d.DB.Select("account_id").Where("account_id = ?", id).Find(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

// GetAccountInfoByUserId 根据user_id 查询账户表信息
func (d *DefaultAccountInfo) GetAccountInfoByUserId(ctx context.Context, userId string) (models.TbAccountInfo, error) {
	var account models.TbAccountInfo
	result := d.DB.Select("account_id,user_name,user_type,user_grade").Where("user_id=?", userId).Find(&account)
	if result.Error != nil {
		return models.TbAccountInfo{}, result.Error
	}
	return account, nil
}

func (d *DefaultAccountInfo) UpdateAccountInfoById(ctx context.Context, id int, perf *pb.UserInfoPerf) error {
	result := d.DB.Model(models.TbAccountInfo{}).Where("account_id = ?", id).
		Updates(models.TbAccountInfo{
			//AccountId:  int(perf.GetId()),
			UserId:     tools.RMu0000(perf.GetUserId()),
			UserName:   tools.RMu0000(perf.GetUserName()),
			UserType:   int(perf.GetUserType()),
			RiskGroup:  int(perf.GetRiskGroup()),
			ParUserId:  perf.GetUuserId(),
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (d *DefaultAccountInfo) CreateAccountInfo(ctx context.Context, perf *pb.UserInfoPerf) error {
	info := models.TbAccountInfo{
		AccountId: int(perf.GetId()),
		UserId:    tools.RMu0000(perf.GetUserId()),
		UserName:  tools.RMu0000(perf.GetUserName()),
		UserType:  int(perf.GetUserType()),
		RiskGroup: int(perf.GetRiskGroup()),
		ParUserId: perf.GetUuserId(),
	}
	result := d.DB.Create(&info)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultAccountInfo) CheckLogin(ctx context.Context, userId string) ([]*models.TbAccountInfo, error) {
	var account []*models.TbAccountInfo
	result := d.DB.Select("account_id, user_passwd,user_type").Where("user_id = ?", userId).Find(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (d *DefaultAccountInfo) GetRanking(ctx context.Context, page, limit int) ([]*models.TbAccountInfo, int64, error) {
	var info []*models.TbAccountInfo
	var count int64
	result := d.DB.Select("user_id,user_name").Where("user_type = 1").Order("account_id").
		Offset((page - 1) * limit).Limit(limit).Find(&info).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return info, count, nil
}

func (d *DefaultAccountInfo) GetProviderAccountInfo(ctx context.Context) ([]*ProviderAccountInfo, error) {
	var info []*ProviderAccountInfo
	d.DB.Model(models.TbAccountInfo{}).Select("b.algo_id, b.provider ,user_id").Joins("inner join tb_algo_info b on user_name = b.provider").
		Scan(&info)
	return info, nil
}

// GetAccountLists 配置菜单,列出用户信息
func (d *DefaultAccountInfo) GetAccountLists(ctx context.Context, userId string, page, limit int) ([]*models.TbAccountInfo, int64, error) {
	var account []*models.TbAccountInfo
	var count int64
	result := d.DB.Select("account_id,user_id,user_name,user_type,user_grade,update_time")
	if userId != "" {
		result = result.Where("user_id like ?", "%"+userId+"%")
	}
	result = result.Offset((page - 1) * limit).Limit(limit).Find(&account).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil && result.RowsAffected != 0 {
		return nil, 0, result.Error
	}
	return account, count, nil
}

// ModifyUserProperty 配置菜单，修改用户级别
func (d *DefaultAccountInfo) ModifyUserProperty(ctx context.Context, userId string, userGrade string) error {
	err := d.DB.Model(models.TbAccountInfo{}).Where("user_id=?", userId).
		Updates(map[string]interface{}{"user_grade": userGrade, "update_time": time.Now()}).Error
	if err != nil {
		return err
	}
	return nil
}

// AddUser 绩效配置，新增用户
func (d *DefaultAccountInfo) AddUser(ctx context.Context, userId, userName, grade string) error {
	info := models.TbAccountInfo{
		UserId:     userId,
		UserName:   userName,
		UserGrade:  grade,
		GradeFixed: 1,
	}
	result := d.DB.Create(&info)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultAccountInfo) DelUser(ctx context.Context, userId string) error {
	// 不会真正删除，先设置user_type 为4
	err := d.DB.Model(models.TbAccountInfo{}).Where("user_id=?", userId).
		Updates(map[string]interface{}{"user_type": 4, "update_time": time.Now()}).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *DefaultAccountInfo) ImportUserUpdate(ctx context.Context, info *proto.UserInfo) error {
	result := d.DB.Model(models.TbAccountInfo{}).Where("user_id=?", info.UserId).
		Updates(models.TbAccountInfo{
			UserId:     info.UserId,
			UserName:   info.UserName,
			UserGrade:  info.UserGrade,
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (d *DefaultAccountInfo) ImportUserCreate(ctx context.Context, info *proto.UserInfo) error {
	t := models.TbAccountInfo{
		UserId:    info.GetUserId(),
		UserName:  info.GetUserName(),
		UserGrade: info.GetUserGrade(),
	}
	result := d.DB.Create(&t)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}
