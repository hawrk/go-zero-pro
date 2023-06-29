// Package main
/*
 Author: hawrkchen
 Date: 2023/3/2 18:11
 Desc:
*/
package main

import (
	"encoding/json"
	"fmt"
)

type Region struct { // 区间配置
	Lower float64 `json:"lower"`
	Upper float64 `json:"upper"`
	Score int32   `json:"score"`
}

type FactorConf struct {
	Weight  int32    `json:"weight"`
	Regions []Region `json:"regions"`
}

type EcoJson struct {
	ProfileType int           `json:"profile_type"`
	Factor      []EconomyConf `json:"factors"`
}

// EconomyConf  经济性算法配置
type EconomyConf struct {
	FactorName string   `json:"factor_name"`
	Weight     int32    `json:"weight"`
	Regions    []Region `json:"regions"`
}

// ProgressConf 完成度算法配置
type ProgressConf struct {
	Progress FactorConf `json:"progress"`
}

// RiskConf 风险度算法配置
type RiskConf struct {
	MinJoion     FactorConf `json:"min_joion"`
	WithDrawRate FactorConf `json:"with_draw_rate"`
	ProfitRate   FactorConf `json:"profit_rate"`
}

// AssessConf 绩效算法配置
type AssessConf struct {
	VwapDev    FactorConf `json:"vwap_dev"`
	ProfitRate FactorConf `json:"profit_rate"`
}

// StableConf 稳定性算法配置
type StableConf struct {
	VwapDevStd    FactorConf `json:"vwap_dev_std"`
	ProfitRateStd FactorConf `json:"profit_rate_std"`
	JoinRate      FactorConf `json:"join_rate"`
}

func BuildEcoJson() []EconomyConf {
	var rgs []Region
	r1 := Region{
		Lower: 10,
		Upper: 50,
		Score: 2,
	}
	r2 := Region{
		Lower: 50,
		Upper: 80,
		Score: 5,
	}
	r3 := Region{
		Lower: 80,
		Upper: 100,
		Score: 8,
	}
	rgs = append(rgs, r1, r2, r3)

	//var ej EcoJson
	var ecs []EconomyConf
	ek := EconomyConf{
		FactorName: "交易量",
		Weight:     20,
		Regions:    rgs,
	}
	ek2 := EconomyConf{
		FactorName: "盈亏收益",
		Weight:     30,
		Regions:    rgs,
	}
	ek3 := EconomyConf{
		FactorName: "收益率",
		Weight:     40,
		Regions:    rgs,
	}
	ek4 := EconomyConf{
		FactorName: "手续费",
		Weight:     50,
		Regions:    rgs,
	}
	ek5 := EconomyConf{
		FactorName: "流量费",
		Weight:     60,
		Regions:    rgs,
	}
	ek6 := EconomyConf{
		FactorName: "撤单率",
		Weight:     70,
		Regions:    rgs,
	}
	ek7 := EconomyConf{
		FactorName: "最小拆单单位",
		Weight:     80,
		Regions:    rgs,
	}
	ecs = append(ecs, ek)
	ecs = append(ecs, ek2)
	ecs = append(ecs, ek3)
	ecs = append(ecs, ek4)
	ecs = append(ecs, ek5)
	ecs = append(ecs, ek6)
	ecs = append(ecs, ek7)
	return ecs
}

func main() {

	tvc := BuildEcoJson()
	out, err := json.Marshal(tvc)
	if err != nil {
		fmt.Println("marshal ec error:", err)
		return
	}
	o := string(out)
	fmt.Println("ecno json str:", o)

	/*

		pc := BuildProgress()
		out2, err := json.Marshal(pc)
		if err != nil {
			fmt.Println("marshal ec error:", err)
			return
		}
		o2 := string(out2)
		fmt.Println("progress json str:", o2)

		rc := BuildRiskConf()
		out3, err := json.Marshal(rc)
		if err != nil {
			fmt.Println("marshal ec error:", err)
			return
		}
		o3 := string(out3)
		fmt.Println("risk json str:", o3)

		ac := BuildAssessConf()
		out4, err := json.Marshal(ac)
		if err != nil {
			fmt.Println("marshal ec error:", err)
			return
		}
		o4 := string(out4)
		fmt.Println("assess json str:", o4)

		sc := BuildStableConf()
		out5, err := json.Marshal(sc)
		if err != nil {
			fmt.Println("marshal ec error:", err)
			return
		}
		o5 := string(out5)
		fmt.Println("stable json str:", o5)

	*/
}

/*
func BuildEconomy() EconomyConf {
	var rgs []Region
	r1 := Region{
		Lower: 10,
		Upper: 50,
		Score: 2,
	}
	r2 := Region{
		Lower: 50,
		Upper: 80,
		Score: 5,
	}
	r3 := Region{
		Lower: 80,
		Upper: 100,
		Score: 8,
	}
	rgs = append(rgs, r1, r2, r3)
	fc := FactorConf{
		Weight:  20,
		Regions: rgs,
	}

	var rpgs []Region
	rp1 := Region{
		Lower: 10,
		Upper: 50,
		Score: 2,
	}
	rp2 := Region{
		Lower: 50,
		Upper: 80,
		Score: 5,
	}
	rp3 := Region{
		Lower: 80,
		Upper: 100,
		Score: 8,
	}
	rpgs = append(rpgs, rp1, rp2, rp3)
	pfc := FactorConf{
		Weight:  30,
		Regions: rpgs,
	}

	var rprgs []Region
	rpr1 := Region{
		Lower: 10,
		Upper: 50,
		Score: 2,
	}
	rpr2 := Region{
		Lower: 50,
		Upper: 80,
		Score: 5,
	}
	rpr3 := Region{
		Lower: 80,
		Upper: 100,
		Score: 8,
	}
	rprgs = append(rprgs, rpr1, rpr2, rpr3)
	prfc := FactorConf{
		Weight:  25,
		Regions: rpgs,
	}

	var rcgs []Region
	rc1 := Region{
		Lower: 10,
		Upper: 50,
		Score: 2,
	}
	rc2 := Region{
		Lower: 50,
		Upper: 80,
		Score: 5,
	}
	rc3 := Region{
		Lower: 80,
		Upper: 100,
		Score: 8,
	}
	rcgs = append(rcgs, rc1, rc2, rc3)
	rfc := FactorConf{
		Weight:  25,
		Regions: rpgs,
	}

	ec := EconomyConf{
		FactorName: "",
		Weight:     0,
		Regions:    nil,
	}
	return ec
}
*/

func BuildProgress() ProgressConf {
	var rgs []Region
	r1 := Region{
		Lower: 10,
		Upper: 30,
		Score: 3,
	}
	r2 := Region{
		Lower: 30,
		Upper: 50,
		Score: 5,
	}
	r3 := Region{
		Lower: 50,
		Upper: 70,
		Score: 7,
	}
	r4 := Region{
		Lower: 70,
		Upper: 100,
		Score: 10,
	}
	rgs = append(rgs, r1, r2, r3, r4)
	fc := FactorConf{
		Weight:  100,
		Regions: rgs,
	}
	pc := ProgressConf{
		Progress: fc,
	}
	return pc
}

func BuildRiskConf() RiskConf {
	var rcgs []Region
	r1 := Region{
		Lower: 0.1,
		Upper: 0.5,
		Score: 3,
	}
	r2 := Region{
		Lower: 0.5,
		Upper: 0.8,
		Score: 5,
	}

	rcgs = append(rcgs, r1, r2)
	fc := FactorConf{
		Weight:  40,
		Regions: rcgs,
	}
	pc := RiskConf{
		MinJoion:     fc,
		WithDrawRate: FactorConf{},
		ProfitRate:   FactorConf{},
	}
	return pc
}

func BuildAssessConf() AssessConf {
	var rcgs []Region
	r1 := Region{
		Lower: 0.1,
		Upper: 0.5,
		Score: 3,
	}
	r2 := Region{
		Lower: 0.5,
		Upper: 0.8,
		Score: 5,
	}

	rcgs = append(rcgs, r1, r2)
	fc := FactorConf{
		Weight:  40,
		Regions: rcgs,
	}
	pc := AssessConf{
		VwapDev:    fc,
		ProfitRate: FactorConf{},
	}
	return pc
}

func BuildStableConf() StableConf {
	var rcgs []Region
	r1 := Region{
		Lower: 0.1,
		Upper: 0.5,
		Score: 3,
	}
	r2 := Region{
		Lower: 0.5,
		Upper: 0.8,
		Score: 5,
	}

	rcgs = append(rcgs, r1, r2)
	fc := FactorConf{
		Weight:  40,
		Regions: rcgs,
	}
	pc := StableConf{
		VwapDevStd:    fc,
		ProfitRateStd: FactorConf{},
		JoinRate:      FactorConf{},
	}
	return pc
}
