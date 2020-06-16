package fst

import (
	_ "unsafe"
)

//go:linkname String runtime.gostringnocopy
func String(b *byte) string
