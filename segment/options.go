package segment

import "time"

type PoolOptions func(pool *Pool)

// WithCapacity sets the future cache capacity.
//
// By default, the capacity is 10*1024*1024.
func WithCapacity(capacity int) PoolOptions {
	return func(p *Pool) {
		p.capacity = capacity
	}
}

// WithInitialCapacity sets the minimum total size for the internal data structures.
// Providing a large enough estimate at construction time avoids the need for expensive resizing operations later,
// but setting this value unnecessarily high wastes memory.
//
// By default, the initial capacity is 1000.
func WithInitialCapacity(initialCapacity int) PoolOptions {
	return func(p *Pool) {
		p.initialCapacity = initialCapacity
	}
}

// WithVariableTTL specifies that each item should be automatically removed from the cache once a duration has elapsed
// after the item's creation. Items are expired based on the custom ttl specified for each item separately.
//
// By default, the ttl is 300 * time.Second.
func WithVariableTTL(variableTTL time.Duration) PoolOptions {
	return func(p *Pool) {
		p.variableTTL = variableTTL
	}
}
