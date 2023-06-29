package logic

import (
	"algo_assess/pkg/tools"
	"context"

	"algo_assess/assess-rpc-server/internal/svc"
	"algo_assess/assess-rpc-server/proto"

	"github.com/zeromicro/go-zero/core/logx"
)

type TotalScoreRankingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTotalScoreRankingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TotalScoreRankingLogic {
	return &TotalScoreRankingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  综合评分排名列表
func (l *TotalScoreRankingLogic) TotalScoreRanking(in *proto.ScoreRankReq) (*proto.ScoreRankRsp, error) {
	// todo: add your logic here and delete this line
	l.Logger.Info("TotalScoreRanking get req:", in)
	var list []*proto.ScoreRankInfo
	rank := (in.GetPage() - 1) * in.GetLimit()
	var total int64
	score := 100 - (in.GetPage()-1)*in.GetLimit()
	if score < 1 {
		score = 1
	}
	if in.GetRankType() == 1 { // dashboard 排名
		out, count, err := l.svcCtx.SummaryRepo.GetTotalScoreRanking(l.ctx, in.GetDate(), in.GetUserId(), int(in.GetPage()), int(in.GetLimit()))
		if err != nil {
			l.Logger.Error("GetTotalScoreRanking error:", err)
			return &proto.ScoreRankRsp{
				Code: 208,
				Msg:  err.Error(),
			}, nil
		}
		for _, v := range out {
			rank++
			info := &proto.ScoreRankInfo{
				Ranking:  rank,
				AlgoName: v.AlgoName,
				Score:    int32(v.CumsumScore),
			}
			list = append(list, info)
		}
		total = count
	} else if in.GetRankType() == 2 { // 高阶股票排名
		out, count, err := l.svcCtx.SecurityRepo.GetRanking(l.ctx, int(in.GetPage()), int(in.GetLimit()))
		if err != nil {
			l.Logger.Error("Get Security Ranking error:", err)
			return &proto.ScoreRankRsp{
				Code: 208,
				Msg:  err.Error(),
			}, nil
		}
		for _, v := range out {
			rank++
			info := &proto.ScoreRankInfo{
				Ranking:  rank,
				AlgoName: "",
				Score:    score,
				SecId:    v.SecurityId,
				SecName:  tools.RMu0000(v.SecurityName),
				UserId:   "",
			}
			if score > 1 {
				score--
			}
			list = append(list, info)
			total = count
		}
	} else if in.GetRankType() == 3 { // 高阶用户排名
		out, count, err := l.svcCtx.UserInfoRepo.GetRanking(l.ctx, int(in.GetPage()), int(in.GetLimit()))
		if err != nil {
			l.Logger.Error("Get user info Ranking error:", err)
			return &proto.ScoreRankRsp{
				Code: 208,
				Msg:  err.Error(),
			}, nil
		}
		for _, v := range out {
			rank++
			info := &proto.ScoreRankInfo{
				Ranking:  rank,
				AlgoName: "",
				Score:    int32(score),
				SecId:    "",
				SecName:  "",
				UserId:   v.UserId,
				UserName: v.UserName,
			}
			if score > 1 {
				score--
			}
			list = append(list, info)
			total = count
		}
	} else {
		l.Logger.Error("unsupported rank type:", in.GetRankType())
		return &proto.ScoreRankRsp{
			Code:  220,
			Msg:   "unsupported rank type",
			Total: 0,
			Info:  nil,
		}, nil
	}

	rsp := &proto.ScoreRankRsp{
		Code:  200,
		Msg:   "success",
		Total: total,
		Info:  list,
	}

	return rsp, nil
}
