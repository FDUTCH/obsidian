package util

import "github.com/FDUTCH/obsidian/internal/pool"

func NewBufferPool(size int) *pool.Pool[[]byte] {
	return pool.NewPool[[]byte](func() []byte {
		return make([]byte, size)
	})
}
