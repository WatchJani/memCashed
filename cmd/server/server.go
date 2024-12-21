package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"root/client"
	"root/cmd/internal/types"
	"root/cmd/memory_allocator"
	"sync"
)

const (
	PayloadSizeLength = 4
	MiB               = 1024 * 1024
	BufferSizeTCP     = MiB + PayloadSizeLength
)

type Server struct {
	Add        string
	MaxConn    int
	ActiveConn int
	sync.RWMutex
	Manager *memory_allocator.SlabManager
}

// New initializes a new Server instance by loading the configuration
// and setting up the slab memory allocator.
func New() *Server {
	config := types.LoadConfiguration()
	newAllocator := config.MemoryAllocator()

	return &Server{
		Add:     config.Port(),
		MaxConn: config.MaxConnection(),
		Manager: memory_allocator.NewSlabManager(
			config.Slabs(newAllocator),
			config.NumberWorker(),
		),
	}
}

// Run starts the server, listens for incoming TCP connections,
// and handles them concurrently. It also enforces a maximum connection limit.
func (s *Server) Run() error {
	ls, err := net.Listen("tcp", s.Add)
	if err != nil {
		return err
	}

	defer Close(ls, "server is close")

	for {
		conn, err := ls.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		s.ActiveConn++

		if s.ActiveConn > s.MaxConn {
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
		}

		go s.HandleConn(conn)
	}
}

// decrease reduces the active connection count by one.
// This method ensures thread-safety using a write lock.
func (s *Server) decrease() {
	s.Lock()
	defer s.Unlock()

	s.ActiveConn--
}

// Close safely closes an io.Closer resource (e.g., a connection or listener)
// and logs an optional message if closing fails.
func Close(c io.Closer, msg string) {
	fmt.Println(msg)

	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

// HandleConn processes an individual TCP connection, reading data,
// allocating slab memory, and delegating requests to a job channel.
func (s *Server) HandleConn(conn net.Conn) {
	defer func() {
		Close(conn, "connection is close")
		s.decrease()
	}()

	bufSize := make([]byte, 4)

	for {
		//read first 4 byte
		_, err := conn.Read(bufSize)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		payloadSize := client.DecodeLength(bufSize)
		slabBlock, index, err := s.Manager.GetSlab(payloadSize, conn)
		if err != nil {
			log.Println(err)
		}

		_, err = conn.Read(slabBlock)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		s.Req(slabBlock, index, conn)
	}
}

// Req sends a processed request to the slab manager's job channel,
// including the payload, index, and connection.
func (s *Server) Req(buf []byte, index int, conn net.Conn) {
	s.Manager.JobCh <- memory_allocator.NewTransfer(buf, index, conn)
}
