package iotoolz

import (
	"bytes"
	"io"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	// t.Run("Small", func(t *testing.T) {
	// 	rb := NewRingBuffer(16)

	// 	t.Run("Write", func(t *testing.T) {
	// 		n, err := rb.Write([]byte("hello"))
	// 		assert.Equal(t, 5, n)
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, []byte("hello"), rb.b[1:6]) // NOTE:  off by one in ring buffer
	// 	})

	// 	t.Run("Read", func(t *testing.T) {
	// 		b := make([]byte, 5)
	// 		n, err := rb.Read(b)
	// 		assert.Equal(t, 5, n)
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, []byte("hello"), b)
	// 	})

	// 	t.Run("Close", func(t *testing.T) {
	// 		assert.NoError(t, rb.Close(), "Contract violation")
	// 		assert.NotNil(t, rb.err.Load())
	// 	})
	// })

	t.Run("Large", func(t *testing.T) {
		rb := NewRingBuffer(8)

		t.Run("Write", func(t *testing.T) {
			n, err := rb.Write([]byte("hello, world!"))
			assert.Equal(t, 13, n)
			assert.NoError(t, err)
			assert.Equal(t, "world!, ", string(rb.b)) // NOTE:  off by one in ring buffer
		})

		t.Run("Read", func(t *testing.T) {
			b := make([]byte, 5)
			n, err := rb.Read(b)
			assert.Equal(t, 5, n)
			assert.NoError(t, err)
			assert.Equal(t, "hello", string(b))
		})

		t.Run("Close", func(t *testing.T) {
			assert.NoError(t, rb.Close(), "Contract violation")
			assert.NotNil(t, rb.err.Load())
		})
	})
}

type mockBuffer struct {
	mu  sync.Mutex
	err error
	bytes.Buffer
}

func (buf *mockBuffer) Read(b []byte) (n int, err error) {
	buf.mu.Lock()
	defer buf.mu.Unlock()
	if err = buf.err; err == nil {
		n, err = buf.Buffer.Read(b)
	}
	return
}

func (buf *mockBuffer) Write(b []byte) (n int, err error) {
	buf.mu.Lock()
	defer buf.mu.Unlock()
	if err = buf.err; err == nil {
		n, err = buf.Buffer.Write(b)
	}
	return
}

func (buf *mockBuffer) Close() error {
	buf.mu.Lock()
	buf.err = io.EOF
	buf.mu.Unlock()
	return nil
}

func TestStream(t *testing.T) {
	var s string
	for i := 0; i < 1; i++ {
		s += "hello"
	}

	t.Run("BytesBuffer", func(t *testing.T) {
		src := bytes.NewBuffer([]byte(s))
		dst := new(bytes.Buffer)

		n, err := StreamBuffer(dst, src, &mockBuffer{})
		assert.Equal(t, len(s), n)
		assert.Equal(t, len(s), dst.Len())
		assert.NoError(t, err)
	})

	// t.Run("RingBuffer", func(t *testing.T) {
	// 	src := bytes.NewBuffer([]byte(s))
	// 	dst := new(bytes.Buffer)

	// 	n, err := Stream(dst, src)
	// 	assert.Equal(t, len(s), n)
	// 	assert.Equal(t, len(s), dst.Len())
	// 	assert.NoError(t, err)
	// })
}
