package logic

import (
	mqservice "algo_assess/assess-mq-server/assessmqservice"
	"context"
	"encoding/json"

	"algo_assess/assess-api-server/internal/svc"
	"algo_assess/assess-api-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Region struct { // 区间配置
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
	Score int32   `json:"score"`
}

// AlgoFactor  算法因子配置
type AlgoFactor struct {
	FactorKey  string   `json:"factor_key"`
	FactorName string   `json:"factor_name"`
	Weight     int32    `json:"weight"`
	Regions    []Region `json:"regions"`
}

type AlgoConfigQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAlgoConfigQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AlgoConfigQueryLogic {
	return &AlgoConfigQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AlgoConfigQueryLogic) AlgoConfigQuery(req *types.QueryAlgoConfigReq) (resp *types.QueryAlgoConfigRsp, err error) {
	l.Logger.Infof("in AlgoConfigQuery, get req:%+v", req)
	qReq := &mqservice.GetAlgoConfigReq{
		ProfileType: req.QueryType,
	}
	rsp, err := l.svcCtx.AssessMQClient.GetAlgoConfig(l.ctx, qReq)
	if err != nil {
		l.Logger.Error("rpc GetAlgoConfig error:", err)
		return &types.QueryAlgoConfigRsp{
			Code:       361,
			Msg:        err.Error(),
			ConfigJson: "",
		}, nil
	}
	var algoJson string
	// 查询不到配置信息时，需要做一个兜底
	if rsp.GetConfig() == "" {
		algoJson = GetDefaultAlgoConfig(req.QueryType)
	} else {
		algoJson = rsp.GetConfig()
	}

	return &types.QueryAlgoConfigRsp{
		Code:       200,
		Msg:        rsp.GetMsg(),
		ConfigJson: algoJson,
	}, nil
}

func GetDefaultAlgoConfig(t int32) string {
	var rgs []Region
	r1 := Region{
		Lower: 0,
		Upper: 100,
		Score: 5,
	}
	rgs = append(rgs, r1)

	var ecs []AlgoFactor
	// 经济性
	ek := AlgoFactor{
		FactorKey:  "trade_vol",
		FactorName: "交易量",
		Weight:     10,
		Regions:    rgs,
	}
	ek2 := AlgoFactor{
		FactorKey:  "profit",
		FactorName: "盈亏收益",
		Weight:     10,
		Regions:    []Region{},
	}
	ek3 := AlgoFactor{
		FactorKey:  "profit_rate",
		FactorName: "收益率",
		Weight:     30,
		Regions:    []Region{},
	}
	ek4 := AlgoFactor{
		FactorKey:  "charge",
		FactorName: "手续费",
		Weight:     10,
		Regions:    []Region{},
	}
	//ek5 := AlgoFactor{
	//	FactorKey:  "cross_fee",
	//	FactorName: "流量费",
	//	Weight:     10,
	//	Regions:    []Region{},
	//}
	ek5 := AlgoFactor{
		FactorKey:  "cancel_rate",
		FactorName: "撤单率",
		Weight:     20,
		Regions:    []Region{},
	}
	ek6 := AlgoFactor{
		FactorKey:  "min_split_order",
		FactorName: "最小拆单单位",
		Weight:     10,
		Regions:    []Region{},
	}
	ek7 := AlgoFactor{
		FactorKey:  "deal_effi",
		FactorName: "成交效率",
		Weight:     10,
		Regions:    []Region{},
	}
	// 完成度
	ek8 := AlgoFactor{
		FactorKey:  "progress",
		FactorName: "完成度",
		Weight:     80,
		Regions:    rgs,
	}
	// 母单贴合度
	ek81 := AlgoFactor{
		FactorKey:  "algo_order_fit",
		FactorName: "母单贴合度",
		Weight:     10,
		Regions:    rgs,
	}
	// 成交量贴合度
	ek82 := AlgoFactor{
		FactorKey:  "trade_vol_fit",
		FactorName: "成交量贴合度",
		Weight:     10,
		Regions:    rgs,
	}
	// 风险度
	ek9 := AlgoFactor{
		FactorKey:  "min_jonit_rate",
		FactorName: "最小贴合度",
		Weight:     30,
		Regions:    rgs,
	}
	ek10 := AlgoFactor{
		FactorKey:  "withdraw_rate",
		FactorName: "回撤比例",
		Weight:     40,
		Regions:    []Region{},
	}
	// 绩效
	ek11 := AlgoFactor{
		FactorKey:  "vwap_dev",
		FactorName: "vwap滑点值",
		Weight:     70,
		Regions:    rgs,
	}
	// 绩效收益率
	ek111 := AlgoFactor{
		FactorKey:  "assess_profitrate",
		FactorName: "绩效收益率",
		Weight:     30,
		Regions:    rgs,
	}
	// 稳定性
	ek12 := AlgoFactor{
		FactorKey:  "vwap_std_dev",
		FactorName: "vwap滑点校准差",
		Weight:     30,
		Regions:    rgs,
	}
	ek13 := AlgoFactor{
		FactorKey:  "profit_rate_std",
		FactorName: "收益率标准差",
		Weight:     30,
		Regions:    []Region{},
	}
	ek14 := AlgoFactor{
		FactorKey:  "joint_rate",
		FactorName: "贴合度",
		Weight:     10,
		Regions:    []Region{},
	}
	ek15 := AlgoFactor{
		FactorKey:  "trade_vol_std_dev",
		FactorName: "成交量贴合度标准差",
		Weight:     20,
		Regions:    []Region{},
	}
	ek16 := AlgoFactor{
		FactorKey:  "time_std_dev",
		FactorName: "时间贴合度标准差",
		Weight:     10,
		Regions:    []Region{},
	}

	if t == 1 {
		ecs = append(ecs, ek, ek2, ek3, ek4, ek5, ek6, ek7)
	} else if t == 2 {
		ecs = append(ecs, ek8, ek81, ek82)
	} else if t == 3 {
		ecs = append(ecs, ek9, ek10, ek3)
	} else if t == 4 {
		ecs = append(ecs, ek11, ek111)
	} else if t == 5 {
		ecs = append(ecs, ek12, ek13, ek14, ek15, ek16)
	}
	out, err := json.Marshal(ecs)
	if err != nil {
		return ""
	}
	o := string(out)
	return o
}
