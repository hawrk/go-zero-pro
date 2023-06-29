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

type MultiAlgoAnalyseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMultiAlgoAnalyseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiAlgoAnalyseLogic {
	return &MultiAlgoAnalyseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// MultiAlgoAnalyse 对比分析
func (l *MultiAlgoAnalyseLogic) MultiAlgoAnalyse(req *types.MultiAnalyseReq) (resp *types.MultiAnalyseRsp, err error) {
	l.Logger.Infof("get req:%+v", req)
	crossDayFlag := false
	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))
	if start != end {
		crossDayFlag = true
	}
	var provider, algoTypeName, algoName string
	var algoNameList []string
	// 初始化页面，无算法名称
	if len(req.AlgoName) <= 0 {
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
			return &types.MultiAnalyseRsp{
				Code:     200,
				Msg:      "",
				CrossDay: crossDayFlag,
				List:     []types.AlgoScore{},
			}, nil
		}
		provider = dRsp.GetProvider()
		algoTypeName = dRsp.GetAlgoTypeName()
		algoName = dRsp.GetAlgoName()
		algoNameList = append(algoNameList, algoName)
	} else {
		for _, v := range req.AlgoName {
			algoNameList = append(algoNameList, v)
		}
	}
	if len(algoNameList) == 0 { // 无数据，直接返回
		l.Logger.Info("no data....")
		return &types.MultiAnalyseRsp{
			Code:     200,
			Msg:      "",
			CrossDay: crossDayFlag,
			List:     []types.AlgoScore{},
		}, nil
	}

	mReq := &assessservice.CompareReq{
		StartTime:    start,
		EndTime:      end,
		UserId:       req.UserId,
		UserType:     int32(req.UserType),
		AlgoName:     algoNameList,
		CrossDayFlag: crossDayFlag,
	}
	mRsp, err := l.svcCtx.AssessClient.CompareMultiAlgo(l.ctx, mReq)
	if err != nil {
		l.Logger.Error("call rpc CompareMultiAlgo error:", err)
		return &types.MultiAnalyseRsp{
			Code:     205,
			Msg:      err.Error(),
			CrossDay: crossDayFlag,
		}, nil
	}
	//l.Logger.Info("mRsp:", mRsp)

	var as []types.AlgoScore
	for _, v := range mRsp.GetAlgoScore() {
		// 拼五个维度的分数
		var di []types.AlgoDimension
		for _, d := range v.GetDimension() {
			i := types.AlgoDimension{
				ProfileType: d.GetProfileType(),
				Score:       int(d.GetProfileScore()),
				Desc:        d.GetProfileDesc(),
			}
			di = append(di, i)
		}
		// 拼时间线图
		var al []types.AnalyseTimeLine
		for _, t := range v.GetTl() {
			a := types.AnalyseTimeLine{
				TimePoint: t.GetTimePoint(),
				Score:     t.GetAssessScore(),
			}
			al = append(al, a)
		}

		a := types.AlgoScore{
			AlgoName:       v.GetAlgoName(),
			CompositeScore: v.GetTotalScore(),
			Ranking:        v.GetRanking(),
			Dimension:      di,
			Data:           al,
		}
		as = append(as, a)
	}

	rsp := &types.MultiAnalyseRsp{
		Code:         200,
		Msg:          "success",
		CrossDay:     crossDayFlag,
		Provider:     provider,
		AlgoTypeName: algoTypeName,
		AlgoName:     algoNameList,
		List:         as,
	}

	return rsp, nil
}
