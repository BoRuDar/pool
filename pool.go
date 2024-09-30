package pool

import (
	"sync"
)

// Pool is a memory pool for T, where T is a struct.
type Pool[T any] struct {
	p        sync.Pool
	mu       *sync.Mutex
	defaultT *T
}

// New creates new Pool.
func New[T any]() Pool[T] {
	return Pool[T]{
		p: sync.Pool{New: func() any {
			return new(T)
		}},
		mu:       &sync.Mutex{},
		defaultT: new(T),
	}
}

// Get return an object from the pool or creates a new one if no more objects for reusing.
func (p *Pool[T]) Get() *T {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.p.Get().(*T)
}

// Return objects to the pool and resets all struct fields to their default values.
func (p *Pool[T]) Return(v *T) {
	p.mu.Lock()
	defer p.mu.Unlock()

	*v = *p.defaultT
	p.p.Put(v)
}
