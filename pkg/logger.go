package pkg

type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})

	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})

	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Debugln(v ...interface{})

	Trace(v ...interface{})
	Tracef(format string, v ...interface{})
	Traceln(v ...interface{})

	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Infoln(v ...interface{})

	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Warnln(v ...interface{})

	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Errorln(v ...interface{})
}
