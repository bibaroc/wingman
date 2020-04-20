package logger_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/bibaroc/wingman/pkg/logger"
)

const (
	short  = "this will be fast"
	medium = "this will be fast"
	large  = "this will be fast"
)

func BenchmarkLoggerCustom(b *testing.B) {
	log := logger.NewLogger(logger.ERROR, ioutil.Discard, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Errorln(short)
		b.SetBytes(int64(len(short)))
	}
}
func BenchmarkLoggerStd(b *testing.B) {
	log := log.New(ioutil.Discard, "std", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println(short)
		b.SetBytes(int64(len(short)))
	}
}
func BenchmarkLoggerCustomWithCallerInfo(b *testing.B) {
	log := logger.NewLogger(logger.ERROR, ioutil.Discard, logger.WithCallerInfo)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Errorln(short)
		b.SetBytes(int64(len(short)))
	}
}
func BenchmarkLoggerStdWithCallerInfo(b *testing.B) {
	log := log.New(ioutil.Discard, "std", log.Lshortfile)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println(short)
		b.SetBytes(int64(len(short)))
	}
}
func BenchmarkLoggerFmt(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Fprintln(ioutil.Discard, short)
		b.SetBytes(int64(len(short)))
	}
}
