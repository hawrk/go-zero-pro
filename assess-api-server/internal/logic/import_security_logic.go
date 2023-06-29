package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"algo_assess/pkg/tools"
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportSecurityLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportSecurityLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportSecurityLogic {
	return &ImportSecurityLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportSecurityLogic) ImportSecurity(r *http.Request) (resp *types.ImportSecurityRsp, err error) {
	// todo: add your logic here and delete this line
	f, header, err := r.FormFile("file")
	if err != nil {
		l.Logger.Error("parse file error:", err)
		return &types.ImportSecurityRsp{
			Code:   310,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	l.Logger.Info("ImportSecurity, parse file name:", header.Filename)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		l.Logger.Error("read file fail:", err)
		return &types.ImportSecurityRsp{
			Code:   320,
			Msg:    err.Error(),
			Result: 0,
		}, nil
	}
	// 解析文件内容
	var security XmlSecurity
	if err := xml.Unmarshal(b, &security); err != nil {
		l.Logger.Error("xml Unmarshal error:", err)
		return &types.ImportSecurityRsp{
			Code:   330,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	list := make([]*mqservice.SecurityInfo, 0, len(security.SecurityData))
	for _, v := range security.SecurityData {
		l := mqservice.SecurityInfo{
			Id:        0,
			SecId:     v.SecurityId,
			SecName:   tools.RMu0000(v.SecurityName),
			Status:    v.SecurityStatus,
			FundType:  v.FundType,
			StockType: v.StockType,
			Liquidity: v.Liquidity, // 流动性
			Industry:  v.Industry,  // 行业类型
		}
		list = append(list, &l)
	}
	// rpc
	rsp, err := l.svcCtx.AssessMQClient.ImportSecurityInfo(l.ctx, &mqservice.ImportSecurityReq{
		List: list,
	})
	if err != nil {
		l.Logger.Error("rpc call ImportSecurityInfo error:", err)
		return &types.ImportSecurityRsp{
			Code:   250,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}

	return &types.ImportSecurityRsp{
		Code:   200,
		Msg:    "success",
		Result: rsp.GetResult(),
	}, nil
}
