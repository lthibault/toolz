package iotoolz

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
)

// CopyStreamBuffered behaves like CopyStream, but allows the user to specify the buffer
// to be used.
func CopyStreamBuffered(dst io.Writer, src io.Reader, buf chan []byte) (int, error) {
	defer close(buf)

	var n int
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() (err error) {
		var nn int
		b := make([]byte, 64)

		for {
			nn, err = src.Read(b)

			if nn > 0 {
				select {
				case buf <- b[:nn]:
				case <-ctx.Done():
					err = ctx.Err()
				}
			}

			if err != nil {
				return
			}
		}
	})

	g.Go(func() (err error) {

		var nn int
		for {
			select {
			case b := <-buf:
				nn, err = dst.Write(b)
				n += nn
			case <-ctx.Done():
				err = ctx.Err()
			}

			if err != nil {
				return
			}
		}
	})

	switch err := g.Wait(); err {
	case nil:
		return n, nil
	case io.EOF:
		return n, nil
	default:
		return n, err
	}
}

// CopyStream behaves like io.Copy except that it passes data continuously from src to
// dst instead of waiting for the former to return io.EOF.
func CopyStream(dst io.Writer, src io.Reader) (int, error) {
	return CopyStreamBuffered(dst, src, make(chan []byte))
}
