package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct {
	Addr       string
	MaxConn    int
	ActiveConn int
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

	// headerSize := 10
	readBuffer := make([]byte, 10)

	for {
		n, err := conn.Read(readBuffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		header := readBuffer[:n]
		operation, keyLength, ttl, bodyLength := Decode(header)
		fmt.Println("operation:", operation)
		// fmt.Println("key length:", keyLength)
		fmt.Println("ttl", ttl)
		// fmt.Println("body length", bodyLength)

		// // bodyLength := int(readBuffer[20])
		// remainingBody := bodyLength - (n - headerSize)

		// if remainingBody > 0 {
		bodyBuffer := make([]byte, bodyLength+keyLength)

		// copy(bodyBuffer, readBuffer[headerSize:n])

		_, err = conn.Read(bodyBuffer)
		if err != nil && err != io.EOF {
			log.Println("Error reading body:", err)
		}

		fmt.Println("Header:", header)
		fmt.Println("Body:", bodyBuffer)
		// } else {
		// 	bodyBuffer := make([]byte, n-headerSize)
		// 	copy(bodyBuffer, readBuffer[headerSize:n])

		// 	fmt.Println("Header:", header)
		// 	fmt.Println("Body:", bodyBuffer)
		// }
	}
}
