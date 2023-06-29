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

type AlgoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoListLogic {
	return &AlgoListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AlgoList 算法 dashboard 根据算法类型查询列表
func (l *AlgoListLogic) AlgoList(req *types.AlgoListReq) (resp *types.AlgoListRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("AlgoList get req:%+v", req)
	//根据算法类型名称反查一下算法类型ID
	tReq := &assessservice.ChooseAlgoReq{
		ChooseType:   6,
		AlgoTypeName: req.AlgoTypeName,
	}
	tRsp, err := l.svcCtx.AssessClient.ChooseAlgoInfo(l.ctx, tReq)
	if err != nil {
		l.Logger.Error("call rpc ChooseAlgoInfo error:", err)
		return &types.AlgoListRsp{
			Code: 205,
			Msg:  err.Error(),
		}, nil
	}

	listCh := make(chan *assessservice.AlgoOrderRsp)
	topCh := make(chan *assessservice.MultiAlgoRsp)

	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	//根据算法类型取列表
	go func() {
		lReq := &assessservice.AlgoOrderReq{
			Date:     start,
			AlgoType: tRsp.GetAlgoType(),
			UserId:   req.UserId,
			UserType: int32(req.UserType),
			Page:     req.Page,
			Limit:    req.Limit,
		}
		lRsp, err := l.svcCtx.AssessClient.GetAlgoOrderSummary(l.ctx, lReq)
		if err != nil {
			l.Logger.Error("call rpc GetAlgoOrderSummary error:", err)
			listCh <- &assessservice.AlgoOrderRsp{}
			return
		}
		listCh <- lRsp
	}()

	// 取top4
	go func() {
		// 只有第一页的请求才会调用这个接口
		if req.Page != 1 {
			topCh <- &assessservice.MultiAlgoRsp{}
			return
		}
		mReq := &assessservice.MultiAlgoReq{
			Date:      start,
			AlgoType:  tRsp.GetAlgoType(),
			UserId:    req.UserId,
			UserType:  int32(req.UserType),
			SceneType: 1,
		}
		mRsp, err := l.svcCtx.AssessClient.GetMultiAlgoAssess(l.ctx, mReq)
		if err != nil {
			l.Logger.Error("call rpc GetMultiAlgoAssess error:", err)
			topCh <- &assessservice.MultiAlgoRsp{}
			return
		}
		topCh <- mRsp
	}()
	listRsp, topRsp := <-listCh, <-topCh
	resp = l.BuildSummaryList(listRsp, topRsp)

	return resp, nil
}

func (l *AlgoListLogic) BuildSummaryList(aRsp *assessservice.AlgoOrderRsp, tRsp *assessservice.MultiAlgoRsp) *types.AlgoListRsp {
	var list []types.AlgoListInfo
	for _, v := range aRsp.GetInfo() {
		l := types.AlgoListInfo{
			Provider:   v.GetProvider(),
			UserCnt:    v.GetUserNum(),
			TradeVol:   float64(v.GetTotalTradeAmount()) / 10000,
			ProfitRate: v.GetProfitRate() * 100,
			OrderCnt:   v.GetOrderNum(),
			Side: types.DBTradeSide{
				BuyRate:  v.GetBuyRate() * 100,
				SellRate: v.GetSellRate() * 100,
			},
		}
		list = append(list, l)
	}
	if len(list) == 0 { // 适配前端，slice类型都返回空值，不返回nil值
		list = []types.AlgoListInfo{}
	}

	// 拼top4
	var sList []types.AlgoAccessInfo
	for _, v := range tRsp.GetSummary() {
		var tls []types.MTTimeLine
		for _, k := range v.GetTl() {
			l := types.MTTimeLine{
				TimePoint: k.GetTimePoint(),
				Score:     k.GetAssessScore(),
			}
			tls = append(tls, l)
		}
		tl := types.AlgoAccessInfo{
			AlgoName:   v.GetAlgoName(),
			TotalScore: v.GetTotalScore(),
			TL:         tls,
		}
		sList = append(sList, tl)
	}
	if len(sList) == 0 {
		sList = []types.AlgoAccessInfo{}
	}

	rsp := &types.AlgoListRsp{
		Code:  200,
		Msg:   "success",
		Total: aRsp.GetTotal(),
		List:  list,
		Infos: sList,
	}

	return rsp
}
