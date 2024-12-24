package main

import (
	"log"

	"github.com/WatchJani/memCashed/tree/master/memcached/server"
)

// main is the entry point of the application. It initializes a new server instance
// and attempts to start it by calling the Run method. If there is an error while
// starting the server, it logs the error.
func main() {
	// Create a new server instance using the New method.
	// Run the server and check if it returns any error.
	if err := server.New().Run(); err != nil {
		// If an error occurs, log the error message.
		log.Println(err)
	}
}
