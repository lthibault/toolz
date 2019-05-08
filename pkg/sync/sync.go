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

type Ctr uint32

func (ctr *Ctr) Incr() uint32 { return atomic.AddUint32((*uint32)(unsafe.Pointer(ctr)), 1) }
func (ctr *Ctr) Decr() uint32 { return atomic.AddUint32((*uint32)(unsafe.Pointer(ctr)), ^uint32(0)) }
func (ctr *Ctr) Num() int     { return int(atomic.LoadUint32((*uint32)(unsafe.Pointer(ctr)))) }

type Flag uint32

func (f *Flag) Set()       { atomic.CompareAndSwapUint32((*uint32)(unsafe.Pointer(f)), 0, 1) }
func (f *Flag) Bool() bool { return atomic.LoadUint32((*uint32)(unsafe.Pointer(f))) != 0 }
