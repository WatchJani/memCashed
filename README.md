# In-Memory Database with LRU, TTL, and Custom Memory Allocator 

This is a high-performance in-memory database designed to efficiently store and manage data with built-in features such as **LRU (Least Recently Used)** caching, **TTL (Time-to-Live)** expiration, and a **custom memory allocator** that works as a **slab allocator**.

## Key Features

- **In-Memory Database**: All data is stored in memory, ensuring fast access times and low latency for data operations. The database is ideal for use cases where speed and efficiency are critical, such as caching, session management, or real-time applications.

- **LRU Cache**: The database uses a **Least Recently Used (LRU)** caching strategy to manage memory usage efficiently. The LRU algorithm ensures that the least recently accessed data is automatically evicted when the memory limit is reached, making room for more frequently accessed data.

- **TTL (Time-to-Live)**: Each entry in the database can have an associated **TTL** value, allowing data to automatically expire after a specified duration. This feature is useful for caching scenarios where data should only be retained for a limited time (e.g., session data, temporary results).

- **Custom Memory Allocator (Slab Allocator)**: The database implements a custom memory allocator that works as a **slab allocator**. This allows for more efficient memory management, especially in scenarios involving frequent memory allocation and deallocation, by reducing fragmentation and improving memory access patterns.

## Core Operations:

- **Get**: Retrieve the value associated with a specific key.
- **Set**: Add or update a key-value pair in the database.
- **Delete**: Remove a key-value pair from the database, freeing up memory.
