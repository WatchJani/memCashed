package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// mux := NewServerMux()

	// mux.HandleFunc("SET", func() {

	// })

	s := Server{
		Addr:    ":5000",
		MaxConn: 100,
		// Handler: mux,
	}

	go func() {
		var (
			buf bytes.Buffer
			err error
		)

		conn, err := net.Dial("tcp", "localhost:5000")
		if err != nil {
			log.Println(err)
		}

		defer func() {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
		}()

		//S_

		//Operation - 1byte [x]
		//Key length - 1byte [x]
		//TTL - 4byte [x]
		//body length - 4byte
		//End of req - \r\n 2byte

		key := []byte("janko")
		value := []byte("super dan za pobjedu\r\n")
		ttl := 1212121

		//set operation
		if err = buf.WriteByte('S'); err != nil {
			fmt.Println("Error writing operation:", err)
			return
		}

		//key length
		keyLength := uint8(len(key))
		err = binary.Write(&buf, binary.LittleEndian, keyLength)
		if err != nil {
			fmt.Println("Error writing key length:", err)
			return
		}

		//set ttl
		TTL := uint32(ttl)
		err = binary.Write(&buf, binary.LittleEndian, TTL)
		if err != nil {
			fmt.Println("Error writing TTL value:", err)
			return
		}

		bodyLength := uint32(len(value))
		err = binary.Write(&buf, binary.LittleEndian, bodyLength)
		if err != nil {
			fmt.Println("Error writing body length:", err)
			return
		}

		if _, err = buf.Write(key); err != nil {
			fmt.Println("Error writing key:", err)
			return
		}

		if _, err = buf.Write(value); err != nil {
			fmt.Println("Error writing value:", err)
			return
		}

		fmt.Printf("Header: %v\n", buf.String())
		fmt.Printf("Header length: %d\n", len(buf.Bytes())) // Ovo bi trebalo da ispiše dužinu headera

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			log.Println(err)
		}
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

type Server struct {
	Addr       string
	MaxConn    int
	ActiveConn int
	// *Handler
}

// type HandlerFn func()

// type Handler struct {
// 	router map[string]HandlerFn
// }

// func (h *Handler) HandleFunc(path string, fn HandlerFn) {
// 	h.router[path] = fn
// }

// func NewServerMux() *Handler {
// 	return &Handler{
// 		router: make(map[string]HandlerFn, 2),
// 	}
// }

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

	headerSize := 24
	readBuffer := make([]byte, 1024)

	for {
		n, err := conn.Read(readBuffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading from connection:", err)
			}

			break
		}

		header := readBuffer[:headerSize]
		bodyLength := int(readBuffer[20])
		remainingBody := bodyLength - (n - headerSize)

		if remainingBody > 0 {
			bodyBuffer := make([]byte, bodyLength)

			copy(bodyBuffer, readBuffer[headerSize:n])

			_, err := conn.Read(bodyBuffer[n-headerSize:])
			if err != nil && err != io.EOF {
				log.Println("Error reading body:", err)
			}

			fmt.Println("Header:", header)
			fmt.Println("Body:", bodyBuffer)
		} else {
			bodyBuffer := make([]byte, n-headerSize)
			copy(bodyBuffer, readBuffer[headerSize:n])

			fmt.Println("Header:", header)
			fmt.Println("Body:", bodyBuffer)
		}
	}
}

//req body

//CMD SET\n
//TTL 6565151\n
//\n
//body.....
