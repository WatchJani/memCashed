package memory_allocator

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"unsafe"

	"github.com/WatchJani/memCashed/client"
	"github.com/WatchJani/memCashed/cmd/constants"
	"github.com/WatchJani/memCashed/cmd/link_list"
	"github.com/WatchJani/memCashed/cmd/stack"
)

var (
	// ErrOperationIsNotSupported is the error returned when an unsupported operation is attempted.
	ErrOperationIsNotSupported = errors.New("operation is not supported")

	// NoReq is a buffer used for reading unprocessed data when memory allocation fails.
	HeaderSize        = 10
	NoReq      []byte = make([]byte, MiB+HeaderSize) // Buffer to read unprocessed data

)

// SlabManager manages slabs, LRU (Least Recently Used) caches, and memory allocation.
type SlabManager struct {
	slabs        []Slab          // Slabs for memory allocation
	lru          []link_list.DLL // Least Recently Used (LRU) cache for each slab
	sync.RWMutex                 // Mutex to protect concurrent access to shared data
	store        sync.Map        // Store to hold key-value pairs for cache management
	JobCh        chan Transfer   // Channel to receive transfer jobs for processing
}

// Transfer represents a data payload and connection information for a transfer task.
type Transfer struct {
	payload []byte   // Data payload
	conn    net.Conn // Network connection
	index   int      // Index of the slab category
}

// Key represents a stored object with its field, TTL (Time-To-Live), and a pointer to its node in the LRU list.
type Key struct {
	field   []byte          // Object data field
	ttl     time.Time       // Time-To-Live for the object
	pointer *link_list.Node // Pointer to the node in the LRU list
}

// NewTransfer creates a new Transfer object with the specified payload, index, and connection.
func NewTransfer(payload []byte, index int, conn net.Conn) Transfer {
	return Transfer{
		payload: payload,
		conn:    conn,
		index:   index,
	}
}

// FreeSpace frees space in the slab's LRU cache by removing the least recently used node.
func (s *SlabManager) FreeSpace(index, slabSize int) ([]byte, string) {
	s.Lock()
	defer s.Unlock()

	lastNode := s.lru[index].LastNode() // Get the last (least recently used) node

	s.lru[index].Delete(lastNode) // Delete the last node in the LRU cache

	// Get free space from LRU after deleting the node
	return s.lru[index].GetLRUFreeSpace(lastNode, slabSize), lastNode.GetKey()
}

// GetSlabIndex returns the slab at the specified index.
func (s *SlabManager) GetSlabIndex(index int) *Slab {
	return &s.slabs[index]
}

// GetLRUIndex returns the LRU cache at the specified index.
func (s *SlabManager) GetLRUIndex(index int) *link_list.DLL {
	return &s.lru[index]
}

// NewSlabManager creates a new SlabManager with the provided slabs and starts worker goroutines.
func NewSlabManager(slabs []Slab, numberOfWorker int) *SlabManager {
	sm := &SlabManager{
		slabs: slabs,
		lru:   make([]link_list.DLL, len(slabs)), // Initialize LRU for each slab
		JobCh: make(chan Transfer),               // Channel for receiving transfer jobs
	}

	// Start a worker goroutine of numberOfWorker
	for range numberOfWorker {
		go sm.Worker()
	}

	return sm
}

// GetSlab allocates a slab of memory based on the payload size, handles errors, and frees space if necessary.
func (s *SlabManager) GetSlab(payloadSize int, conn net.Conn) ([]byte, int, error) {
	slabIndex, chunkSize := s.GetIndex(payloadSize)

	// Attempt to allocate memory from the chosen slab
	slabBlock, err := s.ChoseSlab(slabIndex).AllocateMemory()
	if err != nil {
		// If memory allocation fails, handle the error
		fmt.Println(err)

		// If slab is inactive, notify the client and try to read the rest of the request
		if !s.GetSlabIndex(slabIndex).IsSlabActive() {
			conn.Write([]byte(err.Error()))

			// and if I can't allocate memory to my server I still have to read the req to the end
			_, err := conn.Read(NoReq)
			if err != nil {
				return nil, -1, err
			}
		}

		// If there is no more space in memory, uses LRU
		// (Least Recently Used) policy to free up space.
		s.Lock()
		lastNode := s.lru[slabIndex].LastNode()                           // Get the last LRU node
		s.lru[slabIndex].Delete(lastNode)                                 // Delete last node in
		slabBlock = s.lru[slabIndex].GetLRUFreeSpace(lastNode, chunkSize) // Get free space after deleting the node
		s.Unlock()

		// Deletes the key from the hash table.
		s.store.Delete(lastNode.GetKey())
	}

	return slabBlock, slabIndex, nil
}

// TLLParser converts a TTL value into a time.Time object.
func TLLParser(ttl uint32) time.Time {
	if ttl > 0 {
		return time.Now().Add(time.Duration(ttl) * time.Second) // Add TTL to the current time
	}

	return time.Time{} // Return an empty time if TTL is 0
}

// Worker listens for transfer jobs and processes them based on the payload command.
func (s *SlabManager) Worker() {
	for payload := range s.JobCh {
		switch payload.payload[0] {
		case constants.SetOperation: // Command to store data
			_, keySize, ttl, bodySize := client.Decode(payload.payload) // Decode the payload

			bodyOffset := constants.HeaderSize + keySize
			key := string(payload.payload[constants.HeaderSize:bodyOffset]) // Extract key from the payload

			// Insert the key into the LRU cache
			node := s.lru[payload.index].Inset(link_list.NewValue(unsafe.Pointer(&payload.payload[0]), key))

			// Store the key-value pair in the store with TTL
			s.store.Store(key, Key{
				field:   payload.payload[bodyOffset : bodyOffset+bodySize],
				ttl:     TLLParser(ttl),
				pointer: node,
			})

			if _, err := payload.conn.Write(constants.ObjectInserted); err != nil {
				log.Println(err) // Log any errors that occur while writing to the connection
			}
		case constants.GetOperation: // Command to get data
			_, keySize, _, _ := client.Decode(payload.payload)                                  // Decode the payload
			key := string(payload.payload[constants.HeaderSize : constants.HeaderSize+keySize]) // Extract key from the payload

			s.slabs[payload.index].freeList.Push(unsafe.Pointer(&payload.payload[0])) //delete our header space

			// Fetch the value from the store
			valueObject, isFound := s.store.Load(key)
			if !isFound {
				if _, err := payload.conn.Write(constants.ErrObjectNotFound); err != nil {
					log.Println(err)
				}
				continue
			}

			value := valueObject.(Key)

			// Check if the TTL has expired and delete the object if expired
			if !value.ttl.IsZero() && time.Now().After(value.ttl) {
				s.store.Delete(key)
				s.lru[payload.index].Delete(value.pointer) // Remove the node from LRU
				memoryPointer := value.pointer.GetPointer()
				s.slabs[payload.index].freeList.Push(memoryPointer)

				if _, err := payload.conn.Write(constants.ErrTimeExpire); err != nil {
					log.Println(err)
				}
				continue
			}
			s.lru[payload.index].Read(value.pointer)

			// Return the field data if found
			if _, err := payload.conn.Write(value.field); err != nil {
				log.Println(err)
			}
		case constants.DeleteOperation: // Command to delete data
			_, keySize, _, _ := client.Decode(payload.payload) // Decode the payload
			key := string(payload.payload[10 : 10+keySize])    // Extract key from the payload

			s.slabs[payload.index].freeList.Push(unsafe.Pointer(&payload.payload[0])) //delete our header space

			// Fetch and delete the object from the store
			valueObject, isFound := s.store.Load(key)
			if !isFound {
				if _, err := payload.conn.Write(constants.ErrObjectNotFound); err != nil {
					log.Println(err)
				}
				continue
			}

			value := valueObject.(Key)
			s.store.Delete(key)

			memoryPointer := value.pointer.GetPointer()

			s.lru[payload.index].Delete(value.pointer) // Remove from LRU
			s.slabs[payload.index].freeList.Push(memoryPointer)
			if _, err := payload.conn.Write(constants.ObjectDeleted); err != nil {
				log.Println(err)
			}
		default:
			log.Println(ErrOperationIsNotSupported)
		}
	}
}

// GetIndex performs a binary search to find the appropriate slab index based on the data size.
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

// ChoseSlab returns the slab at the specified index.
func (s *SlabManager) ChoseSlab(index int) *Slab {
	return &s.slabs[index]
}

// Slab represents a memory slab used for allocation.
type Slab struct {
	slabSize     int                         // Size of the slab
	freeList     stack.Stack[unsafe.Pointer] // Stack of free blocks in the slab
	currentPage  []byte                      // Current memory page in the slab
	pagePointer  int                         // Pointer to the current position in the slab
	sync.RWMutex                             // Mutex to protect access to the slab
	*Allocator                               // Memory allocator associated with the slab
}

// IsSlabActive checks if the slab has an active memory page.
func (s *Slab) IsSlabActive() bool {
	return s.currentPage != nil
}

// GetCurrentPage returns the current page of the slab.
func (s *Slab) GetCurrentPage() []byte {
	return s.currentPage
}

// NewSlab creates a new Slab with the specified size and allocator.
func NewSlab(slabSize, maxMemoryAllocate int, allocator *Allocator) Slab {
	return Slab{
		slabSize:  slabSize,
		freeList:  stack.New[unsafe.Pointer](10),
		Allocator: allocator,
	}
}

// AllocateMemory allocates memory for the slab, either by reusing a free block or allocating a new page.
func (s *Slab) AllocateMemory() ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	// Try to pop from the free list if there are free blocks
	if !s.freeList.IsEmpty() {
		ptr, err := s.freeList.Pop()
		return unsafe.Slice((*byte)(ptr), s.slabSize), err
	}

	// Calculate the memory range for the new allocation
	start := s.pagePointer
	end := start + s.slabSize

	// If no active page or insufficient space, allocate a new page
	if s.currentPage == nil || !IsEnoughSpace(end, len(s.currentPage)) {
		block, err := s.AllocateBlock()
		if err != nil {
			return nil, err
		}

		// Update the current page with the new block
		s.UpdatePage(block)
		return s.currentPage[0:s.slabSize], nil //new memory block
	}

	// Return the allocated memory block from the current page
	return s.currentPage[start:end], nil
}

func (s *Slab) UpdatePage(dataBlock []byte) {
	s.currentPage = dataBlock
	s.pagePointer = 0
}
