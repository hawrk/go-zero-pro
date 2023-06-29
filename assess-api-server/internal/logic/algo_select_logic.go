package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AlgoSelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoSelectLogic {
	return &AlgoSelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlgoSelectLogic) AlgoSelect(req *types.AlgoSelectReq) (resp *types.AlgoSelectRsp, err error) {
	// todo: add your logic here and delete this line
	d := cast.ToInt64(time.Unix(req.Date, 0).Format(global.TimeFormatDay))
	l.Logger.Infof("get req:%+v", req)
	q := &assessservice.ChooseAlgoReq{
		ChooseType:   req.ChooseType,
		Provider:     req.Provider,
		AlgoTypeName: req.AlgoType,
		AlgoName:     req.AlgoName,
		Date:         d,
	}
	// 下拉选择列表需要根据用户拥有的算法权限拉取，需求请总线的rpc服务，不再走绩效的rpc服务
	/*  ---先不切
	if (req.ChooseType == 1 || req.ChooseType == 2 || req.ChooseType == 3 || req.ChooseType == 9) && req.UserId != "" {
		l.Logger.Info("into mornano rpc call.....")
		brsp, err := l.svcCtx.MornanoClient.GetAlgoChooseList(l.ctx, &mornanoservice.AlgoChooseReq{
			UserId:       req.UserId,
			SelectType:   req.ChooseType,
			Provider:     req.Provider,
			AlgoTypeName: req.AlgoType,
			AlgoName:     req.AlgoName,
		})
		if err != nil {
			l.Logger.Error("call mornano rpc AlgoChooseList error:", err)
			return &types.AlgoSelectRsp{
				Code: 204,
				Msg:  err.Error(),
			}, nil
		}
		return &types.AlgoSelectRsp{
			Code:     200,
			Msg:      "success",
			Provider: brsp.GetProvider(),
			AlgoType: brsp.GetAlgoTypeName(),
			AlgoName: brsp.GetAlgoName(),
		}, nil
	}
	*/
	// 走原来的绩效rpc服务
	rsp, err := l.svcCtx.AssessClient.ChooseAlgoInfo(l.ctx, q)
	if err != nil {
		l.Logger.Error("call asses rcp choose algo info error:", err)
		return &types.AlgoSelectRsp{
			Code: 204,
			Msg:  err.Error(),
		}, nil
	}
	r := &types.AlgoSelectRsp{
		Code:     200,
		Msg:      "success",
		Provider: rsp.GetProvider(),
		AlgoType: rsp.GetAlgoTypeName(),
		AlgoName: rsp.GetAlgoName(),
	}
	return r, nil
}
