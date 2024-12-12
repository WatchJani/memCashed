package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"root/memory_allocator"
)

const HeaderSize = 10

type Server struct {
	Addr       string
	MaxConn    int
	ActiveConn int
	Slab       memory_allocator.SlabManager
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

	readBuffer := make([]byte, HeaderSize)

	for {
		_, err := conn.Read(readBuffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		header := readBuffer
		operation, keyLength, ttl, bodyLength := Decode(header)
		fmt.Println("operation:", operation)
		fmt.Println("ttl", ttl)

		bodyBuffer := make([]byte, bodyLength+keyLength) //switch with real allocator slab

		n, err := conn.Read(bodyBuffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		fmt.Println("Header:", header)
		fmt.Println("Body:", bodyBuffer[:n])
	}
}
