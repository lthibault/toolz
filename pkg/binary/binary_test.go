package bintoolz

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testMsg = "foo bar foobar"

var (
	msgLenBytes = make([]byte, 2)
)

func init() {
	binary.BigEndian.PutUint16(msgLenBytes, uint16(len(testMsg)))
}

func TestReadWrite(t *testing.T) {
	t.Run("Read", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			b := bytes.NewBuffer(msgLenBytes)
			b.WriteString(testMsg)

			msg, err := ReadString(b, binary.BigEndian)
			assert.NoError(t, err)
			assert.Equal(t, testMsg, msg)
		})

		t.Run("Bytes", func(t *testing.T) {
			b := bytes.NewBuffer(msgLenBytes)
			b.WriteString(testMsg)

			msg, err := Read(b, binary.BigEndian)
			assert.NoError(t, err)
			assert.Equal(t, testMsg, string(msg))
		})
	})

	t.Run("Write", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			b := new(bytes.Buffer)
			assert.NoError(t, WriteString(b, binary.BigEndian, testMsg))
			assert.Equal(t, msgLenBytes, b.Bytes()[:2])
			assert.Equal(t, testMsg, b.String()[2:])
		})

		t.Run("Bytes", func(t *testing.T) {
			b := new(bytes.Buffer)
			assert.NoError(t, Write(b, binary.BigEndian, []byte(testMsg)))
			assert.Equal(t, msgLenBytes, b.Bytes()[:2])
			assert.Equal(t, testMsg, b.String()[2:])
		})
	})

	t.Run("Integration", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			b := new(bytes.Buffer)
			assert.NoError(t, WriteString(b, binary.BigEndian, testMsg))
			msg, err := ReadString(b, binary.BigEndian)
			assert.NoError(t, err)
			assert.Equal(t, testMsg, msg)
		})

		t.Run("Bytes", func(t *testing.T) {
			b := new(bytes.Buffer)
			assert.NoError(t, Write(b, binary.BigEndian, []byte(testMsg)))
			msg, err := Read(b, binary.BigEndian)
			assert.NoError(t, err)
			assert.Equal(t, []byte(testMsg), msg)
		})
	})
}
