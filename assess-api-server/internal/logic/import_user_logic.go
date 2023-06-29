package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ImportUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImportUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImportUserLogic {
	return &ImportUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ImportUserLogic) ImportUser(r *http.Request) (resp *types.ImportUserRsp, err error) {
	// todo: add your logic here and delete this line
	f, header, err := r.FormFile("file")
	if err != nil {
		l.Logger.Error("parse file error:", err)
		return &types.ImportUserRsp{
			Code:   310,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	l.Logger.Info("in ImportUser, get file name:", header.Filename)
	b, err := ioutil.ReadAll(f)
	if err != nil {
		l.Logger.Error("read file fail:", err)
		return &types.ImportUserRsp{
			Code:   320,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	// 解析文件内容
	var user XmlUserInfo
	if err := xml.Unmarshal(b, &user); err != nil {
		l.Logger.Error("xml Unmarshal error:", err)
		return &types.ImportUserRsp{
			Code:   330,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	list := make([]*mqservice.UserInfo, 0, len(user.UserInfo))
	for _, v := range user.UserInfo {
		l := mqservice.UserInfo{
			UserId:    v.UserId,
			UserName:  v.UserName,
			UserType:  v.UserType,
			UserGrade: v.UserGrade,
		}
		list = append(list, &l)
	}
	// rpc
	rsp, err := l.svcCtx.AssessMQClient.ImportUserInfo(l.ctx, &mqservice.ImportUserReq{
		Infos: list,
	})
	if err != nil {
		l.Logger.Error("rpc call ImportUserInfo error:", err)
		return &types.ImportUserRsp{
			Code:   250,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}

	return &types.ImportUserRsp{
		Code:   200,
		Msg:    "success",
		Result: rsp.GetResult(),
	}, nil
}
