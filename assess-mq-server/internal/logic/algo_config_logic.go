package logic

import (
	"algo_assess/assess-mq-server/internal/dao"
	"context"

	"algo_assess/assess-mq-server/internal/svc"
	"algo_assess/assess-mq-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlgoConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAlgoConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoConfigLogic {
	return &AlgoConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AlgoConfig 配置： 算法配置
func (l *AlgoConfigLogic) AlgoConfig(in *proto.AlgoConfigReq) (*proto.AlgoConfigRsp, error) {
	l.Logger.Infof("AlgoConfig, req:%+v", in)
	id, _, err := l.svcCtx.BusiConfigRepo.GetBusiConfigByProfile(l.ctx, in.GetProfileType())
	if err != nil {
		l.Logger.Error("GetBusiConfigByProfile error:", err)
		return &proto.AlgoConfigRsp{
			Code:   200,
			Msg:    err.Error(),
			Result: 2,
		}, nil
	}
	if id <= 0 { // 无记录，新插入一条
		if err := l.svcCtx.BusiConfigRepo.CreateBusiConfig(l.ctx, in.GetProfileType(), in.GetAlgoConfig()); err != nil {
			l.Logger.Error("CreateBusiConfig error:", err)
			return &proto.AlgoConfigRsp{
				Code:   200,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}
	} else { // 有记录，直接更新
		if err := l.svcCtx.BusiConfigRepo.UpdateBusiParam(l.ctx, id, in.GetAlgoConfig()); err != nil {
			l.Logger.Error("UpdateBusiParam error:", err)
			return &proto.AlgoConfigRsp{
				Code:   200,
				Msg:    err.Error(),
				Result: 2,
			}, nil
		}
	}

	// 同步更新本地Cache
	l.SyncGWeight()

	return &proto.AlgoConfigRsp{
		Code:   200,
		Msg:    "success",
		Result: 1,
	}, nil
}

func (l *AlgoConfigLogic) SyncGWeight() {
	infos, err := l.svcCtx.BusiConfigRepo.GetAllBusiConfig(context.Background())
	if err != nil {
		l.Logger.Error("SyncGWeight get config error:", err)
		return
	}
	// 初始化 画像基础数据
	m := make(map[int]string)
	for _, v := range infos {
		if v.BusiType == 4 { // 取画像的权重配置
			m[v.SecType] = v.Params
		}
	}
	dao.GScoreConf = dao.NewWeights(m)
}
