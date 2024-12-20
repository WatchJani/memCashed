package main

import (
	"fmt"
	"log"
	"net"
	"root/client"
	"root/cmd/server"
	"time"
)

func main() {
	s := server.New()

	//single core
	// go func() {
	// 	time.Sleep(100 * time.Millisecond)
	// 	data, err := client.Set([]byte("super mario"), []byte("game"), 2121321321)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	conn, err := net.Dial("tcp", ":5000")
	// 	if err != nil {
	// 		log.Println(err)
	// 	}

	// 	defer func() {
	// 		log.Println("connection close on client side")
	// 		if err := conn.Close(); err != nil {
	// 			log.Println(err)
	// 		}
	// 	}()

	// 	// buff := make([]byte, 4096)
	// 	// go func(conn net.Conn) {
	// 	// 	for {
	// 	// 		_, err := conn.Read(buff)
	// 	// 		if err != nil {
	// 	// 			log.Println(err)
	// 	// 			return
	// 	// 		}

	// 	// 		// log.Println(string(buff[:n]))
	// 	// 	}
	// 	// }(conn)

	// 	for range 1000 {
	// 		if _, err := conn.Write(data); err != nil {
	// 			log.Println(err)
	// 		}
	// 	}

	// 	time.Sleep(400 * time.Millisecond)
	// 	fmt.Println("req:", Manager.GetNumberOfReq())
	// }()

	go func() {
		time.Sleep(time.Millisecond)
		numberOfConnection := 100
		SenderCh := make(chan []byte)
		port := ":5000"

		//Workers
		for range numberOfConnection {
			conn, err := net.Dial("tcp", port)
			if err != nil {
				log.Fatal(err)
				return
			}

			//write data
			go func(conn net.Conn) {
				for {
					payload := <-SenderCh

					if _, err := conn.Write(payload); err != nil {
						log.Println(err)
					}
				}
			}(conn)

			//get response from server
			go func(conn net.Conn) {
				buff := make([]byte, 4096)

				for {
					if _, err := conn.Read(buff); err != nil {
						log.Println(err)
					}
				}
			}(conn)
		}

		time.Sleep(100 * time.Millisecond)

		dataPayload, err := client.Set([]byte("super mario"), []byte("game"), 2121321321)

		if err != nil {
			log.Fatal(err)
			return
		}

		for i := 0; i < 10_000_000; i++ {
			SenderCh <- dataPayload
		}

		time.Sleep(12_000 * time.Millisecond)
		fmt.Println("req:", s.GetNumberOfReq())
	}()

	if err := s.Run(); err != nil {
		log.Println(err)
	}
}
