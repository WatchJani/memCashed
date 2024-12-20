package memory_allocator

import (
	"errors"
	"fmt"
	"log"
	"net"
	"root/client"
	"root/cmd/link_list"
	"root/cmd/stack"
	"sync"
	"time"
	"unsafe"
)

var (
	ErrOperationIsNotSupported = errors.New("operation is not supported")
)

type SlabManager struct {
	slabs []Slab
	lru   []link_list.DLL
	sync.RWMutex
	store   sync.Map
	JobCh   chan Transfer
	counter int
}

type Transfer struct {
	payload []byte
	conn    net.Conn
}

type Key struct {
	field   []byte
	ttl     time.Time
	pointer *link_list.Node
}

func NewTransfer(payload []byte, conn net.Conn) Transfer {
	return Transfer{
		payload: payload,
		conn:    conn,
	}
}

func (s *SlabManager) FreeSpace(index, slabSize int) ([]byte, string) {
	s.Lock()
	defer s.Unlock()

	lastNode := s.lru[index].LastNode()

	s.lru[index].Delete(lastNode) //Delete last node in
	// s.lru[index].Read(lastNode) //set node to root

	return s.lru[index].GetLRUFreeSpace(lastNode, slabSize), lastNode.GetKey()
}

func (s *SlabManager) GetSlabIndex(index int) *Slab {
	return &s.slabs[index]
}

func (s *SlabManager) GetLRUIndex(index int) *link_list.DLL {
	return &s.lru[index]
}

func NewSlabManager(slabs []Slab) *SlabManager {
	sm := &SlabManager{
		slabs: slabs,
		lru:   make([]link_list.DLL, len(slabs)),
		JobCh: make(chan Transfer),
	}

	for range slabs {
		go sm.Worker()
	}

	return sm
}

var NoReq []byte = make([]byte, 1024*1024+10)

func (s *SlabManager) GetSlab(payloadSize int, conn net.Conn) ([]byte, error) {
	slabIndex, chunkSize := s.GetIndex(payloadSize)

	slabBlock, err := s.ChoseSlab(slabIndex).AllocateMemory()
	if err != nil {
		// If an error occurs during memory allocation,
		// prints the error.
		fmt.Println(err)

		// If the slab is not active, sends the error to the
		// client and tries to read the rest of the request.
		if !s.GetSlabIndex(slabIndex).IsSlabActive() {
			conn.Write([]byte(err.Error()))

			// and if I can't allocate memory to my server I still have to read the req to the end
			_, err := conn.Read(NoReq)
			if err != nil {
				return nil, err
			}
		}

		// If there is no more space in memory, uses LRU
		// (Least Recently Used) policy to free up space.

		s.Lock()
		lastNode := s.lru[slabIndex].LastNode()

		s.lru[slabIndex].Delete(lastNode) //Delete last node in
		// s.lru[index].Read(lastNode) //set node to root

		slabBlock = s.lru[slabIndex].GetLRUFreeSpace(lastNode, chunkSize)
		s.Unlock()
		// key := slabBlock[:keyLength]

		// Deletes the key from the hash table.
		s.store.Delete(lastNode.GetKey())
	}

	return slabBlock, nil
}

func TLLParser(ttl uint32) time.Time {
	if ttl > 0 {
		return time.Now().Add(time.Duration(ttl) * time.Second)
	}

	return time.Time{}
}

func (s *SlabManager) GetNumberOfReq() int {
	return s.counter
}

func (s *SlabManager) Worker() {
	for payload := range s.JobCh {
		// operation, keySize, ttl, bodySize := client.Decode(payload.payload)
		index, _ := s.GetIndex(len(payload.payload))
		switch payload.payload[0] {
		case 'S':
			_, keySize, ttl, bodySize := client.Decode(payload.payload)

			key := string(payload.payload[10 : 10+keySize+keySize])

			s.store.Store(key, Key{
				field: payload.payload[10 : 10+bodySize],
				ttl:   TLLParser(ttl),
			})

			//insert in lru

			s.lru[index].Inset(link_list.NewValue(unsafe.Pointer(&payload.payload[0]), key))

			if _, err := payload.conn.Write([]byte("object inserted")); err != nil {
				log.Println(err)
			}

			s.Lock()
			s.counter++
			s.Unlock()
		case 'G':
			_, keySize, _, _ := client.Decode(payload.payload) //another parser i need
			key := string(payload.payload[10 : 10+keySize+keySize])

			valueObject, isFound := s.store.Load(key)
			if !isFound {
				if _, err := payload.conn.Write([]byte("object not found")); err != nil {
					log.Println(err)
				}
				continue
			}

			value := valueObject.(Key)
			if !value.ttl.IsZero() && time.Now().After(value.ttl) {
				s.store.Delete(key)
				s.lru[index].Delete(value.pointer) //delete from lru
				if _, err := payload.conn.Write([]byte("time expire")); err != nil {
					log.Println(err)
				}
				continue
			}

			if _, err := payload.conn.Write(value.field); err != nil {
				log.Println(err)
			}
		case 'D':
			_, keySize, _, _ := client.Decode(payload.payload) //another parser i need
			key := string(payload.payload[10 : 10+keySize+keySize])

			valueObject, isFound := s.store.Load(key)
			if !isFound {
				if _, err := payload.conn.Write([]byte("object not found")); err != nil {
					log.Println(err)
				}
				continue
			}

			value := valueObject.(Key)
			s.store.Delete(key)

			s.lru[index].Delete(value.pointer) //delete from lru
			// obj.lru. //add to stack
			if _, err := payload.conn.Write([]byte("Deleted")); err != nil {
				log.Println(err)
			}
		default:
			log.Println(ErrOperationIsNotSupported)
		}
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

// maxMemoryAllocate not define yet
func NewSlab(slabSize, maxMemoryAllocate int, allocator *Allocator) Slab {
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
