package synctoolz

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// FuncGroup is a WaitGroup that acts on functions
type FuncGroup sync.WaitGroup

// Go runs the specified function in a goroutine
func (g *FuncGroup) Go(f func()) {
	(*sync.WaitGroup)(unsafe.Pointer(g)).Add(1)
	defer (*sync.WaitGroup)(unsafe.Pointer(g)).Done()
	f()
}

// Wait for all goroutines to complete
func (g *FuncGroup) Wait() { (*sync.WaitGroup)(unsafe.Pointer(g)).Wait() }

// Ctr is a lock-free counter
type Ctr uint32

// Incr increments the counter
func (ctr *Ctr) Incr() uint32 { return atomic.AddUint32((*uint32)(unsafe.Pointer(ctr)), 1) }

// Decr decrements the counter
func (ctr *Ctr) Decr() uint32 { return atomic.AddUint32((*uint32)(unsafe.Pointer(ctr)), ^uint32(0)) }

// Num returns the generic-integer value of the counter.  This is useful for integer comparisons.
func (ctr *Ctr) Num() int { return int(atomic.LoadUint32((*uint32)(unsafe.Pointer(ctr)))) }

// Flag is a lock-free boolean flag
type Flag uint32

// Set the flag's value to true
func (f *Flag) Set() { atomic.CompareAndSwapUint32((*uint32)(unsafe.Pointer(f)), 0, 1) }

// Unset the flag, making its value false
func (f *Flag) Unset() { atomic.CompareAndSwapUint32((*uint32)(unsafe.Pointer(f)), 1, 0) }

// Bool evaluates the flag's value
func (f *Flag) Bool() bool { return atomic.LoadUint32((*uint32)(unsafe.Pointer(f))) != 0 }

// Var is a variable that can be set exactly once.  Calls to Get will block until it is
// set.
type Var struct {
	init, setter sync.Once
	v            interface{}
	ch           chan interface{}
	rq           chan struct{}
}

func (v *Var) setup() {
	v.ch = make(chan interface{}, 1)
	v.rq = make(chan struct{})
}

// Ready returns a channel that is closed when the Var's value is set
func (v *Var) Ready() <-chan struct{} { return v.rq }

// Get the value once it is set
func (v *Var) Get() interface{} {
	v.init.Do(v.setup)
	v.setter.Do(func() {
		v.v = <-v.ch
	})

	return v.v
}

// Set the value
func (v *Var) Set(value interface{}) {
	v.init.Do(v.setup)

	v.ch <- value
	close(v.ch)
	close(v.rq)
}
