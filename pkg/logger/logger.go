package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

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
	fptrs []uintptr
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
	const SKIPCALLERS = 4
	const NOTFOUND = "n/a"

	n := runtime.Callers(SKIPCALLERS, l.fptrs)
	if n == 0 {
		return NOTFOUND
	}
	frame := runtime.CallersFrames(l.fptrs[:n])

	f, _ := frame.Next()
	fnName := f.Function
	for i := len(f.Function) - 1; i > 0; i-- {
		if os.IsPathSeparator(f.Function[i]) {
			fnName = f.Function[i+1:]
			break
		}
	}

	return fnName
}

func (l *Log) write(level LogLevel, args ...interface{}) error {
	callerInfo := l.caller()

	l.Lock()
	l.buff = l.buff[:0]

	l.output(level, callerInfo)

	fmt.Fprint(&l.buff, args...)
	_, err := l.out.Write(l.buff)

	l.Unlock()
	return err
}
func (l *Log) writeln(level LogLevel, args ...interface{}) error {
	callerInfo := l.caller()

	l.Lock()
	l.buff = l.buff[:0]

	l.output(level, callerInfo)

	fmt.Fprintln(&l.buff, args...)
	_, err := l.out.Write(l.buff)

	l.Unlock()
	return err
}
func (l *Log) writef(level LogLevel, format string, args ...interface{}) error {
	callerInfo := l.caller()

	l.Lock()
	l.buff = l.buff[:0]

	l.output(level, callerInfo)

	fmt.Fprintf(&l.buff, format, args...)
	_, err := l.out.Write(l.buff)

	l.Unlock()
	return err
}

func (l *Log) output(level LogLevel, callerInfo string) {
	l.buff = append(l.buff, level.String()...)

	l.buff = append(l.buff, ' ')
	l.buff = append(l.buff, callerInfo...)
	l.buff = append(l.buff, ' ')

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
		buff:  make([]byte, 0, 10*1<<10),
		fptrs: make([]uintptr, 1),
	}
}
