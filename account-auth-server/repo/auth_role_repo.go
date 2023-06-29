// Package repo
/*
 Author: hawrkchen
 Date: 2022/10/13 10:19
 Desc:
*/
package repo

import (
	"account-auth/account-auth-server/global"
	"account-auth/account-auth-server/internal/types"
	"account-auth/account-auth-server/models"
	"context"
	"gorm.io/gorm"
	"time"
)

type AuthRoleRepo interface {
	CreateAuthRole(ctx context.Context, req *types.RoleModifyReq) error
	GetRoleList(ctx context.Context, roleId int, roleName string, chanType int, page, limit int) ([]*models.TbAuthRole, int64, error)
	DelAuthRole(ctx context.Context, roleId int) error
	UpdateAuthRole(ctx context.Context, req *types.RoleModifyReq) error
	GetEffectRoleList(ctx context.Context, chanType int) ([]*models.TbAuthRole, error)
	GetAuthByRoleId(ctx context.Context, roleId int) (models.TbAuthRole, error)
	CheckExistRole(ctx context.Context, chType int, roleId int, roleName string) (int, string, error)
}

type DefaultAuthRole struct {
	DB *gorm.DB
}

func NewAuthRole(conn *gorm.DB) AuthRoleRepo {
	return &DefaultAuthRole{
		DB: conn,
	}
}

func (d *DefaultAuthRole) CreateAuthRole(ctx context.Context, req *types.RoleModifyReq) error {
	in := models.TbAuthRole{
		RoleId:   uint(req.RoleId),
		RoleName: req.RoleName,
		RoleAuth: req.RoleAuth,
		ChanType: req.ChanType,
		Status:   1,
	}
	result := d.DB.Create(&in)
	if result.Error != nil || result.RowsAffected != 1 {
		return result.Error
	}
	return nil
}

func (d *DefaultAuthRole) GetRoleList(ctx context.Context, roleId int, roleName string, chanType int, page, limit int) ([]*models.TbAuthRole, int64, error) {
	var out []*models.TbAuthRole
	var count int64
	// modify 只取状态正常的
	result := d.DB.Model(&models.TbAuthRole{}).Select("role_id,role_name,role_auth,status,create_time").
		Where("chan_type=? and status=?", chanType, global.RoleStatusNormal)
	if roleId != 0 {
		result = result.Where("role_id=?", roleId)
	}
	if roleName != "" {
		result = result.Where("role_name like ?", "%"+roleName+"%")
	}
	//result = result.Order("status").Order("create_time desc")
	result = result.Order("create_time desc")
	result = result.Offset((page - 1) * limit).Limit(limit).Find(&out).Offset(-1).Limit(-1).Count(&count)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return out, count, nil
}

func (d *DefaultAuthRole) DelAuthRole(ctx context.Context, roleId int) error {
	result := d.DB.Model(models.TbAuthRole{}).Where("role_id=?", roleId).
		Updates(models.TbAuthRole{
			Status:     2,
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (d *DefaultAuthRole) UpdateAuthRole(ctx context.Context, req *types.RoleModifyReq) error {
	result := d.DB.Model(models.TbAuthRole{}).Where("role_id=?", req.RoleId).
		Updates(models.TbAuthRole{
			RoleName:   req.RoleName,
			RoleAuth:   req.RoleAuth,
			UpdateTime: time.Now(),
		})
	return result.Error
}

func (d *DefaultAuthRole) GetEffectRoleList(ctx context.Context, chanType int) ([]*models.TbAuthRole, error) {
	var out []*models.TbAuthRole
	result := d.DB.Model(&models.TbAuthRole{}).Select("id,role_id,role_name,role_auth,status,create_time").
		Where("chan_type=? and status=1", chanType).Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}

func (d *DefaultAuthRole) GetAuthByRoleId(ctx context.Context, roleId int) (models.TbAuthRole, error) {
	var out models.TbAuthRole
	result := d.DB.Model(&models.TbAuthRole{}).Select("role_auth").Where("role_id=?", roleId).Find(&out)
	if result.Error != nil {
		return models.TbAuthRole{}, result.Error
	}
	return out, nil
}

// CheckExistRole 检查 role_id 或role_name 是否存在
func (d *DefaultAuthRole) CheckExistRole(ctx context.Context, chType int, roleId int, roleName string) (int, string, error) {
	var out models.TbAuthRole
	result := d.DB.Model(&models.TbAuthRole{}).Select("role_id, role_name")
	if roleId != 0 {
		result.Where("role_id = ?", roleId)
	}
	if roleName != "" {
		result.Where("role_name=?", roleName)
	}
	result.Where("chan_type=? and status = 1", chType).Find(&out)
	if result.Error != nil {
		return 0, "", result.Error
	}
	return int(out.RoleId), out.RoleName, nil
}
