package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"root/hash_table"
	"root/memory_allocator"
)

const HeaderSize = 10

type Server struct {
	Addr       string
	MaxConn    int
	ActiveConn int
	Slab       memory_allocator.SlabManager
	*hash_table.Engine
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

		s.ActiveConn++
		go s.ReaderLoop(conn)
	}
}

func (s *Server) ReaderLoop(conn net.Conn) {
	defer func() {
		fmt.Println("connection close...")
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	header := make([]byte, HeaderSize)

	for {
		_, err := conn.Read(header)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		operation, keyLength, ttl, bodyLength := Decode(header)

		slabIndex, chunkSize := s.Slab.GetIndex(int(bodyLength + keyLength))
		slabBlock, err := s.Slab.ChoseSlab(slabIndex).AllocateMemory()
		if err != nil {
			fmt.Println(err)
			//no more space in memory //use LRU for free space
			slabBlock = s.Slab.FreeSpace(slabIndex, chunkSize)
			//add block from lru
		}

		n, err := conn.Read(slabBlock)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}
			fmt.Println("yes")
			break
		}

		switch operation {
		case 'S':
			//pointer on key
			key := slabBlock[:keyLength]
			// fmt.Println("key", key)

			//pointer on field //real data
			field := slabBlock[keyLength:n]
			// fmt.Println("field", field)
			// slabManager.
			s.Insert(key, field, ttl) //lru list
		case 'D':

		case 'G':

		}

		// fmt.Println("Header:", header)
		// fmt.Println("Body:", slabBlock[:n])
	}
}
