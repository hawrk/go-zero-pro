// Code generated by goctl. DO NOT EDIT!
// Source: mornano.proto

package mornanoservice

import (
	"context"

	"algo_assess/mornano-rpc-server/proto"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AlgoChooseReq = proto.AlgoChooseReq
	AlgoChooseRsp = proto.AlgoChooseRsp
	AlgoInfo      = proto.AlgoInfo
	AlgoInfoReq   = proto.AlgoInfoReq
	AlgoInfoRsp   = proto.AlgoInfoRsp
	CapitailRsp   = proto.CapitailRsp
	CapitalReq    = proto.CapitalReq
	LoginReq      = proto.LoginReq
	LoginRsp      = proto.LoginRsp
	StockPosition = proto.StockPosition

	MornanoService interface {
		//  算法选择框数据
		GetAlgoChooseList(ctx context.Context, in *AlgoChooseReq, opts ...grpc.CallOption) (*AlgoChooseRsp, error)
		//  登陆校验   --原则上只返回密码和角色
		LoginCheck(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRsp, error)
		//  查询算法基础信息
		GetAlgoInfo(ctx context.Context, in *AlgoInfoReq, opts ...grpc.CallOption) (*AlgoInfoRsp, error)
		//  查询用户资金和证券持仓信息  (用户画像)
		GetUserCapital(ctx context.Context, in *CapitalReq, opts ...grpc.CallOption) (*CapitailRsp, error)
	}

	defaultMornanoService struct {
		cli zrpc.Client
	}
)

func NewMornanoService(cli zrpc.Client) MornanoService {
	return &defaultMornanoService{
		cli: cli,
	}
}

//  算法选择框数据
func (m *defaultMornanoService) GetAlgoChooseList(ctx context.Context, in *AlgoChooseReq, opts ...grpc.CallOption) (*AlgoChooseRsp, error) {
	client := proto.NewMornanoServiceClient(m.cli.Conn())
	return client.GetAlgoChooseList(ctx, in, opts...)
}

//  登陆校验   --原则上只返回密码和角色
func (m *defaultMornanoService) LoginCheck(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginRsp, error) {
	client := proto.NewMornanoServiceClient(m.cli.Conn())
	return client.LoginCheck(ctx, in, opts...)
}

//  查询算法基础信息
func (m *defaultMornanoService) GetAlgoInfo(ctx context.Context, in *AlgoInfoReq, opts ...grpc.CallOption) (*AlgoInfoRsp, error) {
	client := proto.NewMornanoServiceClient(m.cli.Conn())
	return client.GetAlgoInfo(ctx, in, opts...)
}

//  查询用户资金和证券持仓信息  (用户画像)
func (m *defaultMornanoService) GetUserCapital(ctx context.Context, in *CapitalReq, opts ...grpc.CallOption) (*CapitailRsp, error) {
	client := proto.NewMornanoServiceClient(m.cli.Conn())
	return client.GetUserCapital(ctx, in, opts...)
}