package iotoolz

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyStream(t *testing.T) {
	t.Run("Pipeline", func(t *testing.T) {
		src := bytes.NewBufferString("hello, world!")
		dst := new(bytes.Buffer)

		n, err := CopyStream(dst, src)

		assert.Equal(t, len("hello, world!"), n)
		assert.Equal(t, len("hello, world!"), dst.Len())
		assert.NoError(t, err)
		assert.Equal(t, "hello, world!", dst.String())
	})

	t.Run("Echo", func(t *testing.T) {
		p0, p1 := net.Pipe()
		go CopyStream(p1, p1)

		p0.Write([]byte("hello, world!"))

		b := make([]byte, 16)
		n, err := p0.Read(b)

		assert.Equal(t, len("hello, world!"), n)
		assert.NoError(t, err)
		assert.Equal(t, "hello, world!", string(b[:n]))
	})
}
