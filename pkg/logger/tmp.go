package logger

import (
	_ "unsafe"
)

//go:linkname callers runtime.callers
func callers(skip int, pcbuf []uintptr) int

//go:linkname findfunc runtime.findfunc
func findfunc(pc uintptr) funcInfo

//go:linkname funcname runtime.funcname
func funcname(f funcInfo) string
