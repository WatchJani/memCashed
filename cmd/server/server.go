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

// just for testing
func (s *Server) GetNumberOfReq() int {
	return s.Manager.GetNumberOfReq()
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

		go s.HandleConn(conn)
	}
}

func Close(c io.Closer, msg string) {
	fmt.Println(msg)

	if err := c.Close(); err != nil {
		log.Println(err)
	}
}

func (s *Server) HandleConn(conn net.Conn) {
	defer Close(conn, "connection is close")

	var (
		buf         = make([]byte, BufferSizeTCP)
		pointer     int
		active      bool
		payloadSize = make([]byte, PayloadSizeLength)
		slabBlock   []byte
	)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		if active {
			temp := pointer
			pointer = 0

			var end int
			if temp < PayloadSizeLength {
				pointer = PayloadSizeLength - temp
				temp += copy(payloadSize[temp:], buf[:pointer])
				temp = 0

				end = client.DecodeLength(payloadSize)
				slabBlock, err = s.Manager.GetSlab(end, conn)
				if err != nil {
					log.Println(err)
				}

			} else {
				end = client.DecodeLength(payloadSize)
			}

			copy(slabBlock[temp:], buf[pointer:end-temp+PayloadSizeLength])
			pointer += end - temp
		}

		for {
			if !active {
				if pointer+PayloadSizeLength > n {
					copy(payloadSize, buf[pointer:])
					pointer = n - pointer
					active = true
					break
				} else {
					copy(payloadSize, buf[pointer:pointer+PayloadSizeLength])
				}

				end := pointer + client.DecodeLength(buf[pointer:pointer+PayloadSizeLength]) + PayloadSizeLength
				slabBlock, err = s.Manager.GetSlab(end-pointer, conn)
				if err != nil {
					log.Println(err)
				}

				pointer += PayloadSizeLength
				if end > n {
					copy(slabBlock, buf[pointer:])
					pointer = n - pointer
					active = true
					break
				}

				copy(slabBlock, buf[pointer:end])
				pointer = end
			}

			s.Req(slabBlock, conn)
			active = false
		}
	}
}

func (s *Server) Req(buf []byte, conn net.Conn) {
	s.Manager.JobCh <- memory_allocator.NewTransfer(buf, conn)
}
