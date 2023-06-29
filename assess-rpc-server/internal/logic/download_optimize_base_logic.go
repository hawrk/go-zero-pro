package logic

import (
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadOptimizeBaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDownloadOptimizeBaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadOptimizeBaseLogic {
	return &DownloadOptimizeBaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DownloadOptimizeBase 下载一键优选基础数据
func (l *DownloadOptimizeBaseLogic) DownloadOptimizeBase(in *proto.DownloadOptimizeBaseReq) (*proto.DownloadOptimizeBaseRsp, error) {
	count, optimizeBases, err := l.svcCtx.OptimizeBaseRepo.SelectAllOptimizeBase()
	if err != nil {
		return &proto.DownloadOptimizeBaseRsp{
			Code:  1000,
			Msg:   "导出失败",
			Total: count,
			List:  nil,
		}, err
	}
	return &proto.DownloadOptimizeBaseRsp{
		Code:  0,
		Msg:   "导出成功",
		Total: count,
		List:  ToProtoOptimizeBase(optimizeBases),
	}, nil
}
