package iotoolz

import (
	"io"

	bytetoolz "github.com/lthibault/toolz/pkg/bytes"
	"golang.org/x/sync/errgroup"
)

func CopyStreamBuffered(dst io.Writer, src io.Reader, buf io.ReadWriteCloser) (n int, err error) {

	var g errgroup.Group

	g.Go(func() (e error) {
		defer buf.Close()

		var nn int
		b := make([]byte, 64)

		for {
			nn, e = src.Read(b)

			if _, ee := buf.Write(b[:nn]); ee != nil {
				break
			}

			if e != nil {
				break
			}
		}

		return
	})

	g.Go(func() (e error) {
		defer buf.Close()

		var nn int
		var ee error
		b := make([]byte, 64)

		for {
			nn, ee = buf.Read(b)

			if nn > 0 {
				nn, e = dst.Write(b[:nn])
				n += nn
			}

			if e != nil || ee != nil {
				break
			}
		}

		return
	})

	if err = g.Wait(); err == io.EOF {
		err = nil
	}
	return
}

func CopyStream(dst io.Writer, src io.Reader) (int, error) {
	return CopyStreamBuffered(dst, src, bytetoolz.NewStreamBuffer(0))
}
