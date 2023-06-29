package logic

import (
	"algo_assess/pkg/tools"
	"context"
	"github.com/spf13/cast"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultiDayAnalyseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMultiDayAnalyseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MultiDayAnalyseLogic {
	return &MultiDayAnalyseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MultiDayAnalyse 多日分析
func (l *MultiDayAnalyseLogic) MultiDayAnalyse(in *proto.MultiDayReq) (*proto.MultiDayRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("MultiDayAnalyse req:", in)
	var tls []*proto.TimeLine
	// 区分是当天的还是跨天的
	if in.GetCrossDayFlag() { // 跨天，需要拼多天的时间线
		out, err := l.svcCtx.TimeLineRepo.GetMultiTimeLine(l.ctx, in.GetStartTime(), in.GetEndTime(),
			in.GetUserId(), in.GetUserType(), int(in.GetAlgoId()))
		if err != nil {
			l.Logger.Error("cross day GetMultiTimeLine error:", err)
			return &proto.MultiDayRsp{
				Code: 206,
				Msg:  err.Error(),
				Tl:   nil,
			}, nil
		}
		for _, v := range out {
			tl := &proto.TimeLine{
				TimePoint:   tools.GetTimePointByDay(cast.ToString(v.Date)),
				AssessScore: int32(v.AssessScore),
				Progress:    v.Progress,
				RiskScore:   int32(v.RiskScore),
			}
			tls = append(tls, tl)
		}
		// 时间线补一个初始起点,如果跨天，则需要判断拿到起始的日期，如果没有交易数据，则置为0
		if len(tls) > 0 {
			if tls[0].TimePoint != tools.GetTimePointByDay(cast.ToString(in.GetStartTime())) {
				t := &proto.TimeLine{
					TimePoint:   tools.GetTimePointByDay(cast.ToString(in.GetStartTime())), // 设置起点坐标
					AssessScore: 0,
					Progress:    0,
					RiskScore:   0,
				}
				tls = append(tls, &proto.TimeLine{})
				copy(tls[1:], tls[0:]) // 所有元素后移
				tls[0] = t
			}
		}
	} else { // 只查当天
		out, err := l.svcCtx.TimeLineRepo.GetAlgoTimeLineByAllUser(l.ctx, in.GetStartTime(),
			int(in.GetAlgoId()), in.GetUserId(), in.GetUserType())
		if err != nil {
			l.Logger.Error("cur day GetMultiTimeLine error:", err)
			return &proto.MultiDayRsp{
				Code: 206,
				Msg:  err.Error(),
				Tl:   nil,
			}, nil
		}
		for _, v := range out {
			tl := &proto.TimeLine{
				TimePoint:   tools.GetTimePoint(cast.ToString(v.TransactTime)),
				AssessScore: int32(v.AssessScore),
				Progress:    v.Progress,
				RiskScore:   int32(v.RiskScore),
			}
			tls = append(tls, tl)
		}
		// 时间线补一个初始起点 9：30分的数据，如果没有交易数据，则置为0
		if len(tls) > 0 {
			if tls[0].TimePoint != "09:30" {
				t := &proto.TimeLine{
					TimePoint:   "09:30",
					AssessScore: 0,
					Progress:    0,
					RiskScore:   0,
				}
				tls = append(tls, &proto.TimeLine{})
				copy(tls[1:], tls[0:]) // 所有元素后移
				tls[0] = t
			}
		}
	}
	resp := &proto.MultiDayRsp{
		Code: 200,
		Msg:  "success",
		Tl:   tls,
	}

	return resp, nil
}
