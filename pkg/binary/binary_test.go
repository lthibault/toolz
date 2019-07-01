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
	t.Run("ReadPath", func(t *testing.T) {
		b := bytes.NewBuffer(msgLenBytes)
		b.WriteString(testMsg)

		msg, err := ReadString(b, binary.BigEndian)
		assert.NoError(t, err)
		assert.Equal(t, testMsg, msg)
	})

	t.Run("WritePath", func(t *testing.T) {
		b := new(bytes.Buffer)
		assert.NoError(t, WriteString(b, binary.BigEndian, testMsg))
		assert.Equal(t, msgLenBytes, b.Bytes()[:2])
		assert.Equal(t, testMsg, b.String()[2:])
	})
}
