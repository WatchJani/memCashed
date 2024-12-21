package server

import (
	"log"
	"net"
	"github.com/WatchJani/memCashed/tree/master/client"
	"sync"
	"testing"
)

const Port string = ":5000"

func BenchmarkSynchronous(b *testing.B) {
	b.StopTimer()

	numberOfConnection := 99
	SenderCh := make(chan []byte, numberOfConnection) // Buffered channel to prevent blocking
	var wg sync.WaitGroup

	// Workers
	for i := 0; i < numberOfConnection; i++ {
		conn, err := net.Dial("tcp", Port)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		go func(conn net.Conn) {
			defer conn.Close() // Ensure connection is closed
			defer wg.Done()

			buff := make([]byte, 4096)
			for payload := range SenderCh {
				// Write data to the connection
				if _, err := conn.Write(payload); err != nil {
					log.Println("Write error:", err)
					return
				}

				// Read response from the server
				if _, err := conn.Read(buff); err != nil {
					log.Println("Read error:", err)
					return
				}
			}
		}(conn)
	}

	// Generate payload
	dataPayload, err := client.Set([]byte("super mario"), []byte("game"), 2121321321)
	if err != nil {
		log.Fatal(err)
	}

	b.StartTimer()

	// Send data
	for i := 0; i < b.N; i++ {
		SenderCh <- dataPayload
	}

	close(SenderCh) // Close the channel to signal workers to stop
	wg.Wait()       // Wait for all workers to finish
}
