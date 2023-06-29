package logic

import (
	"algo_assess/global"
	"algo_assess/pkg/tools"
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserProfileLogic {
	return &GetUserProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserProfile 高阶分析：取用户画像绩效信息
func (l *GetUserProfileLogic) GetUserProfile(in *proto.UserProfileReq) (*proto.UserProfileRsp, error) {
	l.Logger.Infof("in GetUserProfile, get req:%+v", in)
	// 1. 取算法画像
	summary, err := l.svcCtx.SummaryRepo.GetAlgoSummary(l.ctx, in.GetUserId(), in.GetUserType(), in.GetCurDay(), in.GetAlgoId())
	if err != nil {
		l.Logger.Error("GetAlgoSummary error:", err)
	}
	// 2. 取盈亏信息
	// 盈亏金额， 交易次数，交易总金额，完成度
	pa, tn, tc, pr := l.GetProfile(in)

	//3. 取用户信息
	account, _ := l.svcCtx.UserInfoRepo.GetAccountInfoByUserId(l.ctx, in.GetUserId())

	grade := l.GetUserProperty(in.GetUserType(), account.UserGrade)

	// 4. 查排名
	// 取综合评分列表
	role := account.UserType
	if role == 0 { // 默认能取到的 1-普通用户 2-算法厂商 3-管理员，如果都不是，那就给一个超管的角色
		role = 4
	}
	scoreList, _ := l.svcCtx.SummaryRepo.GetCumsumList(l.ctx, in.GetCurDay(), role)

	var rsp proto.UserProfileRsp
	var di []*proto.AlgoDimension
	var score int32
	for _, v := range summary {
		di = BuildAlgoDimension(v.EconomyScore, v.ProgressScore, v.RiskScore, v.AssessScore, v.StableScore)
		score = int32(v.CumsumScore)
		break
	}
	rsp = proto.UserProfileRsp{
		Code:        200,
		Msg:         "success",
		Profit:      pa,
		TradeCnt:    int32(tn),
		TradeAmount: tc,
		UserGrade:   grade,
		Progress:    pr,
		Dimension:   di,
		TotalScore:  score,
		Ranking:     tools.GetRanking(scoreList, int(score)),
	}

	return &rsp, nil
}

// GetUserProperty 取用户级别
func (l *GetUserProfileLogic) GetUserProperty(userType int32, grade string) string {
	var userGrade string
	if userType == global.UserTypeAdmin {
		userGrade = "A"
	} else if grade == "" {
		userGrade = "C"
	} else {
		userGrade = grade
	}

	return userGrade
}

// GetProfile 取盈亏金额信息
// return: 盈亏金额, 交易次数， 交易金额，完成度
func (l *GetUserProfileLogic) GetProfile(req *proto.UserProfileReq) (int64, int, int64, float64) {
	p, err := l.svcCtx.SummaryRepo.GetUserProfit(l.ctx, req.GetCurDay(), req.GetUserId(), req.GetUserType(), int(req.GetAlgoId()))
	if err != nil {
		l.Logger.Error("GetUserProfit error:", err)
		return 0, 0, 0, 0
	}
	var progress float64
	if p.EntrustQty <= 0 {
		progress = 0.00
	} else {
		progress = float64(p.DealQty) / float64(p.EntrustQty) * 100
	}

	return int64(p.Perfit), p.OrderNum, p.OrderAmount, progress
}
