package iotoolz

import (
	"io"
	"runtime"
	"sync/atomic"

	turbine "github.com/lthibault/turbine/pkg"

	"golang.org/x/sync/errgroup"
)

// RingBuffer is a fixed-sized buffer of bytes with Read and Write methods.  It is for
// concurrent reads and writes provided there is at most one concurrent call to Read and
// one concurrent call to Write.
//
// Contrary to bytes.Buffer, data is streamed continuously through the buffer.
type RingBuffer struct {
	mask int64
	b    []byte
	t    *turbine.Turbine

	pending      int32
	lower, upper int64
	err          atomic.Value
}

// NewRingBuffer with specified capacity
func NewRingBuffer(size int64) (rb *RingBuffer) {
	if size > 0 && (size&(size-1)) != 0 {
		panic("size must be power of 2")
	}

	rb = new(RingBuffer)
	defer rb.Reset()

	rb.mask = size - 1
	rb.t = turbine.New(size, rb)

	return rb
}

// Close the buffer.  Always returns a nil error.
func (r *RingBuffer) Close() error {
	r.t.Stop()
	if r.err.Load() == nil {
		r.setErr(io.EOF)
	}
	return nil
}

// Reset clears all internal state.  It is not thread-safe.
func (r *RingBuffer) Reset() {
	r.Close()
	r.b = make([]byte, r.mask+1)
	r.err = atomic.Value{}
	r.t.Start()
}

func (r *RingBuffer) Write(b []byte) (n int, err error) {
	var ok bool
	if err, ok = r.err.Load().(error); ok && err != nil {
		return
	}

	n = len(b)
	for _, x := range b {
		seq := r.t.Writer().Reserve(1)
		r.b[seq&r.mask] = x
		r.t.Writer().Commit(seq)
	}

	return
}

func (r *RingBuffer) Read(b []byte) (n int, err error) {
	var ok bool
	for atomic.LoadInt32(&r.pending) != 1 { // either locked or non-pending
		if err, ok = r.err.Load().(error); ok && err != nil {
			return
		}

		runtime.Gosched() // no pending reads
	}

	lenb := int64(len(b))
	for i := r.lower; i <= r.upper; i++ {
		if i-r.lower >= lenb {
			r.lower = i
			return
		}

		b[i-r.lower] = r.b[i&r.mask]
		n++
	}

	atomic.StoreInt32(&r.pending, 0)
	return
}

func (r *RingBuffer) setErr(err error) { r.err.Store(err) }

// Consume values from the circular buffer
func (r *RingBuffer) Consume(lower, upper int64) {
	for !atomic.CompareAndSwapInt32(&r.pending, 0, -1) { // set to locked
		if r.err.Load() != nil {
			return
		}

		runtime.Gosched() // if we get here, there are pending reads
	}

	r.lower = lower
	r.upper = upper

	if !atomic.CompareAndSwapInt32(&r.pending, -1, 1) { // set to pending
		panic("unreachable")
	}
}

// Stream works like io.Copy, except that it streams bytes continuously rather than
// waiting for the writer to close.
func Stream(dst io.Writer, src io.Reader) (written int, err error) {
	return StreamBuffer(dst, src, NewRingBuffer(64))
}

// StreamBuffer works like io.CopyBuffer, except that it streams bytes continuously rather than
// waiting for the writer to close.
func StreamBuffer(dst io.Writer, src io.Reader, buf io.ReadWriteCloser) (written int, err error) {

	var g errgroup.Group
	g.Go(func() error {
		defer buf.Close()

		var eof bool
		b := make([]byte, 64)

		for {
			switch _, e := src.Read(b); e {
			case nil:
			case io.EOF:
				eof = true
			default:
				return e
			}

			n, e := buf.Write(b)
			written += n

			if e != nil || eof {
				return e
			}
		}
	})

	g.Go(func() error {
		var eof bool
		b := make([]byte, 64)

		for {
			n, e := buf.Read(b)

			switch e {
			case nil:
			case io.EOF:
				eof = true
			default:
				return e
			}

			if _, e := dst.Write(b[:n]); e != nil || eof {
				return e
			}
		}
	})

	err = g.Wait()
	return
}
