package logic

import (
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"context"
	"github.com/spf13/cast"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMultiAlgoAssessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMultiAlgoAssessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMultiAlgoAssessLogic {
	return &GetMultiAlgoAssessLogic{
		ctx: ctx,

		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMultiAlgoAssess 取多个算法绩效时间线  (dashboard 时间线）
func (l *GetMultiAlgoAssessLogic) GetMultiAlgoAssess(in *proto.MultiAlgoReq) (*proto.MultiAlgoRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("in GetMultiAlgoAssess get req:", in)
	var result []*models.TbAlgoSummary
	var err error
	if in.GetSceneType() == 1 { // dashboard top4算法绩效
		result, err = l.svcCtx.SummaryRepo.GetTopRankAlgoSummary(l.ctx, in.GetDate(), in.GetUserId(), in.GetUserType(), int(in.GetAlgoType()))
		if err != nil {
			l.Logger.Error("query top rank algo error:", err)
			return &proto.MultiAlgoRsp{}, nil
		}
	}
	out, err := l.GetTimeLineByAlgoId(in.GetDate(), in.GetUserType(), result)
	if err != nil {
		l.Logger.Error("query time line error:", err)

	}
	return &proto.MultiAlgoRsp{
		Code:    0,
		Msg:     "success",
		Summary: out,
	}, nil
}

func (l *GetMultiAlgoAssessLogic) GetTimeLineByAlgoId(date int64, userType int32, result []*models.TbAlgoSummary) ([]*proto.AssessSummary, error) {
	var summary []*proto.AssessSummary
	for k, v := range result {
		tl, err := l.svcCtx.TimeLineRepo.GetAlgoTimeLineByAllUser(l.ctx, date, v.AlgoId, v.UserId, userType)
		if err != nil {
			l.Logger.Error("query time line error:", err)
			return nil, err
		}
		var a []*proto.TimeLine
		for _, tv := range tl {
			tp := tools.GetTimePoint(cast.ToString(tv.TransactTime))
			t := proto.TimeLine{
				TimePoint:   tp,
				AssessScore: int32(tv.AssessScore),
				Progress:    0,
			}
			a = append(a, &t)
		}
		// 时间线补一个初始起点 9：30分的数据，如果没有交易数据，则置为0
		if len(a) > 0 {
			if a[0].TimePoint != "09:30" {
				t := &proto.TimeLine{
					TimePoint:   "09:30",
					AssessScore: 0,
					Progress:    0,
					RiskScore:   0,
				}
				a = append(a, &proto.TimeLine{})
				copy(a[1:], a[0:]) // 所有元素后移
				a[0] = t
			}
		}

		s := proto.AssessSummary{
			AlgoName:      v.AlgoName,
			EconomyScore:  int32(v.EconomyScore),
			ProgressScore: int32(v.ProgressScore),
			RiskScore:     int32(v.RiskScore),
			AssessScore:   int32(v.AssessScore),
			StableScore:   int32(v.StableScore),
			TotalScore:    int32(v.CumsumScore),
			Ranking:       int32(k + 1),
			Desc:          "先乱写",
			Tl:            a,
		}
		summary = append(summary, &s)
	}
	return summary, nil
}
