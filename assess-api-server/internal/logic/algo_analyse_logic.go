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

type AlgoAnalyseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoAnalyseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoAnalyseLogic {
	return &AlgoAnalyseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AlgoAnalyse 单个算法 多日分析
func (l *AlgoAnalyseLogic) AlgoAnalyse(req *types.AnalyseReq) (resp *types.AnalyseRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("get req :%+v", req)

	crossDayFlag := false
	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))
	if start != end { // 判断是否是跨天
		crossDayFlag = true
	}
	var algoId int32
	var provider, algoTypeName, algoName string
	// 如果algo_name 没传，则先返回空数据
	if req.AlgoName == "" {
		l.Logger.Info("algo_name empty, get default algo....")
		dReq := &assessservice.DefaultReq{
			Scene:     1,
			UserId:    req.UserId,
			UserType:  int32(req.UserType),
			StartTime: start,
			EndTime:   end,
		}
		dRsp, err := l.svcCtx.AssessClient.GetDefaultAlgo(l.ctx, dReq)
		if err != nil {
			l.Logger.Error("rpc call GetDefaultAlgo error:", err)
			return &types.AnalyseRsp{
				Code:     200,
				Msg:      "success",
				CrossDay: crossDayFlag,
				Data:     []types.AnalyseLine{},
			}, nil
		}
		algoId = dRsp.GetAlgoId()
		provider = dRsp.GetProvider()
		algoTypeName = dRsp.GetAlgoTypeName()
		algoName = dRsp.GetAlgoName()
	} else {
		// 先反查一下算法ID
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
	if algoId == 0 {
		l.Logger.Info("no data....")
		return &types.AnalyseRsp{
			Code:         200,
			Msg:          "success",
			CrossDay:     crossDayFlag,
			Provider:     provider,
			AlgoTypeName: algoTypeName,
			AlgoName:     "",
			Data:         nil,
		}, nil
	}

	mReq := &assessservice.MultiDayReq{
		StartTime:    start,
		EndTime:      end,
		CrossDayFlag: crossDayFlag,
		UserId:       req.UserId,
		UserType:     int32(req.UserType),
		AlgoId:       algoId,
	}

	mRsp, err := l.svcCtx.AssessClient.MultiDayAnalyse(l.ctx, mReq)
	if err != nil {
		l.Logger.Error("call rpc MultiDayAnalyse error:", err)
		return &types.AnalyseRsp{
			Code:     208,
			Msg:      err.Error(),
			CrossDay: crossDayFlag,
			Data:     nil,
		}, nil
	}

	p := l.BuildMultiDayRsp(mRsp, crossDayFlag, provider, algoTypeName, algoName)

	return p, nil
}

func (l *AlgoAnalyseLogic) BuildMultiDayRsp(r *assessservice.MultiDayRsp, crossFlag bool, provider, algoTypeName, AlgoName string) *types.AnalyseRsp {
	var list []types.AnalyseLine
	// 拼绩效
	var aslist []types.AnalyseTimeLine
	for _, v := range r.GetTl() {
		a := types.AnalyseTimeLine{
			TimePoint: v.GetTimePoint(),
			Score:     v.GetAssessScore(),
		}
		aslist = append(aslist, a)
	}
	al := types.AnalyseLine{
		ProfileType: 4,
		Point:       aslist,
	}
	// 拼完成度
	var pglist []types.AnalyseTimeLine
	for _, v := range r.GetTl() {
		p := types.AnalyseTimeLine{
			TimePoint: v.GetTimePoint(),
			Score:     int32(v.GetProgress()),
		}
		pglist = append(pglist, p)
	}
	pl := types.AnalyseLine{
		ProfileType: 2,
		Point:       pglist,
	}
	// 拼风险度
	var rklist []types.AnalyseTimeLine
	for _, v := range r.GetTl() {
		r := types.AnalyseTimeLine{
			TimePoint: v.GetTimePoint(),
			Score:     v.GetRiskScore(),
		}
		rklist = append(rklist, r)
	}
	rl := types.AnalyseLine{
		ProfileType: 3,
		Point:       rklist,
	}
	// 拼齐三个维度的时间线
	list = append(list, al, pl, rl)

	rsp := &types.AnalyseRsp{
		Code:         200,
		Msg:          "success",
		CrossDay:     crossFlag,
		Provider:     provider,
		AlgoTypeName: algoTypeName,
		AlgoName:     AlgoName,
		Data:         list,
	}
	return rsp
}
