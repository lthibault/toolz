// Package bintoolz contains binary utilities
package bintoolz

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
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

// ReadConn reads a single []byte message from a net.Conn.  It makes use of deadlines in
// the context.  ReadConn calls SetReadDeadline.
func ReadConn(ctx context.Context, conn net.Conn, order binary.ByteOrder) ([]byte, error) {
	d, _ := ctx.Deadline()
	if err := conn.SetReadDeadline(d); err != nil {
		return nil, err
	}

	blen, err := readHeader(conn, order)
	if err != nil {
		return nil, err
	}

	if err := conn.SetReadDeadline(d); err != nil {
		return nil, err
	}

	buf := pool.Get()
	defer pool.Put(buf)

	if err = readBody(buf, conn, blen); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Write bytes to an io.Writer
func Write(w io.Writer, order binary.ByteOrder, b []byte) error {
	buf := pool.Get()
	defer pool.Put(buf)

	if err := writeHeader(buf, order, len(b)); err != nil {
		return err
	}

	buf.Write(b)
	_, err := io.Copy(w, buf)
	return err
}

// WriteConn writes a single []byte message to a net.Conn.  It makes use of deadlines in
// the context.  ReadConn calls SetWriteDeadline.
func WriteConn(ctx context.Context, conn net.Conn, order binary.ByteOrder, b []byte) error {
	d, _ := ctx.Deadline()
	if err := conn.SetWriteDeadline(d); err != nil {
		return err
	}

	buf := pool.Get()
	defer pool.Put(buf)

	if err := writeHeader(buf, order, len(b)); err != nil {
		return err
	}

	if err := conn.SetWriteDeadline(d); err != nil {
		return err
	}

	buf.Write(b)
	_, err := io.Copy(conn, buf)
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
func WriteString(w io.Writer, order binary.ByteOrder, s string) error {
	buf := pool.Get()
	defer pool.Put(buf)

	if err := writeHeader(buf, order, len(s)); err != nil {
		return err
	}

	buf.WriteString(s)

	_, err := io.Copy(w, buf)
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
