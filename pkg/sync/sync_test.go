package synctoolz

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVar(t *testing.T) {
	t.Parallel()

	t.Run("ReadConcurrent", func(t *testing.T) {
		var v Var

		var wg sync.WaitGroup
		wg.Add(4)

		go t.Run("GetA", func(t *testing.T) {
			defer wg.Done()
			assert.Equal(t, 1, v.Get())
		})

		go t.Run("GetB", func(t *testing.T) {
			defer wg.Done()
			assert.Equal(t, 1, v.Get())
		})

		go t.Run("GetC", func(t *testing.T) {
			defer wg.Done()
			assert.Equal(t, 1, v.Get())
		})

		go t.Run("GetD", func(t *testing.T) {
			defer wg.Done()
			assert.Equal(t, 1, v.Get())
		})

		v.Set(1)
		wg.Wait()
	})

	t.Run("SetOnce", func(t *testing.T) {
		var v Var
		v.Set(1)
		assert.Panics(t, func() { v.Set(9001) })
	})
}
