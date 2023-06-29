package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"
	"encoding/xml"
	"net/http"

	"algo_assess/assess-api-server/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExportSecurityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportSecurityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportSecurityLogic {
	return &ExportSecurityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportSecurityLogic) ExportSecurity(w http.ResponseWriter) error {
	// todo: add your logic here and delete this line
	l.Logger.Info("into ExportSecurity......")
	rsp, err := l.svcCtx.AssessMQClient.ExportSecurityInfo(l.ctx, &mqservice.ExportSecurityReq{})
	if err != nil {
		l.Logger.Error("rpc call ExportSecurityInfo error:", err)
		return err
	}
	list := make([]*SecurityData, 0, len(rsp.GetInfos()))
	for _, v := range rsp.GetInfos() {
		s := SecurityData{
			SecurityId:              v.GetSecId(),
			SecuritySource:          "",
			SecurityName:            v.GetSecName(),
			PreClosePx:              0,
			BuyQtyUpperLimit:        0,
			SellQtyUpperLimit:       0,
			MarketBuyQtyUpperLimit:  0,
			MarketSellQtyUpperLimit: 0,
			SecurityStatus:          v.GetStatus(),
			HasPriceLimit:           0,
			LimitType:               0,
			Property:                0,
			UpperLimitPrice:         0,
			LowerLimitPrice:         0,
			BuyQtyUnit:              0,
			SellQtyUnit:             0,
			MarketBuyQtyUnit:        0,
			MarketSellQtyUnit:       0,
			FundType:                v.GetFundType(),
			StockType:               v.GetStockType(),
			Liquidity:               v.GetLiquidity(),
			Industry:                v.GetIndustry(),
		}
		list = append(list, &s)
	}
	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "filename="+"securityInfo.xml")
	output, _ := xml.MarshalIndent(list, " ", " ")
	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(output)
	return nil
}
