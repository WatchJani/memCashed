package server

import (
	"log"
	"net"
	"root/client"
	"testing"
	"time"
)

// no memory space 1420ns
// free memory 1100ns
func Benchmark(b *testing.B) {
	b.StopTimer()

	numberOfConnection := 10
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

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		SenderCh <- dataPayload
	}
}
