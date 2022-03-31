// Package tools
/*
 Author: hawrkchen
 Date: 2022/3/22 15:56
 Desc: 字符串一些操作
*/
package tools

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
	for i :=0; i < len(vs); i++ {
		if s == vs[i] {
			vs = append(vs[:i], vs[i+1:]...)
			i--
		}
	}
	return vs
}
