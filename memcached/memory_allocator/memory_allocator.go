package memory_allocator

import (
	"sync"

	"github.com/WatchJani/memCashed/tree/master/memcached/constants"
)

// Allocator is a memory allocator that manages a slice of bytes and keeps track of the next available index.
type Allocator struct {
	memory []byte
	next   int
	sync.RWMutex
}

// GetNext returns the next available index in the allocator's memory.
func (a *Allocator) GetNext() int {
	return a.next
}

// New creates a new Allocator with the specified capacity.
func New(capacity int) *Allocator {
	return &Allocator{
		memory: make([]byte, capacity),
	}
}

// IsEnoughSpace checks if there is enough space to allocate a block of memory from 'end' position to the given length.
func IsEnoughSpace(end, len int) bool {
	return end <= len
}

// AllocateBlock allocates a block of memory in the allocator.
// It locks the allocator for thread-safety and returns a slice of bytes or an error if there isn't enough space.
func (a *Allocator) AllocateBlock() ([]byte, error) {
	a.Lock()
	defer a.Unlock()

	start := a.next
	end := start + constants.MiB

	if !IsEnoughSpace(end, len(a.memory)) {
		return nil, constants.ErrNotEnoughSpace
	}

	a.next = end
	return a.memory[start:end], nil
}
