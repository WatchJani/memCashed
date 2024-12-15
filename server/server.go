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

type Server struct {
	Addr       string
	MaxConn    int
	ActiveConn int
	Slab       memory_allocator.SlabManager
	*hash_table.Engine
	sync.RWMutex
}

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
		conn, err := ls.Accept()
		if err != nil {
			return err
		}

		s.Lock()
		s.ActiveConn++
		s.Unlock()

		go s.ReaderLoop(conn)
	}
}

func (s *Server) ReaderLoop(conn net.Conn) {
	defer func() {
		fmt.Println("connection close...")
		if err := conn.Close(); err != nil {
			log.Println(err)
		}

		s.Lock()
		s.ActiveConn--
		s.Unlock()
	}()

	header := make([]byte, HeaderSize)
	noMemoryBlock := make([]byte, 1024*1024)

	for {
		_, err := conn.Read(header)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		//header decode
		operation, keyLength, ttl, bodyLength := Decode(header)

		//get slab container
		slabIndex, chunkSize := s.Slab.GetIndex(int(bodyLength + keyLength))
		slabBlock, err := s.Slab.ChoseSlab(slabIndex).AllocateMemory()

		if err != nil {
			fmt.Println(err)

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

			//no more space in memory //use LRU for free space
			//add block from lru
			slabBlock = s.Slab.FreeSpace(slabIndex, chunkSize)
			key := slabBlock[:keyLength]
			//delete key in hash map
			s.Distribute(key, hash_table.NewSysDelete(key))
		}

		n, err := conn.Read(slabBlock)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}
			break
		}

		lru, key := s.Slab.GetLRUIndex(slabIndex), slabBlock[:keyLength]

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
