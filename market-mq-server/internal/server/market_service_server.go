// Code generated by goctl. DO NOT EDIT!
// Source: market.proto

package server

import (
	"context"

	"algo_assess/market-mq-server/internal/logic"
	"algo_assess/market-mq-server/internal/svc"
	"algo_assess/market-mq-server/proto"
)

type MarketServiceServer struct {
	svcCtx *svc.ServiceContext
	proto.UnimplementedMarketServiceServer
}

func NewMarketServiceServer(svcCtx *svc.ServiceContext) *MarketServiceServer {
	return &MarketServiceServer{
		svcCtx: svcCtx,
	}
}

//  获取行情价格数量信息
func (s *MarketServiceServer) GetMarketInfo(ctx context.Context, in *proto.MarkReq) (*proto.MarkRsp, error) {
	l := logic.NewGetMarketInfoLogic(ctx, s.svcCtx)
	return l.GetMarketInfo(in)
}
