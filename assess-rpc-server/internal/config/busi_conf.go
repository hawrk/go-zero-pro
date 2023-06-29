// Package config
/*
 Author: hawrkchen
 Date: 2022/8/2 10:47
 Desc: 算法五个维度的描述文案
*/
package config

import "fmt"

// 算法的五个维度文案描述， 分为四个档次，分别为：
//  [0-3)     极差
//  [3-6)     一般
//  [6-8)     较好
//  [8-10]    极好

// GetEconomyDesc 经济性各档次描述
func GetEconomyDesc(score int) (s string) {
	switch {
	case score < 3:
		s = fmt.Sprintf("该算法经济性评分为 %d 表现较差", score)
	case score >= 3 && score < 6:
		s = fmt.Sprintf("该算法经济性评分为 %d 表现一般", score)
	case score >= 6 && score < 8:
		s = fmt.Sprintf("该算法经济性评分为 %d 表现较好", score)
	case score >= 8:
		s = fmt.Sprintf("该算法经济性评分为 %d 表现极好", score)
	}
	return s
}

// GetProgressDesc 完成度各档次描述
func GetProgressDesc(score int) (s string) {
	switch {
	case score < 3:
		s = fmt.Sprintf("该算法完成度评分为 %d 表现较差", score)
	case score >= 3 && score < 6:
		s = fmt.Sprintf("该算法完成度评分为 %d 表现一般", score)
	case score >= 6 && score < 8:
		s = fmt.Sprintf("该算法完成度评分为 %d 表现较好", score)
	case score >= 8:
		s = fmt.Sprintf("该算法完成度评分为 %d 表现极好", score)
	}
	return s
}

// GetRiskDesc 风险度各档次描述
func GetRiskDesc(score int) (s string) {
	switch {
	case score < 3:
		s = fmt.Sprintf("该算法风险度评分为 %d 表现较差", score)
	case score >= 3 && score < 6:
		s = fmt.Sprintf("该算法风险度评分为 %d 表现一般", score)
	case score >= 6 && score < 8:
		s = fmt.Sprintf("该算法风险度评分为 %d 表现较好", score)
	case score >= 8:
		s = fmt.Sprintf("该算法风险度评分为 %d 表现极好", score)
	}
	return s
}

// GetAssessDesc 算法绩效各档次描述
func GetAssessDesc(score int) (s string) {
	switch {
	case score < 3:
		s = fmt.Sprintf("该算法绩效评分为 %d 表现较差", score)
	case score >= 3 && score < 6:
		s = fmt.Sprintf("该算法绩效评分为 %d 表现一般", score)
	case score >= 6 && score < 8:
		s = fmt.Sprintf("该算法绩效评分为 %d 表现较好", score)
	case score >= 8:
		s = fmt.Sprintf("该算法绩效评分为 %d 表现极好", score)
	}
	return s
}

// GetStabilityDesc 稳定性各档次描述
func GetStabilityDesc(score int) (s string) {
	switch {
	case score < 3:
		s = fmt.Sprintf("该算法稳定性评分为 %d 表现较差", score)
	case score >= 3 && score < 6:
		s = fmt.Sprintf("该算法稳定性评分为 %d 表现一般", score)
	case score >= 6 && score < 8:
		s = fmt.Sprintf("该算法稳定性评分为 %d 表现较好", score)
	case score >= 8:
		s = fmt.Sprintf("该算法稳定性评分为 %d 表现极好", score)
	}
	return s
}
