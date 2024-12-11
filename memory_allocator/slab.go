package memory_allocator

import (
	"root/stack"
	"sync"
)

type Slab struct {
	slabSize    int
	freeList    stack.Stack[[]byte]
	currentPage []byte
	pagePointer int
	sync.RWMutex
	*Allocator
}

func NewSlab(slabSize int, allocator *Allocator) Slab {
	return Slab{
		slabSize:  slabSize,
		freeList:  stack.New[[]byte](10),
		Allocator: allocator,
	}
}

func (s *Slab) AllocateMemory() ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	if !s.freeList.IsEmpty() {
		return s.freeList.Pop()
	}

	start := s.pagePointer
	end := start + s.slabSize

	if s.currentPage == nil || !IsEnoughSpace(end, len(s.currentPage)) {
		block, err := s.AllocateBlock()
		if err != nil {
			return nil, err
		}

		s.UpdatePage(block)
	}

	return s.currentPage[start:end], nil
}

func (s *Slab) UpdatePage(dataBlock []byte) {
	s.currentPage = dataBlock
	s.pagePointer = 0
}
