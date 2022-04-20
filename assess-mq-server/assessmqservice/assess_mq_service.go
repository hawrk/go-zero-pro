// Code generated by goctl. DO NOT EDIT!
// Source: mqassess.proto

package assessmqservice

import (
	"context"

	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AssessInfo    = proto.AssessInfo
	GeneralReq    = proto.GeneralReq
	GeneralRsp    = proto.GeneralRsp
	MarketDataReq = proto.MarketDataReq
	MarketDataRsp = proto.MarketDataRsp

	AssessMqService interface {
		//  获取绩效概况
		GetMqGeneral(ctx context.Context, in *GeneralReq, opts ...grpc.CallOption) (*GeneralRsp, error)
		//  推送行情数据
		PullMarketData(ctx context.Context, in *MarketDataReq, opts ...grpc.CallOption) (*MarketDataRsp, error)
	}

	defaultAssessMqService struct {
		cli zrpc.Client
	}
)

func NewAssessMqService(cli zrpc.Client) AssessMqService {
	return &defaultAssessMqService{
		cli: cli,
	}
}

//  获取绩效概况
func (m *defaultAssessMqService) GetMqGeneral(ctx context.Context, in *GeneralReq, opts ...grpc.CallOption) (*GeneralRsp, error) {
	client := proto.NewAssessMqServiceClient(m.cli.Conn())
	return client.GetMqGeneral(ctx, in, opts...)
}

//  推送行情数据
func (m *defaultAssessMqService) PullMarketData(ctx context.Context, in *MarketDataReq, opts ...grpc.CallOption) (*MarketDataRsp, error) {
	client := proto.NewAssessMqServiceClient(m.cli.Conn())
	return client.PullMarketData(ctx, in, opts...)
}