package server

import (
	"log"
	"net"
	"sync"
	"testing"

	client "github.com/WatchJani/memCashed/client/driver"
)

const (
	Port               = ":5000"
	NumberOfConnection = 99
)

var (
	PayloadKey        = []byte("super mario")
	PayloadValue      = []byte("game")
	PayloadTTLDefault = -1
)

func Workers(SenderCh chan []byte, wg *sync.WaitGroup) {
	for i := 0; i < NumberOfConnection; i++ {
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
}

func BenchmarkSynchronousSet(b *testing.B) {
	b.StopTimer()

	SenderCh := make(chan []byte) // Buffered channel to prevent blocking
	var wg sync.WaitGroup

	// Workers
	Workers(SenderCh, &wg)

	b.StartTimer()

	// Send data
	for i := 0; i < b.N; i++ {
		dataPayload, err := client.Set(PayloadKey, PayloadValue, PayloadTTLDefault)
		if err != nil {
			log.Fatal(err)
		}

		SenderCh <- dataPayload
	}

	close(SenderCh) // Close the channel to signal workers to stop
	wg.Wait()       // Wait for all workers to finish
}

func BenchmarkSynchronousGet(b *testing.B) {
	b.StopTimer()

	SenderCh := make(chan []byte) // Buffered channel to prevent blocking
	var wg sync.WaitGroup

	// Workers
	Workers(SenderCh, &wg)

	// Generate payload
	dataPayload, err := client.Get(PayloadKey)
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

func BenchmarkSynchronousDelete(b *testing.B) {
	b.StopTimer()

	SenderCh := make(chan []byte) // Buffered channel to prevent blocking
	var wg sync.WaitGroup

	// Workers
	Workers(SenderCh, &wg)

	// Generate payload
	dataPayload, err := client.Delete(PayloadKey)
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
