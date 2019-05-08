package iotoolz

import (
	"bytes"
	"testing"

	bytetoolz "github.com/lthibault/toolz/pkg/bytes"

	"github.com/stretchr/testify/assert"
)

func TestCopyStream(t *testing.T) {
	src := bytes.NewBuffer([]byte("hello, world!"))
	dst := new(bytes.Buffer)

	n, err := CopyStreamBuffered(dst, src, bytetoolz.NewStreamBuffer(0))

	assert.Equal(t, len("hello, world!"), n)
	assert.Equal(t, len("hello, world!"), dst.Len())
	assert.NoError(t, err)
	assert.Equal(t, "hello, world!", dst.String())
}
