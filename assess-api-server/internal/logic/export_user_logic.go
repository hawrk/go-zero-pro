package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"
	"encoding/xml"
	"net/http"

	"algo_assess/assess-api-server/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ExportUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExportUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportUserLogic {
	return &ExportUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExportUserLogic) ExportUser(w http.ResponseWriter) error {
	// todo: add your logic here and delete this line
	l.Logger.Info("into ExportUser......")
	rsp, err := l.svcCtx.AssessMQClient.ExportUserInfo(l.ctx, &mqservice.ExportUserReq{})
	if err != nil {
		l.Logger.Error("rpc call ExportUserInfo error:", err)
		return err
	}
	list := make([]*UserData, 0, len(rsp.GetInfos()))
	for _, v := range rsp.GetInfos() {
		s := UserData{
			UserId:           v.GetUserId(),
			UserName:         v.GetUserName(),
			UserPasswd:       "",
			PasswdEnCodeType: 0,
			UserType:         v.GetUserType(),
			RiskGroup:        0,
			UuserId:          "",
			UserGrade:        v.GetUserGrade(),
		}
		list = append(list, &s)
	}
	header := w.Header()
	header.Add("Content-Type", "application/octet-stream")
	header.Add("Content-Disposition", "filename="+"userInfo.xml")
	output, _ := xml.MarshalIndent(list, " ", " ")
	_, _ = w.Write([]byte(xml.Header))
	_, _ = w.Write(output)
	return nil
}
