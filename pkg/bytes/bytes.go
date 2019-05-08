package bytetoolz

import "io"

type ioBuffer struct {
	io.WriteCloser
	io.ReadCloser
}

func (buf ioBuffer) Close() error { return buf.WriteCloser.Close() }

type chanBuffer chan []byte

func (buf chanBuffer) Close() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	close(buf)
	return
}

func (buf chanBuffer) Write(b []byte) (n int, err error) {
	n = len(b)
	buf <- b
	return
}

func (buf chanBuffer) Read(b []byte) (n int, err error) {
	bb, ok := <-buf
	if !ok {
		err = io.ErrNoProgress
	}

	n = len(bb)
	copy(b, bb)
	return
}

// NewStreamBuffer returns a buffer that can be used with iotoolz.CopyStream.
func NewStreamBuffer(size int) io.ReadWriteCloser {
	if size == 0 {
		var buf ioBuffer
		buf.ReadCloser, buf.WriteCloser = io.Pipe()
		return buf
	}

	return chanBuffer(make(chan []byte, size))
}
