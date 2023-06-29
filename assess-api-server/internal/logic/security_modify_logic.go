package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type SecurityModifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSecurityModifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SecurityModifyLogic {
	return &SecurityModifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SecurityModifyLogic) SecurityModify(req *types.ModifySecurityReq) (resp *types.ModifySecurityRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("get SecurityModify req:%+v", req)
	/*
		if err := l.CheckParams(req); err != nil {
			return &types.ModifySecurityRsp{
				Code:   250,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}
	*/

	var ls []*mqservice.SecurityUpdate
	for _, v := range req.Lists {
		l := &mqservice.SecurityUpdate{
			SecId:     v.SecId,
			SecName:   v.SecName,
			FundType:  v.FundType,
			StockType: v.StockType,
			Liquidity: v.Liquidity,
			Industry:  v.Industry,
		}
		ls = append(ls, l)
	}
	rsp, err := l.svcCtx.AssessMQClient.SecurityUpdate(l.ctx, &mqservice.SecurityModifyReq{
		OperType: req.OperType,
		List:     ls,
	})
	if err != nil {
		l.Logger.Error("rpc call SecurityUpdate error:", err)
		return &types.ModifySecurityRsp{
			Code:   220,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	return &types.ModifySecurityRsp{
		Code:   200,
		Msg:    rsp.GetMsg(),
		Result: rsp.GetResult(),
	}, nil
}

func (l *SecurityModifyLogic) CheckParams(req *types.ModifySecurityReq) error {
	if req.OperType == 1 { // 新增 必须保证参数完整
		for _, v := range req.Lists {
			if v.SecName == "" || v.FundType == 0 || v.StockType == 0 {
				return errors.New("field sec_name/fund_type/stock_type not set")
			}
		}
	} else if req.OperType == 2 { // 修改
		//for _, v := range req.Lists {
		//	if v.FundType == 0 && v.StockType == 0 {  // 增加行业属性和流动性
		//		return errors.New("fund_type or stock_type must be set")
		//	}
		//}
	}
	return nil
}
