package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"root/hash_table"
	"root/internal/types"
	"root/memory_allocator"
	"sync"
)

const (
	HeaderSize = 10
)

// The `Server` structure defines the core server
// parameters, including the address, maximum number
// of connections, active connections, memory
// allocator, and hash table engine.
type Server struct {
	Addr       string
	MaxConn    int
	ActiveConn int
	Slab       memory_allocator.SlabManager
	*hash_table.Engine
	sync.RWMutex
}

// NewServer creating new server instance with
// load configuration
func NewServer() *Server {
	config := types.LoadConfiguration()
	newAllocator := config.MemoryAllocator()

	return &Server{
		Addr:    config.Port(),
		MaxConn: config.MaxConnection(),
		Slab: memory_allocator.NewSlabManager(
			config.Slabs(newAllocator),
		),
		Engine: config.Workers(),
	}
}

// The `Run` function starts the server and
// listens for incoming TCP connections.
// The main server loop continues to run until
// an error occurs or the connection is closed.
func (s *Server) Run() error {
	ls, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Println(err)
	}

	defer func() {
		fmt.Println("server is close...")
		if err := ls.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		// Accepting an incoming TCP connection. If an error
		// occurs while accepting the connection, the
		// function returns an error.
		conn, err := ls.Accept()
		if err != nil {
			return err
		}

		s.ActiveConn++

		// Limit on the number of active connections.
		// If the limit is exceeded, the new connection
		// is immediately closed.
		if s.ActiveConn > s.MaxConn {
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
		}

		// Each new connection is handled in a
		// separate goroutine, allowing multiple
		// connections to be processed simultaneously.
		go s.ReaderLoop(conn)
	}
}

// ReaderLoop is the function that starts a loop
// for reading data from the connection.
func (s *Server) ReaderLoop(conn net.Conn) {

	// Defers (delays) the function which is called at the
	// end - closing the connection and releasing resources.
	defer func() {
		fmt.Println("connection close...")

		// Attempts to close the connection,
		// and if an error occurs, it logs it.
		if err := conn.Close(); err != nil {
			log.Println(err)
		}

		s.Lock()
		s.ActiveConn-- // Decreases the number of active connections.
		s.Unlock()
	}()

	// Creates a byte slice for the header with size HeaderSize.
	header := make([]byte, HeaderSize)

	// Creates a byte slice for allocating memory (1 MB) in case of error.
	noMemoryBlock := make([]byte, 1024*1024)

	// Starts an infinite loop for reading data from permanent the connection.
	for {
		// Reads data from the connection into the header.
		_, err := conn.Read(header)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		// Decodes the header and gets operation, key length,
		//  TTL (time-to-live), and body length.
		operation, keyLength, ttl, bodyLength := Decode(header)

		// Gets the slab index and chunk size.
		slabIndex, chunkSize := s.Slab.GetIndex(int(bodyLength + keyLength))

		// Attempts to allocate memory in the appropriate slab.
		slabBlock, err := s.Slab.ChoseSlab(slabIndex).AllocateMemory()

		if err != nil {
			// If an error occurs during memory allocation,
			// prints the error.
			fmt.Println(err)

			// If the slab is not active, sends the error to the
			// client and tries to read the rest of the request.
			if !s.Slab.GetSlabIndex(slabIndex).IsSlabActive() {
				conn.Write([]byte(err.Error()))

				//and if I can't allocate memory to my server I still have to read the req to the end
				_, err := conn.Read(noMemoryBlock)
				if err != nil {
					if err != io.EOF {
						log.Println("Error reading from connection:", err)
					}
					break
				}

				continue
			}

			// If there is no more space in memory, uses LRU
			// (Least Recently Used) policy to free up space.
			slabBlock = s.Slab.FreeSpace(slabIndex, chunkSize)
			key := slabBlock[:keyLength]

			// Deletes the key from the hash table.
			s.Distribute(key, hash_table.NewSysDelete(key))
		}

		// Reads data directly into the slabBlock.
		n, err := conn.Read(slabBlock)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}
			break
		}

		// Gets the LRU working struct and key from the slabBlock.
		lru, key := s.Slab.GetLRUIndex(slabIndex), slabBlock[:keyLength]

		// Based on the operation (S, D, G), distributes the corresponding request.
		switch operation {
		case 'S': //set operation
			field := slabBlock[keyLength:n]
			s.Distribute(key, hash_table.NewSetReq(key, conn, lru, field, ttl))
		case 'D': //delete operation
			s.Distribute(key, hash_table.NewDeleteReq(operation, key, conn, lru)) //add stack for delete
		case 'G': //get operation
			s.Distribute(key, hash_table.NewGetReq(operation, key, conn, lru))
		}
	}
}
