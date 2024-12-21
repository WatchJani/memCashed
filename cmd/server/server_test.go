package server

import (
	"log"
	"net"
	"root/client"
	"sync"
	"testing"
)

const Port string = ":5000"

func CreateConnectionWithServer() (net.Conn, error) {
	return net.Dial("tpc", Port)
}

func SetReq(key, value string, conn net.Conn) error {
	dataPayload, err := client.Set([]byte(key), []byte(value), 2121321321)
	if err != nil {
		return err
	}

	if _, err := conn.Write(dataPayload); err != nil {
		return err
	}

	return nil
}

func TestGetReq(t *testing.T) {
	conn, err := CreateConnectionWithServer()
	if err != nil {
		t.Fail()
	}

	if err := SetReq("super mario", "game", conn); err != nil {
		t.Error(err)
	}

}

// func BenchmarkSetReqPerSecond(b *testing.B) {
// 	b.StopTimer()

// 	numberOfConnection := 100
// 	SenderCh := make(chan []byte)

// 	//Workers
// 	for range numberOfConnection {
// 		conn, err := net.Dial("tcp", Port)
// 		if err != nil {
// 			log.Fatal(err)
// 			return
// 		}

// 		//write data
// 		go func(conn net.Conn) {
// 			for {
// 				payload := <-SenderCh

// 				if _, err := conn.Write(payload); err != nil {
// 					log.Println(err)
// 				}
// 			}
// 		}(conn)

// 		//get response from server
// 		go func(conn net.Conn) {
// 			buff := make([]byte, 4096)

// 			for {
// 				if _, err := conn.Read(buff); err != nil {
// 					log.Println(err)
// 				}
// 			}
// 		}(conn)
// 	}

// 	time.Sleep(100 * time.Millisecond)

// 	dataPayload, err := client.Set([]byte("super mario"), []byte("game"), 2121321321)

// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	b.StartTimer()

// 	for i := 0; i < b.N; i++ {
// 		SenderCh <- dataPayload
// 	}

// 	time.Sleep(time.Second)
// }

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
