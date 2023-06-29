// Package tools
/*
 Author: hawrkchen
 Date: 2022/3/22 15:56
 Desc: 字符串一些操作
*/
package tools

import (
	"algo_assess/global"
	"encoding/json"
	"github.com/spf13/cast"
	"reflect"
	"sort"
	"strconv"
	"unsafe"
)

// 去掉 u0000, 保留正常的空格
func RMu0000(s string) string {
	str := make([]rune, 0, len(s))
	for _, v := range []rune(s) {
		if v == 0 {
			continue
		}
		str = append(str, v)
	}
	return string(str)
}

// 删除切片指定元素
func DelSlinceElem(vs []string, s string) []string {
	for i := 0; i < len(vs); i++ {
		if s == vs[i] {
			vs = append(vs[:i], vs[i+1:]...)
			i--
		}
	}
	return vs
}

// DeleteU64SliceElms 从 []uint64 过滤指定元素。注意：不修改原切片。
func DeleteU64SliceElms(i []uint64, elms ...uint64) []uint64 {
	// 构建 map set。
	m := make(map[uint64]struct{}, len(elms))
	for _, v := range elms {
		m[v] = struct{}{}
	}
	// 创建新切片，过滤掉指定元素。
	t := make([]uint64, 0, len(i))
	for _, v := range i {
		if _, ok := m[v]; !ok {
			t = append(t, v)
		}
	}
	return t
}

func InSlice(items []int, item int) bool {
	for _, v := range items {
		if item == v {
			return true
		}
	}
	return false
}

// 零拷贝 string 转byte
func String2Bytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: stringHeader.Data,
		Len:  stringHeader.Len,
		Cap:  stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

// 零拷贝 byte 转 string
func Bytes2String(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer((&b)))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}

// Json2Str 结构体类型直接转换成String类型，
// 注意这里直接忽略了Json失败的case
func Json2Str(i interface{}) string {
	v, _ := json.Marshal(i)
	return Bytes2String(v)
}

func Str2FundRate(str string) global.FundRate {
	var fr global.FundRate
	_ = json.Unmarshal([]byte(str), &fr)
	return fr
}

func Str2TradeVolDir(s string) global.TradeVolDirect {
	var tv global.TradeVolDirect
	_ = json.Unmarshal([]byte(s), &tv)
	return tv
}

func Str2StockType(s string) global.StockType {
	var st global.StockType
	_ = json.Unmarshal([]byte(s), &st)
	return st
}

func Str2TradeVolRate(s string) global.TradeVolRate {
	var tr global.TradeVolRate
	_ = json.Unmarshal([]byte(s), &tr)
	return tr
}

// Hex2Binary 十六进制字符串转化为二进制字符串
func Hex2Binary(x string) string {
	base, _ := strconv.ParseInt(x, 16, 64)
	b := strconv.FormatInt(base, 2)
	return b
}

// GetDoubleTimePoint 根据传入的时间数组，经过排序后，输出最小时间和最大时间点
// 入参 202306151041
func GetDoubleTimePoint(t []int64) (int64, int64) {
	var strKey []string
	m := make(map[string]struct{})
	for _, v := range t {
		d := cast.ToString(v)[:8]
		if _, exist := m[d]; !exist {
			m[d] = struct{}{}
			strKey = append(strKey, d) // 时间精确到天
		}
	}
	//fmt.Println("get len:", len(strKey))
	if len(strKey) == 0 {
		return 0, 0
	} else if len(strKey) == 1 { // 只有一条，即是当天的数据
		start := TimeStr2TimeUnix(strKey[0] + "0000")
		end := TimeStr2TimeUnix(strKey[0] + "2359")
		return start, end
	}
	// 有2个以上元素，说明跨天了
	//fmt.Println("before sort:", strKey)
	sort.Strings(strKey)
	//fmt.Println("after sort:", strKey)
	// 取第一条和最后一条
	start := TimeStr2TimeUnix(strKey[0] + "0000")
	end := TimeStr2TimeUnix(strKey[len(strKey)-1] + "2359")
	return start, end
}
