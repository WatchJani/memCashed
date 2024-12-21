package main

import (
	"fmt"
	"log"
	"root/client"
	"root/cmd/server"
	"time"
)

func main() {
	s := server.New()

	// go func() {
	// 	time.Sleep(time.Millisecond)
	// 	numberOfConnection := 15
	// 	SenderCh := make(chan []byte)

	// 	//Workers
	// 	for range numberOfConnection {
	// 		conn, err := net.Dial("tcp", ":5000")
	// 		if err != nil {
	// 			log.Fatal(err)
	// 			return
	// 		}

	// 		//write data
	// 		go func(conn net.Conn) {
	// 			buff := make([]byte, 4096)
	// 			for {
	// 				payload := <-SenderCh
	// 				if _, err := conn.Write(payload); err != nil {
	// 					log.Println(err)
	// 				}

	// 				_, err := conn.Read(buff)
	// 				if err != nil {
	// 					log.Println(err)
	// 				}
	// 			}
	// 		}(conn)
	// 	}

	// 	dataPayload, err := client.Set([]byte("super mario"), []byte("game"), 2121321321)

	// 	if err != nil {
	// 		log.Fatal(err)
	// 		return
	// 	}

	// 	start := time.Now()

	// 	for i := 0; i < 1_350_000; i++ {
	// 		SenderCh <- dataPayload
	// 	}

	// 	fmt.Println(time.Since(start))
	// }()

	go func() {
		time.Sleep(100 * time.Millisecond)

		driver := client.New(":5000", 15)

		if err := driver.Init(); err != nil {
			fmt.Println(err)
		}

		var (
			key   = []byte("key")
			value = []byte("value")
			ttl   = -1
		)

		start := time.Now()
		response, err := driver.SetReq(key, value, ttl)
		if err != nil {
			fmt.Println(err)
		}

		data := <-response
		fmt.Println(time.Since(start))
		fmt.Println(string(data))
	}()

	if err := s.Run(); err != nil {
		log.Println(err)
	}
}
