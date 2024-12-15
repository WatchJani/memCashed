package main

import (
	"log"
	"net"
	"root/client"
	"root/server"
	"time"
)

func main() {
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

	if err := server.NewServer().Run(); err != nil {
		log.Fatal(err)
	}
}
