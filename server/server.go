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

		//header decode
		operation, keyLength, ttl, bodyLength := Decode(header)

		//get slab container
		slabIndex, chunkSize := s.Slab.GetIndex(int(bodyLength + keyLength))
		slabBlock, err := s.Slab.ChoseSlab(slabIndex).AllocateMemory()

		key := slabBlock[:keyLength]

		if err != nil {
			fmt.Println(err)
			//no more space in memory //use LRU for free space
			//add block from lru
			slabBlock = s.Slab.FreeSpace(slabIndex, chunkSize)

			//delete key in hash map
			s.Distribute(key, hash_table.NewSysDelete(key))
		}

		n, err := conn.Read(slabBlock)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}
			fmt.Println("yes")
			break
		}

		lru := s.Slab.GetLRUIndex(slabIndex)

		switch operation {
		case 'S': //set operation
			field := slabBlock[keyLength:n]
			s.Distribute(key, hash_table.NewSetReq(key, conn, lru, field, ttl))
		case 'D': //delete operation
			s.Distribute(key, hash_table.NewDeleteReq(operation, key, conn, lru))
		case 'G': //get operation
			s.Distribute(key, hash_table.NewGetReq(operation, key, conn, lru))
		}
	}
}
