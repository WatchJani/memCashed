package main

import (
	"log"
	"net/http"
	"os"

	driver "github.com/WatchJani/memCashed/client/driver"
)

func main() {
	driver, err := driver.New()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	store := InitStore(driver)

	mux := http.NewServeMux()

	mux.HandleFunc("/set", store.Set)

	http.ListenAndServe(":5001", mux)
}

type InMemoryStore struct {
	*driver.Driver
}

func InitStore(driver *driver.Driver) InMemoryStore {
	return InMemoryStore{
		Driver: driver,
	}
}

func (s *InMemoryStore) Set(w http.ResponseWriter, r *http.Request) {
	resMsg, err := s.SetReq([]byte("key"), []byte("value"), -1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dbResponse := <-resMsg

	_ = dbResponse
}
