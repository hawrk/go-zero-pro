package logic

import (
	"context"
	"strings"

	"algo_assess/mornano-rpc-server/internal/svc"
	"algo_assess/mornano-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserCapitalLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserCapitalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserCapitalLogic {
	return &GetUserCapitalLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserCapital 查询用户资金和证券持仓信息  (用户画像)
func (l *GetUserCapitalLogic) GetUserCapital(in *proto.CapitalReq) (*proto.CapitailRsp, error) {
	l.Logger.Infof("in GetUserCapital, get req:%+v", in)
	if in.GetUserId() == "" {
		l.Logger.Error("user_id not set...")
		return &proto.CapitailRsp{
			Code:          351,
			Msg:           "user_id not set",
			Available:     0,
			StockPosition: nil,
		}, nil
	}
	// 根据user_id 到 user_info表找出ID
	uuserId, err := l.svcCtx.UserInfoRepo.GetIdByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("GetIdByUserId error:", err)
		return &proto.CapitailRsp{
			Code:          352,
			Msg:           err.Error(),
			Available:     0,
			StockPosition: nil,
		}, nil
	}
	// 取可用资金
	// 1. 根据user_id 从user_info 表中找出对应的id
	asset, err := l.svcCtx.AssetRepo.GetAssetInfoById(l.ctx, uuserId)
	if err != nil {
		l.Logger.Error("GetAssetInfoById error:", err)
	}
	// 资金账户算可用资金
	available := asset.Balance - asset.Frozen
	l.Logger.Info("get available:", available)

	// 2. 取持仓信息
	position, err := l.svcCtx.UserPositionRepo.GetUserPosition(l.ctx, uuserId)
	if err != nil {
		l.Logger.Error("GetUserPosition error:", err)
	}
	var us []*proto.StockPosition
	for _, v := range position {
		mc := (float64(v.AvgPrice) / 10000) * (float64(v.PositionQty) / 100)          // 市值= 持仓数量*最新价
		cost := (float64(v.OriginOpenPrice) / 10000) * (float64(v.PositionQty) / 100) // 成本价 = 持仓数量*平均开仓价
		if mc == 0.00 && cost == 0.00 {
			continue
		}
		u := &proto.StockPosition{
			SecId:     strings.TrimSpace(v.SecurityId),
			SecName:   strings.TrimSpace(v.SecurityName),
			MarketCap: mc,
			Cost:      cost,
		}
		us = append(us, u)
	}
	return &proto.CapitailRsp{
		Code:          200,
		Msg:           "success",
		Available:     available,
		StockPosition: us,
	}, nil

}
