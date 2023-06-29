// Package tools
/*
 Author: hawrkchen
 Date: 2022/3/31 15:27
 Desc:
*/
package tools

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
	"net"
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

// GzipEncode gzip 压缩，可配合json 使用
func GzipEncode(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write(input)
	if err != nil {
		_ = gzipWriter.Close()
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GzipDecode gzip 解压缩
func GzipDecode(input []byte) ([]byte, error) {
	bytesReader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = gzipReader.Close()
	}()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(gzipReader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalWithGzip(jsonData interface{}) []byte {
	dataAfterMarshal, _ := json.Marshal(jsonData)
	dataAfterGzip, err := GzipEncode(dataAfterMarshal)
	if err != nil {
		return nil
	}
	return dataAfterGzip
}

func UnmarshalWithGzip(msg []byte) (interface{}, error) {
	dataAfterDecode, err := GzipDecode(msg)
	if err != nil {
		return nil, err
	}
	var data interface{}
	if err := json.Unmarshal(dataAfterDecode, data); err != nil {
		return nil, err
	}
	return data, nil
}

// GetLocalIP 获取本地的IP地址
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
