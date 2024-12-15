package main

import (
	"log"
	"net"
	"os"
	"root/client"
	"root/hash_table"
	"root/memory_allocator"
	"root/server"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	MaxConnection int    = 100
	Port          string = ":5000"
)

type ServerConfig struct {
	Port          int `yaml:"port"`                  //By default 5001
	MaxConnection int `yaml:"max_number_connection"` //By default 100
}

type Config struct {
	Server ServerConfig `yaml:"server"`

	MemoryAllocate int `yaml:"memory_for_allocate"` //By Default 5GiB

	NumberOfWorker int `yaml:"number_of_worker"`

	DefaultSlab []CustomSlab `yaml:"custom_slabs"`
}

type CustomSlab struct {
	Capacity          int `yaml:"chunk_capacity"`
	MaxMemoryAllocate int `yaml:"max_allocate_memory"`
}

func main() {
	conf := Config{
		Server: ServerConfig{
			Port:          5000,
			MaxConnection: 100,
		},
		MemoryAllocate: 5 * 1024 * 1024 * 1024,

		NumberOfWorker: 15,

		DefaultSlab: []CustomSlab{
			{
				Capacity:          64,
				MaxMemoryAllocate: 0, //by default
			}, {
				Capacity: 128,
			},
			{
				Capacity: 256,
			},
			{
				Capacity: 512,
			},
			{
				Capacity: 1024,
			},
			{
				Capacity: 2048,
			},
			{
				Capacity: 4096,
			},
			{
				Capacity: 4096,
			},
			{
				Capacity: 8192,
			},
			{
				Capacity: 16384,
			},
			{
				Capacity: 32768,
			},
			{
				Capacity: 65536,
			},
			{
				Capacity: 131072,
			},
			{
				Capacity: 262144,
			},
			{
				Capacity: 524288,
			},
			{
				Capacity: 1048576,
			},
		},
	}

	output, err := yaml.Marshal(&conf)
	if err != nil {
		log.Println(err)
	}

	if err := os.WriteFile("./config.yaml", output, 0777); err != nil {
		log.Println(err)
	}

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
			for {
				n, err := conn.Read(buff)
				if err != nil {
					log.Println(err)
					return
				}

				log.Println(string(buff[:n]))
			}
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
