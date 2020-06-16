package fst_test

import (
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
}
func TestConvertedStringIsNOTImmutable(t *testing.T) {
	c := []byte{'w', 'h', 'a', 't', ' ', 'i', 's', ' ', 'i', 'm', 'm', 'u', 't', 'a', 'b', 'i', 'l', 'i', 't', 'y'}
	want := string(c)
	got := fst.String(&c[0])

	// this capitalizes both Ms in iMMutability
	c[9], c[10] = 'M', 'M'
	if diff := cmp.Diff(want, got); diff == "" {
		t.Error("A string created this way should not be immutable if you modify the underlying array")
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
