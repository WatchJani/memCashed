# In-Memory Database with LRU, TTL, and Custom Memory Allocator

This is a high-performance in-memory database designed to efficiently store and manage data with built-in features such as **LRU (Least Recently Used)** caching, **TTL (Time-to-Live)** expiration, and a **custom memory allocator** that works as a **slab allocator**.

## Key Features

- **In-Memory Database**: All data is stored in memory, ensuring fast access times and low latency for data operations. The database is ideal for use cases where speed and efficiency are critical, such as caching, session management, or real-time applications.

- **LRU Cache**: The database uses a **Least Recently Used (LRU)** caching strategy to manage memory usage efficiently. The LRU algorithm ensures that the least recently accessed data is automatically evicted when the memory limit is reached, making room for more frequently accessed data.

- **TTL (Time-to-Live)**: Each entry in the database can have an associated **TTL** value, allowing data to automatically expire after a specified duration. This feature is useful for caching scenarios where data should only be retained for a limited time (e.g., session data, temporary results).

- **Custom Memory Allocator (Slab Allocator)**: The database implements a custom memory allocator that works as a **slab allocator**. This allows for more efficient memory management, especially in scenarios involving frequent memory allocation and deallocation, by reducing fragmentation and improving memory access patterns.

- **Full Vertical Scalability**: The system is designed to scale efficiently with the hardware. It supports **vertical scaling**, meaning it can take full advantage of multi-core processors and scale up performance by utilizing all available CPU cores for parallel processing. This ensures high throughput and low latency even as the data size or workload increases.

## Core Operations:

- **Get**: Retrieve the value associated with a specific key.
- **Set**: Add or update a key-value pair in the database.
- **Delete**: Remove a key-value pair from the database, freeing up memory.

## Benefits

- **Speed**: As an in-memory database, operations like reading, writing, and deleting data are extremely fast, with low latency.
- **Memory Efficiency**: The LRU cache and custom memory allocator ensure that the system uses memory efficiently, even with large datasets.
- **Data Expiry**: The TTL feature provides a powerful mechanism to ensure that stale data is automatically removed.
- **Vertical Scalability**: The system can scale effectively with the hardware by leveraging multiple CPU cores, ensuring maximum utilization and performance under high-load scenarios.

## Performance Considerations

- The system is optimized for speed and efficiency, using in-memory storage and an efficient custom memory allocator.
- LRU ensures that only the most frequently accessed data remains in memory, which improves performance during heavy workloads.
- The ability to fully utilize all CPU cores provides excellent parallelization and enhances the overall performance of the system as the data or workload grows.


# Installation Guide for the Go memCached Server

This guide will walk you through the steps to install and run the memCached server. 

## Prerequisites

- Ensure you have **Golang** installed on your system. You can download it from [golang.org](https://golang.org).
- Optionally, have `make` installed for building the project.

---

## Installation Steps

1. **Download the Server**

   Use `go get` to fetch the server code from the repository:
   ```bash
   go get github.com/WatchJani/memCashed/memcached
   ```

2. **Build the Executable**
	Use the `make` command to build the executable:
	```bash
    make build
    ```
3. **Run the Server**
	```bash
    ./memcached
    ```


# How the Server Works

The server uses a multi-threaded architecture to efficiently handle client requests and execute operations. Here's a brief explanation:

1. **Connection Handling**:  
   Each new connection creates a dedicated thread that parses and processes incoming requests.

2. **Worker Pool**:  
   Parsed requests are sent to a pool of workers. These workers are responsible for performing the actual operations, such as reading, writing, or deleting data.

3. **Optimized Workflow**:  
   This separation between connection handling and request execution ensures better performance and scalability.

---

## Architecture Diagram

Below is a simplified representation of how the server processes requests:

![Server Architecture Diagram](https://github.com/WatchJani/memCashed/blob/master/assets/server.png)

1. **Threads**: Handle client connections and parse requests.
2. **Workers**: Perform the actual operations based on parsed requests.

---

This design allows the server to handle multiple requests concurrently while ensuring efficient resource utilization.


# Go Driver Installation and Usage

A specific driver written for the Go programming language allows seamless integration with the server. Here's how to add it to your project and use it effectively.

---

## Adding the Driver to Your Project

To include the driver in your project, use the following command:

```bash
go get github.com/WatchJani/memCashed/client
```


## Example Usage

```go
package main

import (
	"log"
	"net/http"
	"os"

	driver "github.com/WatchJani/memCashed/client/driver"
)

func main() {
	driver, err := driver.New(":5000", 15)
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
```