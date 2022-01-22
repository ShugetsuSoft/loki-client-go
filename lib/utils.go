package lib

import "unsafe"

func StringOut(bye []byte) string {
	return *(*string)(unsafe.Pointer(&bye))
}

func StringIn(strings string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&strings))
	return *(*[]byte)(unsafe.Pointer(&[3]uintptr{x[0], x[1], x[1]}))
}
