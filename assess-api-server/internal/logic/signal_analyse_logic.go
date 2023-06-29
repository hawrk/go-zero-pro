package logic

import (
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"context"
	"github.com/spf13/cast"
	"k8s.io/apimachinery/pkg/util/rand"
	"time"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SignalAnalyseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSignalAnalyseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SignalAnalyseLogic {
	return &SignalAnalyseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SignalAnalyse 高阶分析 信号分析
func (l *SignalAnalyseLogic) SignalAnalyse(req *types.SignalReq) (resp *types.SignalRsp, err error) {
	l.Logger.Infof("in SignalAnalyse, get req:%+v", *req)
	if l.svcCtx.Config.WorkControl.EnableFakeMsg {
		return l.BuildVirtualRsp(req)
	}
	var algoId int32
	var algoTypeName, algoName, provider string

	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))
	if req.AlgoName == "" { // 算法名称为空的， 取兜底数据
		l.Logger.Info("algo_name empty, get default algo....")
		dReq := &assessservice.DefaultReq{
			Scene:     4,
			UserId:    req.UserId,
			UserType:  int32(req.UserType),
			StartTime: start,
			EndTime:   end,
		}
		dRsp, err := l.svcCtx.AssessClient.GetDefaultAlgo(l.ctx, dReq)
		if err != nil {
			l.Logger.Error("rpc call GetDefaultAlgo error:", err)
			return &types.SignalRsp{
				Code: 350,
				Msg:  err.Error(),
			}, nil
		}
		algoId = dRsp.GetAlgoId()
		algoTypeName = dRsp.GetAlgoTypeName()
		algoName = dRsp.GetAlgoName()
		provider = dRsp.GetProvider()
	} else { // 根据算法名称反查algoId
		alReq := &assessservice.ChooseAlgoReq{
			ChooseType: 4,
			AlgoName:   req.AlgoName,
		}

		alRsp, err := l.svcCtx.AssessClient.ChooseAlgoInfo(l.ctx, alReq)
		if err != nil {
			l.Logger.Error("call rpc ChooseAlgoInfo error:", err)
			return nil, err
		}
		algoId = alRsp.GetAlgoId()
	}
	if algoId <= 0 { // 无数据时直接返回
		l.Logger.Info("no data....")
		return &types.SignalRsp{
			Code:    200,
			Msg:     "no data",
			Signals: nil,
		}, nil
	}

	sReq := &proto.SignalReq{
		AlgoId:   algoId,
		StartDay: start,
		EndDay:   end,
		UserId:   req.UserId,
		UserType: int32(req.UserType),
	}
	sRsp, err := l.svcCtx.AssessClient.GetSignal(l.ctx, sReq)
	if err != nil {
		l.Logger.Error("rpc call GetSignal error:", err)
		return &types.SignalRsp{
			Code:    205,
			Msg:     err.Error(),
			Signals: nil,
		}, nil
	}
	// 拼返回报文
	var signal []types.SignalInfo
	for _, v := range sRsp.GetInfo() {
		s := types.SignalInfo{
			Day:      v.GetDay(),
			OrderCnt: int(v.GetOrderNum()),
			Progress: v.GetProgress(),
		}
		signal = append(signal, s)
	}
	return &types.SignalRsp{
		Code:         200,
		Msg:          "success",
		Provider:     provider,
		AlgoTypeName: algoTypeName,
		AlgoName:     algoName,
		Signals:      signal,
	}, nil
}

func (l *SignalAnalyseLogic) BuildVirtualRsp(req *types.SignalReq) (resp *types.SignalRsp, err error) {
	//TODO: 拼返回数据
	var ss []types.SignalInfo
	for i := -14; i <= -1; i++ {
		d := time.Now().AddDate(0, 0, i).Format(global.TimeFormatDaySp)
		s := types.SignalInfo{
			Day:      d,
			OrderCnt: rand.IntnRange(5, 150),
			Progress: RandFloat64(10.0, 100.0),
		}
		ss = append(ss, s)
	}
	return &types.SignalRsp{
		Code:    200,
		Msg:     "success",
		Signals: ss,
	}, nil
}
