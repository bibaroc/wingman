package fst_test

import (
	"testing"

	"github.com/bibaroc/wingman/pkg/fst"
	"github.com/google/go-cmp/cmp"
)

func TestSlice(t *testing.T) {
	str := "hi there, i'm ur string"
	want := []byte(str)
	got := fst.Slice(&str)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Simple conversion mismatch (-want +got):\n%s", diff)
	}
}

func BenchmarkSlice32(b *testing.B)      { benchmarkSlice(b, str(32)) }
func BenchmarkSlice64(b *testing.B)      { benchmarkSlice(b, str(64)) }
func BenchmarkSlice128(b *testing.B)     { benchmarkSlice(b, str(128)) }
func BenchmarkSlice256(b *testing.B)     { benchmarkSlice(b, str(256)) }
func BenchmarkSlc32(b *testing.B)        { benchmarkSlc(b, str(32)) }
func BenchmarkSlc64(b *testing.B)        { benchmarkSlc(b, str(64)) }
func BenchmarkSlc128(b *testing.B)       { benchmarkSlc(b, str(128)) }
func BenchmarkSlc256(b *testing.B)       { benchmarkSlc(b, str(256)) }
func BenchmarkCastBytea32(b *testing.B)  { benchmarkCastBytea(b, str(32)) }
func BenchmarkCastBytea64(b *testing.B)  { benchmarkCastBytea(b, str(64)) }
func BenchmarkCastBytea128(b *testing.B) { benchmarkCastBytea(b, str(128)) }
func BenchmarkCastBytea256(b *testing.B) { benchmarkCastBytea(b, str(256)) }

func benchmarkSlice(b *testing.B, str string) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := fst.Slice(&str)
		noop(got)
	}
}

func benchmarkSlc(b *testing.B, str string) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := fst.Slc(&str, len(str))
		noop(got)
	}
}

func benchmarkCastBytea(b *testing.B, str string) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got := []byte(str)
		noop(got)
	}
}
