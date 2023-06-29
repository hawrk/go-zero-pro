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

type UserPositionRepo interface {
	GetUserPosition(ctx context.Context, id int64) ([]*UserPosition, error)
}

type DefaultUserPosition struct {
	DB *gorm.DB
}

func NewUserPosition(conn *gorm.DB) UserPositionRepo {
	return &DefaultUserPosition{
		DB: conn,
	}
}

type UserPosition struct {
	SecurityId      string `json:"security_id"`
	SecurityName    string `json:"security_name"`
	PositionQty     int64  `json:"position_qty"`
	OriginOpenPrice int64  `json:"origin_open_price"`
	AvgPrice        int64  `json:"avg_price"`
}

// GetUserPosition 联表 tb_security_info
func (d *DefaultUserPosition) GetUserPosition(ctx context.Context, id int64) ([]*UserPosition, error) {
	var out []*UserPosition
	result := d.DB.Model(busmodels.TbUserPosition{}).Select("tb_user_position.security_id,b.security_name, "+
		"origin_open_price, position_qty,avg_price").
		Joins("inner join tb_security_info b on uuser_id = ? and tb_user_position.security_id = b.security_id ", id).
		Find(&out)
	if result.Error != nil {
		return nil, result.Error
	}
	return out, nil
}
