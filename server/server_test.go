package server

import (
	"log"
	"net"
	"root/client"
	"testing"
	"time"
)

func Benchmark(b *testing.B) {
	b.StopTimer()

	numberOfConnection := 10
	SenderCh := make(chan []byte)
	port := ":5000"

	for range numberOfConnection {
		go func() {
			conn, err := net.Dial("tcp", port)
			if err != nil {
				log.Fatal(err)
				return
			}

			for {
				payload := <-SenderCh

				if _, err := conn.Write(payload); err != nil {
					log.Println(err)
				}
			}
		}()
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
