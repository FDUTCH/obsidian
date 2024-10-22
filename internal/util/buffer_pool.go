package util

import "github.com/FDUTCH/obsidian/internal/pool"

// NewBufferPool returns new *pool.Pool[[]byte]
func NewBufferPool(size int) *pool.Pool[[]byte] {
	return pool.NewPool[[]byte](func() []byte {
		return make([]byte, size)
	})
}
