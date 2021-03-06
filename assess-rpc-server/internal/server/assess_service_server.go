// Code generated by goctl. DO NOT EDIT!
// Source: assess.proto

package server

import (
	"context"

	"algo_assess/assess-rpc-server/internal/logic"
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
)

type AssessServiceServer struct {
	svcCtx *svc.ServiceContext
	proto.UnimplementedAssessServiceServer
}

func NewAssessServiceServer(svcCtx *svc.ServiceContext) *AssessServiceServer {
	return &AssessServiceServer{
		svcCtx: svcCtx,
	}
}

//  获取绩效概况
func (s *AssessServiceServer) GetGeneral(ctx context.Context, in *proto.GeneralReq) (*proto.GeneralRsp, error) {
	l := logic.NewGetGeneralLogic(ctx, s.svcCtx)
	return l.GetGeneral(in)
}
