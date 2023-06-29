package logic

import (
	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"
	"algo_assess/assess-rpc-server/assessservice"
	"algo_assess/global"
	"algo_assess/mornano-rpc-server/mornanoservice"
	"algo_assess/mornano-rpc-server/proto"
	"context"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/threading"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserProfile 高阶分析:用户画像
func (l *UserProfileLogic) UserProfile(req *types.UserSummaryReq) (resp *types.UserSummaryRsp, err error) {
	// todo: add your logic here and delete this line
	l.Logger.Infof("in UserProfile, get req:%+v", req)
	if l.svcCtx.Config.WorkControl.EnableFakeMsg {
		return l.BuildVirtualRsp(req)
	}

	capitalCh := make(chan *mornanoservice.CapitailRsp)
	assessCh := make(chan *assessservice.UserProfileRsp)
	tlCh := make(chan *assessservice.TimeLineRsp)

	start := cast.ToInt64(time.Unix(req.StartTime, 0).Format(global.TimeFormatDay))
	//end := cast.ToInt64(time.Unix(req.EndTime, 0).Format(global.TimeFormatDay))

	var algoId int32
	var algoTypeName, algoName, provider string
	if req.AlgoName != "" {
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
	} else { // 如果没传入算法名称，则做一个兜底，默认返回一条
		l.Logger.Info("algo_name empty, get default algo....")
		dReq := &assessservice.DefaultReq{
			Scene:     1,
			UserId:    req.UserId,
			UserType:  int32(req.UserType),
			StartTime: start,
			EndTime:   start,
		}
		dRsp, err := l.svcCtx.AssessClient.GetDefaultAlgo(l.ctx, dReq)
		if err != nil {
			l.Logger.Error("rpc call GetDefaultAlgo error:", err)
			return &types.UserSummaryRsp{
				Code:         350,
				Msg:          err.Error(),
				Fund:         0,
				Profit:       0,
				TradeCnt:     0,
				CurTradeVol:  0,
				CurRollHold:  0,
				Progress:     0,
				LoginCnt:     0,
				UserGrade:    "",
				FundList:     []types.UserFund{},
				TotalScore:   0,
				Dimensions:   []types.UDimensionInfo{},
				AssessLine:   types.UDemensionLine{},
				ProgressLine: types.UDemensionLine{},
			}, nil

		}
		algoId = dRsp.GetAlgoId()
		algoTypeName = dRsp.GetAlgoTypeName()
		algoName = dRsp.GetAlgoName()
		provider = dRsp.GetProvider()
	}

	//1. redis取登陆次数
	cnt := l.GetLoginCount(req.UserId)
	//2. 取用户资金持仓信息
	threading.GoSafe(func() {
		mReq := &mornanoservice.CapitalReq{
			UserId: req.UserId,
		}
		mRsp, err := l.svcCtx.MornanoClient.GetUserCapital(l.ctx, mReq)
		if err != nil {
			l.Logger.Error("call mornano rpc GetUserCapital error:", err)
			capitalCh <- &mornanoservice.CapitailRsp{}
			return
		}
		capitalCh <- mRsp
	})

	// 3.取用户算法信息
	threading.GoSafe(func() {
		uReq := &assessservice.UserProfileReq{
			AlgoId:   algoId,
			AlgoName: req.AlgoName,
			UserId:   req.UserId,
			UserType: int32(req.UserType),
			CurDay:   start,
		}
		uRsp, err := l.svcCtx.AssessClient.GetUserProfile(l.ctx, uReq)
		if err != nil {
			l.Logger.Error("rpc call GetUserProfile error:", err)
			assessCh <- &assessservice.UserProfileRsp{}
			return
		}
		assessCh <- uRsp
	})

	//3.取算法实时绩效信息
	threading.GoSafe(func() {
		tlReq := &assessservice.TimeLineReq{
			LineType:     12,
			StartTime:    start,
			EndTime:      start,
			UserId:       req.UserId,
			UserType:     int32(req.UserType),
			AlgoId:       algoId,
			CrossDayFlag: false,
		}
		tlRsp, err := l.svcCtx.AssessClient.GetAlgoTimeLine(l.ctx, tlReq)
		if err != nil {
			l.Logger.Error("call assess rcp GetAlgoTimeLine error :", err)
			tlCh <- &assessservice.TimeLineRsp{}
			return
		}
		tlCh <- tlRsp
	})

	cRsp, uRsp, tlRsp := <-capitalCh, <-assessCh, <-tlCh

	return l.BuildUserProfileRsp(cnt, cRsp, uRsp, tlRsp, provider, algoTypeName, algoName), nil

}

// GetLoginCount 取登陆次数
func (l *UserProfileLogic) GetLoginCount(userId string) int {
	LoginKey := global.LoginCount + ":" + time.Now().Format(global.TimeFormatDay) + ":" + userId
	loginCnt, err := l.svcCtx.HRedisClient.Get(l.ctx, LoginKey).Result()
	if err != nil {
		l.Logger.Error("get login count error:", err)
	}
	l.Logger.Info("get loginCnt :", loginCnt)
	return cast.ToInt(loginCnt)
}

func (l *UserProfileLogic) BuildUserProfileRsp(loginCnt int, cRsp *mornanoservice.CapitailRsp, uRsp *assessservice.UserProfileRsp,
	tlRsp *assessservice.TimeLineRsp, provider, algoTypeName, algoName string) *types.UserSummaryRsp {
	// 拼时间线
	aTL, pTL := l.ParseTimeLine(tlRsp.Line)
	// 拼五个维度的评分
	var di []types.UDimensionInfo
	for _, v := range uRsp.Dimension {
		d := types.UDimensionInfo{
			ProfileType: v.GetProfileType(),
			Score:       v.GetProfileScore(),
			Desc:        v.GetProfileDesc(),
		}
		di = append(di, d)
	}
	cRsp.GetStockPosition()

	o := &types.UserSummaryRsp{
		Code:         200,
		Msg:          "success",
		Provider:     provider,
		AlgoTypeName: algoTypeName,
		AlgoName:     algoName,
		Fund:         cRsp.GetAvailable(), // 可用资金
		Profit:       float64(uRsp.GetProfit()) / 10000,
		TradeCnt:     uRsp.GetTradeCnt(),
		CurTradeVol:  float64(uRsp.GetTradeAmount()) / 10000,
		Progress:     uRsp.GetProgress(),
		LoginCnt:     int32(loginCnt),
		UserGrade:    uRsp.GetUserGrade(),
		FundList:     l.GetFundList(cRsp.GetStockPosition()),
		TotalScore:   uRsp.GetTotalScore(),
		Ranking:      uRsp.GetRanking(),
		Dimensions:   di,
		AssessLine:   aTL,
		ProgressLine: pTL,
	}
	return o
}

// GetFundList 拼装资金持仓列表
func (l *UserProfileLogic) GetFundList(in []*proto.StockPosition) []types.UserFund {
	var out []types.UserFund
	for _, v := range in {
		o := types.UserFund{
			SecId:   v.GetSecId(),
			SecName: v.GetSecName(),
			Hold:    v.GetMarketCap(),
			Cost:    v.GetCost(),
		}
		out = append(out, o)
	}
	return out
}

// ParseTimeLine 时间线绩效数据
func (l *UserProfileLogic) ParseTimeLine(line []*assessservice.TimeLine) (types.UDemensionLine, types.UDemensionLine) {
	var assTL []types.UTimeLine
	var pgTL []types.UTimeLine
	for _, v := range line {
		atl := types.UTimeLine{
			TimePoint: v.TimePoint,
			Score:     v.AssessScore,
		}
		ptl := types.UTimeLine{
			TimePoint: v.TimePoint,
			Score:     int32(v.Progress),
		}
		assTL = append(assTL, atl)
		pgTL = append(pgTL, ptl)
	}
	assess := types.UDemensionLine{
		ProfileType: 4,
		Point:       assTL,
	}

	progress := types.UDemensionLine{
		ProfileType: 2,
		Point:       pgTL,
	}
	return assess, progress
}

func (l *UserProfileLogic) BuildVirtualRsp(req *types.UserSummaryReq) (*types.UserSummaryRsp, error) {
	var fundList []types.UserFund
	f1 := types.UserFund{
		SecId:   "600519",
		SecName: "贵州茅台",
		Hold:    12678.90,
		Cost:    13679.30,
	}
	f2 := types.UserFund{
		SecId:   "000001",
		SecName: "平安银行",
		Hold:    2000.5,
		Cost:    3000.4,
	}
	f3 := types.UserFund{
		SecId:   "000034",
		SecName: "神州数码",
		Hold:    140.9,
		Cost:    120.6,
	}

	fundList = append(fundList, f1, f2, f3)

	var diList []types.UDimensionInfo
	d1 := types.UDimensionInfo{
		ProfileType: 1,
		Score:       8,
		Desc:        "该算法经济性评分为 8 表现较好",
	}
	d2 := types.UDimensionInfo{
		ProfileType: 2,
		Score:       6,
		Desc:        "该算法完成度评分为 6 表现一般",
	}
	d3 := types.UDimensionInfo{
		ProfileType: 3,
		Score:       9,
		Desc:        "该算法风险度评分为 9 表现极好",
	}
	d4 := types.UDimensionInfo{
		ProfileType: 4,
		Score:       3,
		Desc:        "该算法绩效评分为 3 表现较差",
	}
	d5 := types.UDimensionInfo{
		ProfileType: 5,
		Score:       10,
		Desc:        "该算法稳定性评分为 10 表现极好",
	}
	diList = append(diList, d1, d2, d3, d4, d5)

	var aLine []types.UTimeLine
	a1 := types.UTimeLine{
		TimePoint: "09:48",
		Score:     8,
	}
	a2 := types.UTimeLine{
		TimePoint: "10:36",
		Score:     6,
	}
	a3 := types.UTimeLine{
		TimePoint: "14:30",
		Score:     5,
	}
	aLine = append(aLine, a1, a2, a3)
	al := types.UDemensionLine{
		ProfileType: 4,
		Point:       aLine,
	}

	var pLine []types.UTimeLine
	p1 := types.UTimeLine{
		TimePoint: "10:03",
		Score:     4,
	}
	p2 := types.UTimeLine{
		TimePoint: "11:05",
		Score:     8,
	}
	p3 := types.UTimeLine{
		TimePoint: "14:30",
		Score:     6,
	}
	pLine = append(pLine, p1, p2, p3)
	pl := types.UDemensionLine{
		ProfileType: 2,
		Point:       pLine,
	}

	rsp := &types.UserSummaryRsp{
		Code:         200,
		Msg:          "success",
		Fund:         1024.00,
		Profit:       9465,
		TradeCnt:     25,
		CurTradeVol:  36,
		CurRollHold:  12356.90,
		Progress:     30,
		LoginCnt:     8,
		UserGrade:    "B",
		FundList:     fundList,
		TotalScore:   83,
		Dimensions:   diList,
		AssessLine:   al,
		ProgressLine: pl,
	}

	return rsp, nil
}
