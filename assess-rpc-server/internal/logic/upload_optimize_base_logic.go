package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadOptimizeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadOptimizeBaseLogic {
	return &UploadOptimizeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 导入一键优选基础数据
func (l *UploadOptimizeBaseLogic) UploadOptimizeBase(in *proto.UploadOptimizeBaseReq) (*proto.UploadOptimizeBaseRsp, error) {
	err := l.svcCtx.OptimizeBaseRepo.BatchUpload(in)
	if err != nil {
		return &proto.UploadOptimizeBaseRsp{
			Code: 10000,
			Msg:  "上传失败",
		}, err
	}
	return &proto.UploadOptimizeBaseRsp{
		Code: 0,
		Msg:  "上传成功",
	}, nil
}
