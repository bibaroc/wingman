package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
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

const (
	WithCallerInfo = 1 << iota
)

var (
	_ pkg.Logger = (*Log)(nil)
)

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
	for i := len(fnName) - 1; i >= 0; i-- {
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

func NewDebugLogger() *Log {
	return NewLogger(
		DEBUG,
		os.Stdout,
		WithCallerInfo,
	)
}
