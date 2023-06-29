package util

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

func BytesExtend(data []byte, size int) []byte {
	if size > len(data) {
		var ext []byte = make([]byte, size-len(data))
		var buffer bytes.Buffer
		buffer.Write(data)
		buffer.Write(ext)
		return buffer.Bytes()
	}
	return data
}

func StringToBytes(s string) []byte {
	tmp1 := (*[2]uintptr)(unsafe.Pointer(&s))
	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}
	return *(*[]byte)(unsafe.Pointer(&tmp2))
}

func byteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

//整形转换成字节
func Int16ToBytes(n int16) []byte {
	x := int16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt16(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int16
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}

//整形转换成字节
func Int32ToBytes(n int32) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt32(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}

//整形转换成字节
func Int64ToBytes(n int64, err error) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt64(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return int(x)
}
