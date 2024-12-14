package memory_allocator

import (
	"root/link_list"
	"root/stack"
	"sync"
)

type SlabManager struct {
	slabs []Slab
	lru   []link_list.DLL
	sync.RWMutex
}

func (s *SlabManager) FreeSpace(index, slabSize int) []byte {
	s.Lock()
	defer s.Unlock()

	lastNode := s.lru[index].LastNode()
	s.lru[index].Delete(lastNode) //Delete last node in
	// s.lru[index].Read(lastNode) //set node to root

	return s.lru[index].GetLRUFreeSpace(lastNode, slabSize)
}

func (s *SlabManager) GetSlabIndex(index int) *Slab {
	return &s.slabs[index]
}

func (s *SlabManager) GetLRUIndex(index int) *link_list.DLL {
	return &s.lru[index]
}

func NewSlabManager(slabs []Slab) SlabManager {
	return SlabManager{
		slabs: slabs,
		lru:   make([]link_list.DLL, len(slabs)),
	}
}

func (s *SlabManager) GetIndex(dataSize int) (int, int) {
	low, high := 0, len(s.slabs)-1
	result := high

	slabs := s.slabs

	for low <= high {
		mid := low + (high-low)/2
		if slabs[mid].slabSize >= dataSize {
			result = mid
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	return result, slabs[result].slabSize
}

func (s *SlabManager) ChoseSlab(index int) *Slab {
	return &s.slabs[index]
}

type Slab struct {
	slabSize    int
	freeList    stack.Stack[[]byte]
	currentPage []byte
	pagePointer int
	sync.RWMutex
	*Allocator
}

func (s *Slab) IsSlabActive() bool {
	return s.currentPage != nil
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
		return s.currentPage[0:s.slabSize], nil //new memory block
	}

	return s.currentPage[start:end], nil
}

func (s *Slab) UpdatePage(dataBlock []byte) {
	s.currentPage = dataBlock
	s.pagePointer = 0
}
