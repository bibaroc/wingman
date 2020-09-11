package fst_test

import (
	"crypto/rand"
	"testing"

	"github.com/bibaroc/wingman/pkg/fst"
	"github.com/google/go-cmp/cmp"
)

func TestSimpleConversion(t *testing.T) {
	c := []byte{'w', 'h', 'a', 't', ' ', 'i', 's', ' ', 'i', 'm', 'm', 'u', 't', 'a', 'b', 'i', 'l', 'i', 't', 'y'}
	want := string(c)
	got := fst.String(&c[0])
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Simple conversion mismatch (-want +got):\n%s", diff)
	}
	got = fst.Str(&c[0], len(c))
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Simple conversion mismatch (-want +got):\n%s", diff)
	}
}
func TestConvertedStringIsNOTImmutable(t *testing.T) {
	c := []byte{'w', 'h', 'a', 't', ' ', 'i', 's', ' ', 'i', 'm', 'm', 'u', 't', 'a', 'b', 'i', 'l', 'i', 't', 'y'}
	want := string(c)
	got := fst.String(&c[0])

	// this capitalizes both Ms in iMMutability
	c[9], c[10] = 'M', 'M'
	if diff := cmp.Diff(want, got); diff == "" {
		t.Errorf("A string created this way should not be immutable if you modify the underlying array\n")
	}
	got = fst.Str(&c[0], len(c))
	if diff := cmp.Diff(want, got); diff == "" {
		t.Errorf("A string created this way should not be immutable if you modify the underlying array\n")
	}
}

func TestConvertedResizing(t *testing.T) {
	c := []byte{'w', 'h', 'a', 't', ' ', 'i', 's', ' ', 'i', 'm', 'm', 'u', 't', 'a', 'b', 'i', 'l', 'i', 't', 'y', 0, 0, 0, 0, 0, 0, 0}

	c[7] = 0
	got := fst.String(&c[0])
	if diff := cmp.Diff("what is", got); diff != "" {
		t.Errorf("Setting a null byte should resize the string down, mismatch (-want +got):\n%s", diff)
	}

	c[7] = '-'
	got = fst.String(&c[0])
	if diff := cmp.Diff("what is-immutability", got); diff != "" {
		t.Errorf("Setting a previously null byte should resize the string up to the next one, mismatch (-want +got):\n%s", diff)
	}

	for i := 20; i < len(c); i++ {
		c[i] = 'a'
	}
	got = fst.String(&c[0])
	if diff := cmp.Diff("what is-immutabilityaaaaaaa", got); diff != "" {
		t.Errorf("Filling the array up to capacity should resize the string up, mismatch (-want +got):\n%s", diff)
	}
}

func BenchmarkString32(b *testing.B)      { benchmarkString(b, bytes(32)) }
func BenchmarkString64(b *testing.B)      { benchmarkString(b, bytes(64)) }
func BenchmarkString128(b *testing.B)     { benchmarkString(b, bytes(128)) }
func BenchmarkString256(b *testing.B)     { benchmarkString(b, bytes(256)) }
func BenchmarkStr32(b *testing.B)         { benchmarkStr(b, bytes(32)) }
func BenchmarkStr64(b *testing.B)         { benchmarkStr(b, bytes(64)) }
func BenchmarkStr128(b *testing.B)        { benchmarkStr(b, bytes(128)) }
func BenchmarkStr256(b *testing.B)        { benchmarkStr(b, bytes(256)) }
func BenchmarkCastString32(b *testing.B)  { benchmarkCastString(b, bytes(32)) }
func BenchmarkCastString64(b *testing.B)  { benchmarkCastString(b, bytes(64)) }
func BenchmarkCastString128(b *testing.B) { benchmarkCastString(b, bytes(128)) }
func BenchmarkCastString256(b *testing.B) { benchmarkCastString(b, bytes(256)) }

func benchmarkString(b *testing.B, arr []byte) {
	ln := len(arr)
	c := make([]byte, ln)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(c, arr)

		c[7] = 0
		got := fst.String(&c[0])
		noop(got)
		c[7] = '-'
		got = fst.String(&c[0])
		noop(got)

		for i := 0; i < len(c); i++ {
			c[i] = 'a'
		}
		got = fst.String(&c[0])
		noop(got)
		c = c[:cap(c)]
	}
}

func benchmarkStr(b *testing.B, arr []byte) {
	ln := len(arr)
	c := make([]byte, ln)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(c, arr)

		c[7] = 0
		got := fst.Str(&c[0], len(c))
		noop(got)
		c[7] = '-'
		got = fst.Str(&c[0], len(c))
		noop(got)

		for i := 0; i < len(c); i++ {
			c[i] = 'a'
		}
		got = fst.Str(&c[0], len(c))
		noop(got)
		c = c[:cap(c)]
	}
}

func benchmarkCastString(b *testing.B, arr []byte) {
	ln := len(arr)
	c := make([]byte, ln)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(c, arr)

		c[7] = 0
		got := string(c)
		noop(got)
		c[7] = '-'
		got = string(c)

		noop(got)

		for i := 0; i < len(c); i++ {
			c[i] = 'a'
		}
		got = string(c)
		noop(got)
		c = c[:cap(c)]
	}
}

func noop(v interface{}) {}
func bytes(count int64) []byte {
	b := make([]byte, count)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
func str(count int64) string {
	b := make([]byte, count)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return string(b)
}
