package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SecurityListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSecurityListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SecurityListLogic {
	return &SecurityListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SecurityListLogic) SecurityList(req *types.ListSecurityReq) (resp *types.ListSecurityRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("get req:%+v", req)
	sReq := &mqservice.SecurityListReq{
		SecId: req.SecId,
		Page:  req.Page,
		Limit: req.Limit,
	}
	rsp, err := l.svcCtx.AssessMQClient.SecurityList(l.ctx, sReq)
	if err != nil {
		l.Logger.Error("rpc call SecurityList error:", err)
		return &types.ListSecurityRsp{
			Code:  210,
			Msg:   err.Error(),
			Total: 0,
			Infos: nil,
		}, nil
	}
	var list []types.SecurityInfo
	for _, v := range rsp.GetInfos() {
		l := types.SecurityInfo{
			Id:         v.GetId(),
			SecId:      v.GetSecId(),
			SecName:    v.GetSecName(),
			Status:     v.GetStatus(),
			FundType:   v.GetFundType(),
			StockType:  v.GetStockType(),
			Liquidity:  v.GetLiquidity(),
			Industry:   v.GetIndustry(),
			UpdateTime: v.GetUpdateTime(),
		}
		list = append(list, l)
	}
	return &types.ListSecurityRsp{
		Code:  200,
		Msg:   "succes",
		Total: rsp.GetTotal(),
		Infos: list,
	}, nil

}
