// Package tools
/*
 Author: hawrkchen
 Date: 2022/3/22 15:54
 Desc: 元转分， 分转元
*/
package tools

import "github.com/shopspring/decimal"

var oneHundredDecimal decimal.Decimal = decimal.NewFromInt(100)

// 分转元
func Fen2Yuan(fen int64) float64 {
	y, _ := decimal.NewFromInt(fen).Div(oneHundredDecimal).Truncate(2).Float64()
	return y
}

// 元转分
func Yuan2Fen(yuan float64) int64 {
	f, _ := decimal.NewFromFloat(yuan).Mul(oneHundredDecimal).Truncate(0).Float64()
	return int64(f)

}

// MulPercent float64运算都会出现精度问题，进行百分比转换时，*100，需要用该函数进行转换
func MulPercent(in float64) float64 {
	f, _ := decimal.NewFromFloat(in).Mul(oneHundredDecimal).Float64()
	return f
}
