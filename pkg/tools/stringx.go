// Package tools
/*
 Author: hawrkchen
 Date: 2022/3/22 15:56
 Desc: 字符串一些操作
*/
package tools

import (
	"reflect"
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
