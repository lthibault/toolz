// Package bintoolz contains binary utilities
package bintoolz

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"sync"
	"unsafe"
)

var pool = bpool{sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}}

type bpool struct{ sync.Pool }

func (bp *bpool) Get() *bytes.Buffer { return bp.Pool.Get().(*bytes.Buffer) }
func (bp *bpool) Put(b *bytes.Buffer) {
	b.Reset()
	bp.Pool.Put(b)
}

// Read bytes from an io.Reader
func Read(r io.Reader, order binary.ByteOrder) ([]byte, error) {
	blen, err := readHeader(r, order)
	if err != nil {
		return nil, err
	}

	buf := pool.Get()
	defer pool.Put(buf)

	if err = readBody(buf, r, blen); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Write bytes to an io.Writer
func Write(w io.Writer, order binary.ByteOrder, b []byte) error {
	if err := writeHeader(w, order, len(b)); err != nil {
		return err
	}

	_, err := io.Copy(w, strings.NewReader(*(*string)(unsafe.Pointer(&b))))
	return err
}

// ReadString from an io.Reader
func ReadString(r io.Reader, order binary.ByteOrder) (string, error) {
	blen, err := readHeader(r, order)
	if err != nil {
		return "", err
	}

	buf := pool.Get()
	defer pool.Put(buf)

	if err = readBody(buf, r, blen); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// WriteString to an io.Writer
func WriteString(w io.Writer, order binary.ByteOrder, b string) error {
	if err := writeHeader(w, order, len(b)); err != nil {
		return err
	}

	_, err := io.Copy(w, strings.NewReader(b))
	return err
}

func readBody(buf *bytes.Buffer, r io.Reader, blen uint16) error {
	_, err := io.Copy(buf, io.LimitReader(r, int64(blen)))
	return err
}

func readHeader(r io.Reader, order binary.ByteOrder) (blen uint16, err error) {
	err = binary.Read(r, order, &blen)
	return
}

func writeHeader(w io.Writer, order binary.ByteOrder, n int) error {
	return binary.Write(w, order, uint16(n))
}
