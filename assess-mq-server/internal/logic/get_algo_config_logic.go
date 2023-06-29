package logic

import (
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoConfigLogic {
	return &GetAlgoConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAlgoConfig 配置： 算法配置查询
func (l *GetAlgoConfigLogic) GetAlgoConfig(in *proto.GetAlgoConfigReq) (*proto.GetAlgoConfigRsp, error) {
	l.Logger.Infof("GetAlgoConfig, req:%+v", in)
	_, param, err := l.svcCtx.BusiConfigRepo.GetBusiConfigByProfile(l.ctx, in.GetProfileType())
	if err != nil {
		l.Logger.Error("GetBusiConfigByProfile error:", err)
		return &proto.GetAlgoConfigRsp{
			Code:   200,
			Msg:    err.Error(),
			Config: "",
		}, nil
	}

	return &proto.GetAlgoConfigRsp{
		Code:   200,
		Msg:    "success",
		Config: param,
	}, nil
}
