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

type ScoreRankingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewScoreRankingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ScoreRankingLogic {
	return &ScoreRankingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ScoreRankingLogic) ScoreRanking(req *types.RankingReq) (resp *types.RankingRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("get req:%+v", req)
	start := cast.ToInt64(time.Unix(req.Date, 0).Format(global.TimeFormatDay))
	q := &assessservice.ScoreRankReq{
		Date:     start,
		RankType: req.RankingType,
		UserId:   req.UserId,
		Page:     req.Page,
		Limit:    req.Limit,
	}
	p, err := l.svcCtx.AssessClient.TotalScoreRanking(l.ctx, q)
	if err != nil {
		l.Logger.Error("call rpc TotalScoreRanking error:", err)
		return &types.RankingRsp{
			Code: 209,
			Msg:  err.Error(),
		}, nil
	}
	var list []types.TotalScore
	for _, v := range p.GetInfo() {
		l := types.TotalScore{
			Ranking:  v.GetRanking(),
			AlgoName: v.GetAlgoName(),
			Score:    v.GetScore(),
			SecId:    v.GetSecId(),
			SecName:  v.GetSecName(),
			UserId:   v.GetUserId(),
			UserName: v.GetUserName(),
		}
		list = append(list, l)
	}

	rsp := &types.RankingRsp{
		Code:  200,
		Msg:   "success",
		Total: p.GetTotal(),
		Info:  list,
	}

	return rsp, nil
}
