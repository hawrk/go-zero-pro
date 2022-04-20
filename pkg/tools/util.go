// Package tools
/*
 Author: hawrkchen
 Date: 2022/3/31 15:27
 Desc:
*/
package tools

import (
	"log"
	"runtime"
)

func PrintPanicStack(extras ...interface{}) {
	if x := recover(); x != nil {
		i := 0
		funcName, file, line, ok := runtime.Caller(i)

		for ok {
			log.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			i++
			funcName, file, line, ok = runtime.Caller(i)
		}
	}
}
