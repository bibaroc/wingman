package fst

import (
	"unsafe"
	_ "unsafe"
)

//go:linkname String runtime.gostringnocopy
func String(b *byte) string

func Str(str *byte, ln int) string {
	ss := stringStruct{str: unsafe.Pointer(str), len: ln}
	s := *(*string)(unsafe.Pointer(&ss))
	return s
}

type stringStruct struct {
	str unsafe.Pointer
	len int
}
