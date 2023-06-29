package logic

import (
	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"
	"algo_assess/global"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"
	"github.com/spf13/cast"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAlgoTimeLineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAlgoTimeLineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAlgoTimeLineLogic {
	return &GetAlgoTimeLineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAlgoTimeLine 时间线图
func (l *GetAlgoTimeLineLogic) GetAlgoTimeLine(in *proto.TimeLineReq) (*proto.TimeLineRsp, error) {
	l.Logger.Info("in GetAlgoTimeLine, get req:", in)
	var out []*proto.TimeLine
	if in.GetCrossDayFlag() { // 跨天
		var result []*models.TbAlgoTimeLine
		var err error
		if in.GetSourceFrom() == global.SourceFromOrigin || in.GetSourceFrom() == global.SourceFromImport {
			result, err = l.svcCtx.TimeLineOrigRepo.GetMultiTimeLineBatch(l.ctx, in.GetBatchNo(), in.GetUserId(), in.GetUserType(), in.GetAlgoId())
			if err != nil {
				l.Logger.Error("in GetMultiTimeLineBatch error:", err)
				return &proto.TimeLineRsp{}, nil
			}
		} else {
			result, err = l.svcCtx.TimeLineRepo.GetMultiTimeLine(l.ctx, in.GetStartTime(), in.GetEndTime(),
				in.GetUserId(), in.GetUserType(), int(in.GetAlgoId()))
			if err != nil {
				l.Logger.Error("in GetAlgoTimeLine,GetMultiTimeLine error:", err)
				return &proto.TimeLineRsp{}, nil
			}
		}
		for _, v := range result {
			tl := &proto.TimeLine{
				TimePoint:   tools.GetTimePointByDay(cast.ToString(v.Date)),
				AssessScore: int32(v.AssessScore),
				Progress:    v.Progress,
				RiskScore:   int32(v.RiskScore),
			}
			out = append(out, tl)
		}
		// 时间线补一个初始起点,如果跨天，则需要判断拿到起始的日期，如果没有交易数据，则置为0
		if len(out) > 0 {
			if out[0].TimePoint != tools.GetTimePointByDay(cast.ToString(in.GetStartTime())) {
				t := &proto.TimeLine{
					TimePoint:   tools.GetTimePointByDay(cast.ToString(in.GetStartTime())), // 设置起点坐标
					AssessScore: 0,
					Progress:    0,
					RiskScore:   0,
				}
				out = append(out, &proto.TimeLine{})
				copy(out[1:], out[0:]) // 所有元素后移
				out[0] = t
			}
		}
	} else { // 当天
		var result []*models.TbAlgoTimeLine
		var err error
		if in.SourceFrom == global.SourceFromImport || in.SourceFrom == global.SourceFromOrigin { // 订单导入、原始订单
			result, err = l.svcCtx.TimeLineOrigRepo.GetImportAlgoTimeLine(l.ctx, in.GetStartTime(), in.GetBatchNo(), in.GetUserId(), in.GetUserType())
			if err != nil {
				l.Logger.Error("GetImportAlgoTimeLine error:", err)
				return &proto.TimeLineRsp{}, nil
			}
		} else {
			result, err = l.svcCtx.TimeLineRepo.GetAlgoTimeLineByAllUser(l.ctx, in.GetStartTime(), int(in.GetAlgoId()), in.GetUserId(), in.GetUserType())
			if err != nil {
				l.Logger.Error("get time line data error:", err)
				return &proto.TimeLineRsp{}, nil
			}
		}

		for _, v := range result {
			var i proto.TimeLine
			if in.GetLineType() == 1 { // 绩效
				i = proto.TimeLine{
					TimePoint:   tools.GetTimePoint(cast.ToString(v.TransactTime)),
					AssessScore: int32(v.AssessScore),
					Progress:    0,
				}
			} else if in.GetLineType() == 2 { // 完成度
				i = proto.TimeLine{
					TimePoint:   tools.GetTimePoint(cast.ToString(v.TransactTime)),
					AssessScore: 0,
					Progress:    v.Progress,
				}
			} else if in.GetLineType() == 12 { // 绩效+完成度
				i = proto.TimeLine{
					TimePoint:   tools.GetTimePoint(cast.ToString(v.TransactTime)),
					AssessScore: int32(v.AssessScore),
					Progress:    v.Progress,
				}
			}
			out = append(out, &i)
		}
		// 时间线补一个初始起点 9：30分的数据，如果没有交易数据，则置为0
		if len(out) > 0 {
			if out[0].TimePoint != "09:30" {
				t := &proto.TimeLine{
					TimePoint:   "09:30",
					AssessScore: 0,
					Progress:    0,
					RiskScore:   0,
				}
				out = append(out, &proto.TimeLine{})
				copy(out[1:], out[0:]) // 所有元素后移
				out[0] = t
			}
		}
	}

	l.Logger.Info("get rsp len:", len(out))
	rsp := &proto.TimeLineRsp{
		Line: out,
	}

	return rsp, nil
}
