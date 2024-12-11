package memory_allocator

import (
	"errors"
	"sync"
)

const MB = 1024 * 1024

var (
	ErrNotEnoughSpace = errors.New("there is not enough space")
)

type Allocator struct {
	memory []byte
	next   int
	sync.RWMutex
}

func (a *Allocator) GetNext() int {
	return a.next
}

func New(capacity int) *Allocator {
	return &Allocator{
		memory: make([]byte, capacity),
	}
}

func IsEnoughSpace(end, len int) bool {
	return end <= len
}

func (a *Allocator) AllocateBlock() ([]byte, error) {
	a.Lock()
	defer a.Unlock()

	start := a.next
	end := start + MB

	if !IsEnoughSpace(end, len(a.memory)) {
		return nil, ErrNotEnoughSpace
	}

	a.next = end
	return a.memory[start:end], nil
}
