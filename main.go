package main

import (
	"log"
	"net"
	"root/client"
	"root/hash_table"
	"root/memory_allocator"
	"root/server"
	"time"
)

const (
	MaxConnection int    = 100
	Port          string = ":5000"
)

func main() {
	newAllocator := memory_allocator.New(5 * 1024 * 1024 * 1024)

	slabCapacity := []int{
		64,
		128,
		256,
		512,
		1024,
		2048,
		4096,
		8192,
		16384,
		32768,
		65536,
		131072,
		262144,
		524288,
		1048576,
	}

	slabAllocator := make([]memory_allocator.Slab, len(slabCapacity))
	for i := range slabAllocator {
		slabAllocator[i] = memory_allocator.NewSlab(slabCapacity[i], newAllocator)
	}

	s := server.Server{
		Addr:    Port,
		MaxConn: MaxConnection,
		Slab:    memory_allocator.NewSlabManager(slabAllocator),
		Engine:  hash_table.NewEngine(15),
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		data, err := client.Set([]byte("super mario"), []byte("game"), 2121321321)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.Dial("tcp", Port)
		if err != nil {
			log.Println(err)
		}

		defer func() {
			log.Println("connection close on client side")
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
		}()

		buff := make([]byte, 4096)
		go func(conn net.Conn) {
			n, err := conn.Read(buff)
			if err != nil {
				log.Println(err)
			}

			log.Println(string(buff[:n]))
		}(conn)

		// for {
		if _, err := conn.Write(data); err != nil {
			log.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
		// }
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
