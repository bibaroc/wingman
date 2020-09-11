package fst

import (
	"unsafe"
	_ "unsafe"
)

func Slice(str *string) []byte {
	s := (*stringStruct)(unsafe.Pointer(str))
	sl := slice{array: s.str, len: s.len, cap: s.len}
	return *(*[]byte)(unsafe.Pointer(&sl))
}

func Slc(str *string, ln int) []byte {
	s := (*stringStruct)(unsafe.Pointer(str))
	sl := slice{array: s.str, len: ln, cap: ln}
	return *(*[]byte)(unsafe.Pointer(&sl))
}

type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
