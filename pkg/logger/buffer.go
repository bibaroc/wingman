package logger

import (
	"io"
)

var _ io.Writer = (*cbuff)(nil)

type cbuff []byte

func (c *cbuff) Write(bytes []byte) (int, error) {
	*c = append(*c, bytes...)
	return len(bytes), nil
}
