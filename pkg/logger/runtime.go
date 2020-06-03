package logger

import (
	"runtime"
	"unsafe"
)

// there functions are implemented in the runtime package
// here I prefer to use the functions directly in order to speed up caller lookup
// as the public functionality exposed by runtime pkg is ~30% slower than using the methods directly

// callers will populate pcbuf with pcs of calling functions
//  returns the total number of pc read
//go:linkname callers runtime.callers
func callers(skip int, pcbuf []uintptr) int

//go:linkname findfunc runtime.findfunc
func findfunc(pc uintptr) funcInfo

//go:linkname funcname runtime.funcname
func funcname(f funcInfo) string

// funcInfo is copy-pasted from src/runtime/symtab.go
//  methods are copied becouse they influence the memory layout
type funcInfo struct {
	*_func
	datap *moduledata
}

func (f funcInfo) valid() bool {
	return f._func != nil
}

func (f funcInfo) _Func() *runtime.Func {
	return (*runtime.Func)(unsafe.Pointer(f._func))
}

// _func is copy-pasted from src/runtime/runtime2.go
type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args        int32  // in/out args size
	deferreturn uint32 // offset of start of a deferreturn call instruction from entry, if any.

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	funcID    funcID  // set for certain special runtime functions
	_         [2]int8 // unused
	nfuncdata uint8   // must be last
}

// funcID is copy-pasted from src/runtime/symtab.go
type funcID uint8

// moduledata is copy-pasted from src/runtime/symtab.go
type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []*itab

	ptab []ptabEntry

	pluginpath string
	pkghashes  []modulehash

	modulename   string
	modulehashes []modulehash

	hasmain uint8 // 1 if module contains the main function, 0 otherwise

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]*_type // offset to *_rtype in previous module

	bad bool // module failed to load and should be ignored

	next *moduledata
}

// functab is copy-pasted from src/runtime/symtab.go
type functab struct {
	entry   uintptr
	funcoff uintptr
}

// textsect is copy-pasted from src/runtime/symtab.go
type textsect struct {
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

// itab is copy-pasted from src/runtime/runtime2.go
type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

// interfacetype is copy-pasted from src/runtime/type.go
type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

// _type is copy-pasted from src/runtime/type.go
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}

// tflag nameOff typeOff are copy-pasted from src/runtime/type.go
type tflag uint8
type nameOff int32
type typeOff int32

// name is copy-pasted from src/runtime/type.go
type name struct {
	bytes *byte
}

// imethod is copy-pasted from src/runtime/type.go
type imethod struct {
	name nameOff
	ityp typeOff
}

// ptanEntry is copy-pasted from src/runtime/plugin.go
type ptabEntry struct {
	name nameOff
	typ  typeOff
}

// modulehash is copy-pasted from src/runtime/symtab.go
type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}

// bitvector is copy-pasted from src/runtime/stack.go
type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}
