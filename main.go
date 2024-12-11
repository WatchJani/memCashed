package main

import (
	"log"
	"root/server"
)

func main() {
	s := server.Server{
		Addr:    ":5000",
		MaxConn: 100,
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
