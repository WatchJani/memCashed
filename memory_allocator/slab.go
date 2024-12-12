package memory_allocator

import (
	"root/stack"
	"sync"
)

type SlabManager struct {
	slabs []Slab
}

func (s *SlabManager) GetSlabIndex(index int) *Slab {
	return &s.slabs[index]
}

func NewSlabManager(slabs []Slab) SlabManager {
	return SlabManager{
		slabs: slabs,
	}
}

func (s *SlabManager) ChoseSlab(dataSize int) *Slab {
	low, high := 0, len(s.slabs)-1
	result := -1

	slabs := s.slabs

	for low <= high {
		mid := low + (high-low)/2
		if slabs[mid].slabSize > dataSize {
			result = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return &slabs[result]
}

type Slab struct {
	slabSize    int
	freeList    stack.Stack[[]byte]
	currentPage []byte
	pagePointer int
	sync.RWMutex
	*Allocator
}

func (s *Slab) GetCurrentPage() []byte {
	return s.currentPage
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
