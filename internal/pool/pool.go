package pool

import "sync"

func NewPool[T any](new func() T) *Pool[T] {
	return &Pool[T]{pool: sync.Pool{New: func() any { return new() }}}
}

// Pool wrapper around sync.Pool
type Pool[T any] struct {
	pool sync.Pool
}

func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

func (p *Pool[T]) Put(val T) {
	p.pool.Put(val)
}
