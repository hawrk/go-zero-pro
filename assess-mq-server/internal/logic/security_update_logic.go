package logic

import (
	"algo_assess/assess-mq-server/internal/dao"
	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type SecurityUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSecurityUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SecurityUpdateLogic {
	return &SecurityUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  配置： 证券属性修改
func (l *SecurityUpdateLogic) SecurityUpdate(in *proto.SecurityModifyReq) (*proto.SecurityModifyRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("SecurityUpdate, get req:", in)
	if len(in.GetList()) == 0 {
		l.Logger.Error("no data update")
		return &proto.SecurityModifyRsp{
			Code:   220,
			Msg:    errors.New("no data").Error(),
			Result: 2,
		}, nil
	}
	if in.GetOperType() == 1 {
		for _, v := range in.GetList() {
			if v.GetSecId() == "" {
				l.Logger.Error("security_id empty...")
				continue
			}
			if err := l.svcCtx.SecurityRepo.AddSecurity(l.ctx, v); err != nil {
				l.Logger.Error("error add security property :", err)
				continue
			}
			// 更新全局证券信息
			l.UpdateGlobalSecurity(v.GetSecId(), v.GetFundType(), v.GetStockType(), v.GetLiquidity(), v.GetIndustry())
		}
	} else if in.GetOperType() == 2 {
		for _, v := range in.GetList() {
			if v.GetSecId() == "" {
				l.Logger.Error("security_id empty...")
				continue
			}
			if err := l.svcCtx.SecurityRepo.ModifySecurityProperty(l.ctx, v.GetSecId(),
				v.GetFundType(), v.GetStockType(), v.GetIndustry(), v.GetLiquidity()); err != nil {
				l.Logger.Error("error modify security property :", err)
				continue
			}
			// 更新全局证券信息
			l.UpdateGlobalSecurity(v.GetSecId(), v.GetFundType(), v.GetStockType(), v.GetLiquidity(), v.GetIndustry())
		}
	} else if in.GetOperType() == 3 {
		for _, v := range in.GetList() {
			if v.GetSecId() == "" {
				l.Logger.Error("security_id empty...")
				continue
			}
			if err := l.svcCtx.SecurityRepo.DelSecurityProperty(l.ctx, v.GetSecId()); err != nil {
				l.Logger.Error("error del security property :", err)
				continue
			}
			// 直接删除全局变量的Key
			dao.GSecurityMap.RWMutex.Lock()
			delete(dao.GSecurityMap.SecurityBase, v.GetSecId())
			dao.GSecurityMap.RWMutex.Unlock()
		}
	}
	return &proto.SecurityModifyRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}

// UpdateGlobalSecurity 更新全局证券信息
func (l *SecurityUpdateLogic) UpdateGlobalSecurity(secId string, fundType, stockType, Liqui int32, industry string) {
	dao.GSecurityMap.RWMutex.Lock()
	s := dao.GSecurityMap.SecurityBase[secId]
	if fundType != 0 {
		s.FundType = int(fundType)
	}
	if stockType != 0 {
		s.StockType = int(stockType)
	}
	if Liqui != 0 {
		s.Liquidity = int(Liqui)
	}
	if industry != "" {
		s.Industry = industry
	}
	dao.GSecurityMap.SecurityBase[secId] = s
	dao.GSecurityMap.RWMutex.Unlock()
}
