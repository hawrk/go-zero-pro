// Package tools
/*
 Author: hawrkchen
 Date: 2022/3/22 15:52
 Desc: 产生随机字符串
*/
package tools

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// Krand 随机字符串
func Krand(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

func IntToBytes(n uint64) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, n)
	return bytebuf.Bytes()
}

// MinInt 取两个整型中的最小值
func MinInt(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func MinFloat(x, y float64) float64 {
	if x < y {
		return x
	}
	return y
}

// MaxInt 取两个整型中的最大值
func MaxInt(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// BinarySearch 二分查找
func BinarySearch(arr []int, start, end int, ele int) int {
	if start > end {
		return 0
	}
	middle := (start + end) / 2
	if arr[ele] > ele {
		BinarySearch(arr, start, middle-1, ele)
	} else if arr[ele] < ele {
		BinarySearch(arr, middle+1, end, ele)
	} else {
		return middle
	}
	return 0
}

// GetRanking 取排名, arr 为倒序排列的分数列表，  score为分数， 返回值为第几位
func GetRanking(arr []int, score int) int32 {
	for k, v := range arr {
		if score == v {
			return int32(k) + 1
		}
	}
	return 0
}
