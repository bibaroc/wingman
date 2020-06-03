package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"github.com/bibaroc/wingman/pkg"
)

type LogLevel int

const (
	FATAL   LogLevel = iota // [FATA]
	ERROR                   // [ERRO]
	WARNING                 // [WARN]
	INFO                    // [INFO]
	DEBUG                   // [DEBU]
	TRACE                   // [TRAC]
)

type CallerEnabled bool

const (
	WithCallerInfo = 1 << iota
)

var _ pkg.Logger = (*Log)(nil)

type Log struct {
	level LogLevel
	out   io.Writer
	flags int
	buff  cbuff
	b     sync.Pool
	sync.Mutex
}

func (l *Log) Panic(v ...interface{}) {
	panic(l.write(ERROR, v...))
}
func (l *Log) Panicf(format string, v ...interface{}) {
	panic(l.writef(ERROR, format, v...))
}
func (l *Log) Panicln(v ...interface{}) {
	panic(l.writeln(ERROR, v...))
}

func (l *Log) Fatal(v ...interface{}) {
	if err := l.write(ERROR, v...); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}
func (l *Log) Fatalf(format string, v ...interface{}) {
	if err := l.writef(ERROR, format, v...); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}
func (l *Log) Fatalln(v ...interface{}) {
	if err := l.writeln(ERROR, v...); err != nil {
		fmt.Println(err)
	}
	os.Exit(1)
}

func (l *Log) Error(v ...interface{}) {
	if l.level < ERROR {
		return
	}
	if err := l.write(ERROR, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Errorf(format string, v ...interface{}) {
	if l.level < ERROR {
		return
	}
	if err := l.writef(ERROR, format, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Errorln(v ...interface{}) {
	if l.level < ERROR {
		return
	}
	if err := l.writeln(ERROR, v...); err != nil {
		fmt.Println(err)
	}
}

func (l *Log) Warn(v ...interface{}) {
	if l.level < WARNING {
		return
	}
	if err := l.write(WARNING, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Warnf(format string, v ...interface{}) {
	if l.level < WARNING {
		return
	}
	if err := l.writef(WARNING, format, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Warnln(v ...interface{}) {
	if l.level < WARNING {
		return
	}
	if err := l.writeln(WARNING, v...); err != nil {
		fmt.Println(err)
	}
}

func (l *Log) Info(v ...interface{}) {
	if l.level < INFO {
		return
	}
	if err := l.write(INFO, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Infof(format string, v ...interface{}) {
	if l.level < INFO {
		return
	}
	if err := l.writef(INFO, format, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Infoln(v ...interface{}) {
	if l.level < INFO {
		return
	}
	if err := l.writeln(INFO, v...); err != nil {
		fmt.Println(err)
	}
}

func (l *Log) Debug(v ...interface{}) {
	if l.level < DEBUG {
		return
	}
	if err := l.write(DEBUG, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Debugf(format string, v ...interface{}) {
	if l.level < DEBUG {
		return
	}
	if err := l.writef(DEBUG, format, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Debugln(v ...interface{}) {
	if l.level < DEBUG {
		return
	}
	if err := l.writeln(DEBUG, v...); err != nil {
		fmt.Println(err)
	}
}

func (l *Log) Trace(v ...interface{}) {
	if l.level < TRACE {
		return
	}
	if err := l.write(TRACE, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Tracef(format string, v ...interface{}) {
	if l.level < TRACE {
		return
	}
	if err := l.writef(TRACE, format, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Traceln(v ...interface{}) {
	if l.level < TRACE {
		return
	}
	if err := l.writeln(TRACE, v...); err != nil {
		fmt.Println(err)
	}
}

func (l *Log) Print(v ...interface{}) {
	if err := l.write(ERROR, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Printf(format string, v ...interface{}) {
	if err := l.writef(ERROR, format, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) Println(v ...interface{}) {
	if err := l.writeln(ERROR, v...); err != nil {
		fmt.Println(err)
	}
}
func (l *Log) caller() string {
	if l.flags&WithCallerInfo == 0 {
		return "-"
	}
	const SKIPCALLERS = 3
	const NOTFOUND = "n/a"

	c := []uintptr{0}
	if callers(SKIPCALLERS, c) == 0 {
		return NOTFOUND
	}

	fnc := findfunc(c[0])
	if !fnc.valid() {
		return NOTFOUND
	}

	fnName := funcname(fnc)
	for i := len(fnName) - 1; i > 0; i-- {
		if os.IsPathSeparator(fnName[i]) {
			return fnName[i+1:]
		}
	}
	return fnName
}

func (l *Log) write(level LogLevel, args ...interface{}) error {
	callerInfo := l.caller()

	buff := l.b.Get().(*bytes.Buffer)
	buff.Reset()

	buff.WriteString(level.String())
	buff.WriteByte(' ')
	buff.WriteString(callerInfo)
	buff.WriteByte(' ')

	fmt.Fprint(buff, args...)

	l.Lock()
	_, err := l.out.Write(buff.Bytes())
	l.Unlock()

	l.b.Put(buff)
	return err
}
func (l *Log) writeln(level LogLevel, args ...interface{}) error {
	callerInfo := l.caller()

	buff := l.b.Get().(*bytes.Buffer)
	buff.Reset()

	buff.WriteString(level.String())
	buff.WriteByte(' ')
	buff.WriteString(callerInfo)
	buff.WriteByte(' ')

	fmt.Fprintln(buff, args...)

	l.Lock()
	_, err := l.out.Write(buff.Bytes())
	l.Unlock()

	l.b.Put(buff)

	return err
}
func (l *Log) writef(level LogLevel, format string, args ...interface{}) error {
	callerInfo := l.caller()

	buff := l.b.Get().(*bytes.Buffer)
	buff.Reset()

	buff.WriteString(level.String())
	buff.WriteByte(' ')
	buff.WriteString(callerInfo)
	buff.WriteByte(' ')

	fmt.Fprintf(buff, format, args...)
	l.Lock()
	_, err := l.out.Write(buff.Bytes())
	l.Unlock()

	l.b.Put(buff)
	return err
}

func NewLogger(
	level LogLevel,
	writer io.Writer,
	flags int,
) *Log {
	return &Log{
		level: level,
		out:   writer,
		flags: flags,
		b: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 10*1<<10))
			},
		},
	}
}

// //go:linkname findMe runtime.findfunc
// func findMe(pc uintptr) *funcInfo

// //go:nosplit
// func gostringnocopy(str *byte) string {
// 	ss := stringStruct{str: unsafe.Pointer(str), len: findnull(str)}
// 	s := *(*string)(unsafe.Pointer(&ss))
// 	return s
// }

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

type functab struct {
	entry   uintptr
	funcoff uintptr
}

type textsect struct {
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

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

type tflag uint8
type nameOff int32
type typeOff int32
type textOff int32

type name struct {
	bytes *byte
}
type imethod struct {
	name nameOff
	ityp typeOff
}
type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}
type ptabEntry struct {
	name nameOff
	typ  typeOff
}

// Information from the compiler about the layout of stack frames.
type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}
type findfuncbucket struct {
	idx        uint32
	subbuckets [16]byte
}

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
type funcID uint8

//go:nosplit
func findnull(s *byte) int {
	if s == nil {
		return 0
	}

	// pageSize is the unit we scan at a time looking for NULL.
	// It must be the minimum page size for any architecture Go
	// runs on. It's okay (just a minor performance loss) if the
	// actual system page size is larger than this value.
	const pageSize = 4096

	offset := 0
	ptr := unsafe.Pointer(s)
	// IndexByteString uses wide reads, so we need to be careful
	// with page boundaries. Call IndexByteString on
	// [ptr, endOfPage) interval.
	safeLen := int(pageSize - uintptr(ptr)%pageSize)

	for {
		t := *(*string)(unsafe.Pointer(&stringStruct{ptr, safeLen}))
		// Check one page at a time.
		if i := strings.IndexByte(t, 0); i != -1 {
			return offset + i
		}
		// Move to next page
		ptr = unsafe.Pointer(uintptr(ptr) + uintptr(safeLen))
		offset += safeLen
		safeLen = pageSize
	}
}

type stringStruct struct {
	str unsafe.Pointer
	len int
}
