// Package repo
/*
 Author: hawrkchen
 Date: 2022/10/12 15:12
 Desc:
*/
package repo

import (
	"account-auth/account-auth-server/internal/types"
	"account-auth/account-auth-server/models"
	"account-auth/account-auth-server/pkg/tools"
	"context"
	"gorm.io/gorm"
	"time"
)

type AuthUserRepo interface {
	CreateAuthUser(ctx context.Context, u *types.UserModfiyReq) error
	GetUserList(ctx context.Context, userName string, chanType int, page, limit int) ([]*models.TbAuthUser, int64, error)
	DelAuthUser(ctx context.Context, userId string, chanType int) error
	UpdateAuthUser(ctx context.Context, infos *types.UserModfiyReq) error
	UpdateAuthUserPasswd(ctx context.Context, infos *types.UserModfiyReq) error
	CheckUserPasswd(ctx context.Context, userId string, chanType int) (string, error)
	GetAuthUserById(ctx context.Context, userId string, chanType int, status int) (models.TbAuthUser, error)
	GetUserByRoleId(ctx context.Context, roleId int, chanType int) ([]*models.TbAuthUser, error)
	UpdateExpireTime(ctx context.Context, userId string, chanType int) error
}

type DefaultAuthUser struct {
	DB *gorm.DB
}

func NewAuthUserRepo(conn *gorm.DB) AuthUserRepo {
	return &DefaultAuthUser{
		DB: conn,
	}
}

// CreateAuthUser 新增绩效平台登陆用户
func (d *DefaultAuthUser) CreateAuthUser(ctx context.Context, u *types.UserModfiyReq) error {
	in := models.TbAuthUser{
		UserId:     u.UserId,
		UserName:   u.UserName,
		UserPasswd: u.Password,
		RoleId:     int(u.RoleId),
		RoleName:   u.RoleName,
		UserType:   int(u.UserType), // 该UserType 在绩效平台新增时，只区分普通用户(0)和管理员(1)，从总线平台同步过来的统一为0
		ChanType:   u.ChanType,
		Status:     1,
		//ExpireTime: time.Now().AddDate(0, 0, 180).Unix(),
		ExpireTime: tools.QuarterExpireTime(),
	}
	result := d.DB.Create(&in)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

//
func (d *DefaultAuthUser) GetUserList(ctx context.Context, userName string, chanType int, page, limit int) ([]*models.TbAuthUser, int64, error) {
	var out []*models.TbAuthUser
	var count int64
	// modify 只返回 正常和禁用的
	result := d.DB.Model(&models.TbAuthUser{}).Select("user_id,user_name,role_id,role_name,user_type,status,create_time").
		Where("chan_type=? and status in(1,3)", chanType).Order("create_time")
	if userName != "" {
		result = result.Where("user_name like ?", "%"+userName+"%")
	}
	result = result.Offset((page - 1) * limit).Limit(limit).Find(&out).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return out, count, nil
}

func (d *DefaultAuthUser) DelAuthUser(ctx context.Context, userId string, chanType int) error {
	result := d.DB.Model(models.TbAuthUser{}).Where("user_id=? and chan_type=?", userId, chanType).
		Updates(models.TbAuthUser{
			Status:     2,
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (d *DefaultAuthUser) UpdateAuthUser(ctx context.Context, infos *types.UserModfiyReq) error {
	result := d.DB.Model(models.TbAuthUser{}).Where("user_id=? and chan_type=?", infos.UserId, infos.ChanType).
		Updates(models.TbAuthUser{
			UserName:   infos.UserName,
			UserPasswd: infos.Password,
			RoleId:     int(infos.RoleId),
			RoleName:   infos.RoleName,
			UserType:   int(infos.UserType),
			Status:     infos.Status,
			UpdateTime: time.Now(),
		})
	return result.Error
}

// UpdateAuthUserPasswd 密码有改动时，需要同步更新FirstLogin，expireTime字段
func (d *DefaultAuthUser) UpdateAuthUserPasswd(ctx context.Context, infos *types.UserModfiyReq) error {
	result := d.DB.Model(models.TbAuthUser{}).Where("user_id=? and chan_type=?", infos.UserId, infos.ChanType).
		Updates(models.TbAuthUser{
			UserName:   infos.UserName,
			UserPasswd: infos.Password,
			RoleId:     int(infos.RoleId),
			RoleName:   infos.RoleName,
			UserType:   int(infos.UserType),
			Status:     infos.Status,
			FirstLogin: 1,
			//ExpireTime: time.Now().AddDate(0, 0, 180).Unix(),
			ExpireTime: tools.QuarterExpireTime(),
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (d *DefaultAuthUser) CheckUserPasswd(ctx context.Context, userId string, chanType int) (string, error) {
	var out string
	result := d.DB.Model(models.TbAuthUser{}).Select("user_passwd").Where("user_id=? and chan_type=?", userId, chanType).Find(&out)
	if result.Error != nil {
		return "", result.Error
	}
	return out, nil
}

// GetAuthUserById 取用户信息， 密码校验，并返回用户类型
func (d *DefaultAuthUser) GetAuthUserById(ctx context.Context, userId string, chanType int, status int) (models.TbAuthUser, error) {
	var out models.TbAuthUser
	result := d.DB.Model(models.TbAuthUser{}).Select("user_id,user_name,user_passwd,role_id,role_name, user_type,status,"+
		"first_login,expire_time").
		Where("user_id=? and chan_type =?", userId, chanType)
	if status != 0 {
		result.Where("status=?", status)
	}
	result.Find(&out)
	if result.Error != nil {
		return models.TbAuthUser{}, result.Error
	}
	return out, nil
}

// GetUserByRoleId 根据RoleId 查用户列表
func (d *DefaultAuthUser) GetUserByRoleId(ctx context.Context, roleId int, chanType int) ([]*models.TbAuthUser, error) {
	var out []*models.TbAuthUser
	result := d.DB.Model(models.TbAuthUser{}).Select("user_id").Where("role_id=? and chan_type =? and status =1", roleId, chanType).
		Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultAuthUser) UpdateExpireTime(ctx context.Context, userId string, chanType int) error {
	result := d.DB.Model(models.TbAuthUser{}).Where("user_id=? and chan_type=?", userId, chanType).
		Updates(models.TbAuthUser{
			//ExpireTime: time.Now().AddDate(0, 0, 180).Unix(),
			ExpireTime: tools.QuarterExpireTime(),
			UpdateTime: time.Now(),
		})
	return result.Error
}
