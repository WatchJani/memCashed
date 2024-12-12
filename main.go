package main

import (
	"log"
	"net"
	"root/client"
	"root/server"
	"time"
)

func main() {
	s := server.Server{
		Addr:    ":5000",
		MaxConn: 100,
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		data, err := client.Set([]byte("super mario"), []byte("game"), 2121321321)
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.Dial("tcp", ":5000")
		if err != nil {
			log.Println(err)
		}

		defer func() {
			log.Println("connection close on client side")
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
		}()

		if _, err := conn.Write(data); err != nil {
			log.Println(err)
		}
	}()

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
