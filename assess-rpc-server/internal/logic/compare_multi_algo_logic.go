package logic

import (
	"algo_assess/assess-rpc-server/internal/config"
	"algo_assess/models"
	"algo_assess/pkg/tools"
	"algo_assess/repo"
	"context"
	"errors"
	"github.com/spf13/cast"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompareMultiAlgoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompareMultiAlgoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompareMultiAlgoLogic {
	return &CompareMultiAlgoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CompareMultiAlgo 多算法对比
func (l *CompareMultiAlgoLogic) CompareMultiAlgo(in *proto.CompareReq) (*proto.CompareRsp, error) {
	l.Logger.Info("CompareMultiAlgo get req:", in)
	// algoName 先转换成algoID
	algoIds, err := l.svcCtx.AlgoInfoRepo.GetAlgoIdsByAlgoNames(l.ctx, in.GetAlgoName())
	if err != nil {
		l.Logger.Error("GetAlgoIdsByAlgoNames error:", err)
		return &proto.CompareRsp{
			Code: 205,
			Msg:  err.Error(),
		}, nil
	}
	if len(algoIds) <= 0 {
		l.Logger.Error("GetAlgoIdsByAlgoNames, no data")
		return &proto.CompareRsp{
			Code: 205,
			Msg:  errors.New("algo summary no data").Error(),
		}, nil
	}
	// 先根据用户ID到账户表反查用户名称和角色权限
	// 传入的userType 只区分普通用户(0) 和管理员(1),account表区分普通用户1，厂商2，管理员3，超管4
	account, err := l.svcCtx.UserInfoRepo.GetAccountInfoByUserId(l.ctx, in.GetUserId())
	if err != nil {
		l.Logger.Error("get account info error:", err)
	}

	var list []*proto.CompareAlgoScore
	if in.GetCrossDayFlag() { // 跨天请求
		// 查询评分页面的汇总数据
		summary, err := l.svcCtx.SummaryRepo.GetCrossDaySummaryByAlgoIds(l.ctx, algoIds,
			in.GetStartTime(), in.GetEndTime(), in.GetUserId(), in.GetUserType())
		if err != nil {
			l.Logger.Error("GetCrossDaySummaryByAlgoIds error:", err)
			return &proto.CompareRsp{
				Code: 205,
				Msg:  err.Error(),
			}, nil
		}
		l.Logger.Info("get result len:", len(summary))
		// 取综合评分列表
		scoreList, err := l.svcCtx.SummaryRepo.GetCumsumListCrossDay(l.ctx, in.GetStartTime(), in.GetEndTime(), account.UserType)
		if err != nil {
			l.Logger.Error("GetCumsumListCrossDay error:", err)
			return &proto.CompareRsp{
				Code: 205,
				Msg:  err.Error(),
			}, nil
		}
		//TODO: 这里更高一点的效率就是查TimeLine表时，根据date和algoId分组汇总查出来，再根据algoId 匹配，只需要一次查询
		for _, v := range summary {
			s := &proto.CompareAlgoScore{
				AlgoName:   v.AlgoName,
				TotalScore: int32(v.CumsumScore),
				Ranking:    tools.GetRanking(scoreList, int(v.CumsumScore)),
				Dimension:  l.GetMultiAlgoDimension(v),
				Tl:         l.GetCrossTimeLines(int(v.AlgoId), in.GetUserId(), in.GetUserType(), in.GetStartTime(), in.GetEndTime()),
			}
			list = append(list, s)
		}

	} else { // 当天
		summary, err := l.svcCtx.SummaryRepo.GetAlgoSummaryByAlgoIds(l.ctx, algoIds, in.GetStartTime(), in.GetUserId(), in.GetUserType())
		if err != nil {
			l.Logger.Error("GetAlgoSummaryByAlgoIds error:", err)
			return &proto.CompareRsp{
				Code: 205,
				Msg:  err.Error(),
			}, nil
		}
		role := account.UserType
		if role == 0 { // 默认能取到的 1-普通用户 2-算法厂商 3-管理员，如果都不是，那就给一个超管的角色
			role = 4
		}
		// 取综合评分列表
		scoreList, err := l.svcCtx.SummaryRepo.GetCumsumList(l.ctx, in.GetStartTime(), role)
		if err != nil {
			l.Logger.Error("GetCumsumList error:", err)
			return &proto.CompareRsp{
				Code: 205,
				Msg:  err.Error(),
			}, nil
		}
		//拿到评分数据后， 再根据算法ID到时间线表查
		tl := make(map[string][]*proto.TimeLine)
		l.GetAssessTimeLine(in.GetStartTime(), in.GetUserId(), in.GetUserType(), summary, tl)
		for _, v := range summary {
			s := &proto.CompareAlgoScore{
				AlgoName:   v.AlgoName,
				TotalScore: int32(v.CumsumScore),
				Ranking:    tools.GetRanking(scoreList, v.CumsumScore),
				Dimension:  l.GetAlgoDimension(v),
				Tl:         tl[v.AlgoName],
			}
			list = append(list, s)
		}
	}

	rsp := &proto.CompareRsp{
		Code:      200,
		Msg:       "success",
		AlgoScore: list,
	}
	return rsp, nil
}

func (l *CompareMultiAlgoLogic) GetAssessTimeLine(date int64, userId string, userType int32,
	sum []*models.TbAlgoSummary, tlMap map[string][]*proto.TimeLine) {
	for _, v := range sum {
		var sli []*proto.TimeLine
		out, err := l.svcCtx.TimeLineRepo.GetAlgoTimeLineByAllUser(l.ctx, date, v.AlgoId, userId, userType)
		if err != nil {
			l.Logger.Error("GetAlgoTimeLine error:", err)
		}
		for _, o := range out {
			s := proto.TimeLine{
				TimePoint:   tools.GetTimePoint(cast.ToString(o.TransactTime)),
				AssessScore: tools.ScoreRound(o.AssessScore),
			}
			sli = append(sli, &s)
		}
		// 时间线补一个初始起点 9：30分的数据，如果没有交易数据，则置为0
		if len(sli) > 0 {
			if sli[0].TimePoint != "09:30" {
				t := &proto.TimeLine{
					TimePoint:   "09:30",
					AssessScore: 0,
					Progress:    0,
					RiskScore:   0,
				}
				sli = append(sli, &proto.TimeLine{})
				copy(sli[1:], sli[0:]) // 所有元素后移
				sli[0] = t
			}
		}
		tlMap[v.AlgoName] = sli
	}
}

func (l *CompareMultiAlgoLogic) GetCrossTimeLines(algoId int, userId string, userType int32, start, end int64) []*proto.TimeLine {
	if algoId == 0 {
		return nil
	}
	var tls []*proto.TimeLine
	out, err := l.svcCtx.TimeLineRepo.GetMultiTimeLine(l.ctx, start, end, userId, userType, algoId)
	if err != nil {
		l.Logger.Error("in GetCrossTimeLines , GetMultiTimeLine error:", err)
		return nil
	}
	for _, v := range out {
		tl := &proto.TimeLine{
			TimePoint:   tools.GetTimePointByDay(cast.ToString(v.Date)),
			AssessScore: tools.ScoreRound(v.AssessScore),
			Progress:    v.Progress,
			RiskScore:   tools.ScoreRound(v.RiskScore),
		}
		tls = append(tls, tl)
	}
	// 时间线补一个初始起点,如果跨天，则需要判断拿到起始的日期，如果没有交易数据，则置为0
	if len(tls) > 0 {
		if tls[0].TimePoint != tools.GetTimePointByDay(cast.ToString(start)) {
			t := &proto.TimeLine{
				TimePoint:   tools.GetTimePointByDay(cast.ToString(start)), // 设置起点坐标
				AssessScore: 0,
				Progress:    0,
				RiskScore:   0,
			}
			tls = append(tls, &proto.TimeLine{})
			copy(tls[1:], tls[0:]) // 所有元素后移
			tls[0] = t
		}
	}
	return tls
}

func (l *CompareMultiAlgoLogic) GetAlgoDimension(v *models.TbAlgoSummary) []*proto.AlgoDimension {
	var rada []*proto.AlgoDimension
	r1 := &proto.AlgoDimension{ // 经济性
		ProfileType:  1,
		ProfileScore: int32(v.EconomyScore),
		ProfileDesc:  config.GetEconomyDesc(v.EconomyScore),
	}
	r2 := &proto.AlgoDimension{ // 完成度
		ProfileType:  2,
		ProfileScore: int32(v.ProgressScore),
		ProfileDesc:  config.GetProgressDesc(v.ProgressScore),
	}
	r3 := &proto.AlgoDimension{ // 风险度
		ProfileType:  3,
		ProfileScore: int32(v.RiskScore),
		ProfileDesc:  config.GetRiskDesc(v.RiskScore),
	}
	r4 := &proto.AlgoDimension{ // 绩效
		ProfileType:  4,
		ProfileScore: int32(v.AssessScore),
		ProfileDesc:  config.GetAssessDesc(v.AssessScore),
	}
	r5 := &proto.AlgoDimension{ // 稳定性
		ProfileType:  5,
		ProfileScore: int32(v.StableScore),
		ProfileDesc:  config.GetStabilityDesc(v.StableScore),
	}
	rada = append(rada, r1, r2, r3, r4, r5)
	return rada
}

func (l *CompareMultiAlgoLogic) GetMultiAlgoDimension(v *repo.AvgSummary) []*proto.AlgoDimension {
	var rada []*proto.AlgoDimension
	r1 := &proto.AlgoDimension{ // 经济性
		ProfileType:  1,
		ProfileScore: tools.ScoreRound(v.EconomyScore),
		ProfileDesc:  config.GetEconomyDesc(int(tools.ScoreRound(v.EconomyScore))),
	}
	r2 := &proto.AlgoDimension{ // 完成度
		ProfileType:  2,
		ProfileScore: tools.ScoreRound(v.ProgressScore),
		ProfileDesc:  config.GetProgressDesc(int(tools.ScoreRound(v.ProgressScore))),
	}
	r3 := &proto.AlgoDimension{ // 风险度
		ProfileType:  3,
		ProfileScore: tools.ScoreRound(v.RiskScore),
		ProfileDesc:  config.GetRiskDesc(int(tools.ScoreRound(v.RiskScore))),
	}
	r4 := &proto.AlgoDimension{ // 绩效
		ProfileType:  4,
		ProfileScore: tools.ScoreRound(v.AssessScore),
		ProfileDesc:  config.GetAssessDesc(int(tools.ScoreRound(v.AssessScore))),
	}
	r5 := &proto.AlgoDimension{ // 稳定性
		ProfileType:  5,
		ProfileScore: tools.ScoreRound(v.StableScore),
		ProfileDesc:  config.GetStabilityDesc(int(tools.ScoreRound(v.StableScore))),
	}
	rada = append(rada, r1, r2, r3, r4, r5)
	return rada
}
