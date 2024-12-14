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