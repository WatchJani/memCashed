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

func New() *Server {
	config := types.LoadConfiguration()
	newAllocator := config.MemoryAllocator()

	return &Server{
		Add:     config.Port(),
		MaxConn: config.MaxConnection(),
		Manager: memory_allocator.NewSlabManager(
			config.Slabs(newAllocator),
		),
	}
}

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

// Decreases the number of active connections.
func (s *Server) decrease() {
	s.Lock()
	defer s.Unlock()

	s.ActiveConn--
}

func Close(c io.Closer, msg string) {
	fmt.Println(msg)

	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	defer func() {
		Close(conn, "connection is close")
		s.decrease()
	}()

	bufSize := make([]byte, 4)

	for {
		//read first 4 byte
		n, err := conn.Read(bufSize)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		if n < 4 {
			continue
		}

		payloadSize := client.DecodeLength(bufSize)
		slabBlock, err := s.Manager.GetSlab(payloadSize, conn)
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

		s.Req(slabBlock, conn)
	}
}

func (s *Server) Req(buf []byte, conn net.Conn) {
	s.Manager.JobCh <- memory_allocator.NewTransfer(buf, conn)
}

// just for testing
func (s *Server) GetNumberOfReq() int {
	return s.Manager.GetNumberOfReq()
}
