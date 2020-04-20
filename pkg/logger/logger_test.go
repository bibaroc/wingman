package logger_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/bibaroc/wingman/pkg/logger"
	pglog "github.com/go-playground/log/v7"
	"github.com/go-playground/log/v7/handlers/console"
)

func BenchmarkLoggerParallelCustom(b *testing.B) {
	log := logger.NewLogger(logger.ERROR, ioutil.Discard, 0)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Errorln(short)
			b.SetBytes(int64(len(short)))
		}

	})
}
func BenchmarkLoggerParallelStd(b *testing.B) {
	log := log.New(ioutil.Discard, "std", 0)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Println(short)
			b.SetBytes(int64(len(short)))
		}

	})
}

func BenchmarkLoggerParallelPlayground(b *testing.B) {
	csole := console.New(false)
	csole.SetDisplayColor(false)
	csole.SetWriter(ioutil.Discard)
	pglog.AddHandler(csole, pglog.AllLevels...)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pglog.Debug(short)
			b.SetBytes(int64(len(short)))
		}
	})
}
func BenchmarkLoggerParallelCustomWithCallerInfo(b *testing.B) {
	log := logger.NewLogger(logger.ERROR, ioutil.Discard, logger.WithCallerInfo)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Errorln(short)
			b.SetBytes(int64(len(short)))
		}
	})
}
func BenchmarkLoggerParallelStdWithCallerInfo(b *testing.B) {
	log := log.New(ioutil.Discard, "std", log.Lshortfile)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Println(short)
			b.SetBytes(int64(len(short)))
		}
	})
}
func BenchmarkLoggerParallelFmt(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Fprintln(ioutil.Discard, short)
			b.SetBytes(int64(len(short)))
		}

	})
}
