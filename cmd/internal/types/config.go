package types

import (
	"fmt"
	"log"
	"os"
	"github.com/WatchJani/memCashed/cmd/internal/cli"
	"github.com/WatchJani/memCashed/cmd/memory_allocator"

	"gopkg.in/yaml.v3"
)

// Constant values used for memory operations and server configuration.
const (
	KiB                       = 1024 // 1 MiB in bytes
	MinimumNumberOfConnection = 5           // Minimum number of connections to the server
	IntDefaultValue           = 0           // Default value for integers
	DefaultPort               = 5001        // Default server port
	DefaultNumberOfWorkers    = 15          // Default number of worker threads
)

// Configuration structure containing server and memory details.
type Config struct {
	Server         ServerConfig `yaml:"server"`              // Server configuration
	MemoryAllocate int          `yaml:"memory_for_allocate"` // Amount of memory allocated (default 5GiB)
	NumberOfWorker int          `yaml:"number_of_worker"`    // Number of worker threads for the server
	DefaultSlab    []CustomSlab `yaml:"custom_slabs"`        // Default slab sizes
}

// Creates and returns a new instance of the `Config` structure.
func NewConfig() *Config {
	return &Config{}
}

// Server configuration structure, including port and max connection count.
type ServerConfig struct {
	Port          int `yaml:"port"`                  // Port the server listens on (default 5001)
	MaxConnection int `yaml:"max_number_connection"` // Maximum number of connections to the server (default 100)
}

// Defines slab structures with capacities and maximum memory allocations.
type CustomSlab struct {
	Capacity          int `yaml:"chunk_capacity"`      // Capacity of each slab (in bytes)
	MaxMemoryAllocate int `yaml:"max_allocate_memory"` // Maximum memory that can be allocated to the slab
}

// Function that loads configuration from a YAML file.
func LoadConfiguration() *Config {
	path := cli.ParseFlag() // Parses the file path from command line flags

	// Reading configuration file data
	configData, err := os.ReadFile(path) // Logs error if there is an issue reading the file
	if err != nil {
		log.Println(err)
	}

	config := NewConfig()

	// Parsing the YAML data into the `Config` structure
	if err := yaml.Unmarshal(configData, config); err != nil {
		log.Fatal(err) // If there is an error during parsing, the program terminates
	}

	return config // Returns the loaded configuration
}

// Function that returns a memory allocator based on the `MemoryAllocate` value from the configuration.
func (c *Config) MemoryAllocator() *memory_allocator.Allocator {
	memorySize := c.MemoryAllocate

	// If no memory size is defined, allocate a minimum of 1 MiB
	if memorySize == IntDefaultValue {
		memorySize = 1 // can load at least 1 MiB
	}

	// Returns a new memory allocator with the specified memory size
	return memory_allocator.New(memorySize * KiB)
}

// Returns the default slabs with predefined capacities and maximum memory allocations.
func DefaultSlabs() []CustomSlab {
	return []CustomSlab{
		{64, 0},
		{128, 0},
		{256, 0},
		{512, 0},
		{1024, 0},
		{2048, 0},
		{4096, 0},
		{8192, 0},
		{16384, 0},
		{32768, 0},
		{65536, 0},
		{131072, 0},
		{262144, 0},
		{524288, 0},
		{1048576, 0},
	}
}

// Configures and returns a list of slabs based on the current configuration and memory allocator.
func (c *Config) Slabs(allocator *memory_allocator.Allocator) []memory_allocator.Slab {
	slabs := c.DefaultSlab

	// If no slabs are defined in the configuration, use the default slabs.
	if len(c.DefaultSlab) == IntDefaultValue {
		slabs = DefaultSlabs()
	}

	// Allocate memory for slabs based on the specified configuration.
	slabAllocator := make([]memory_allocator.Slab, len(slabs))
	for i := range slabAllocator {
		slab := slabs[i]
		slabAllocator[i] = memory_allocator.NewSlab(slab.Capacity, slab.MaxMemoryAllocate, allocator)
	}

	return slabAllocator // Return the configured slabs
}

// Returns the maximum number of connections, ensuring it meets the minimum required value.
func (c *Config) MaxConnection() int {
	maxConnection := c.Server.MaxConnection

	// If the configured value is less than the minimum allowed, set it to the minimum.
	if maxConnection < MinimumNumberOfConnection {
		maxConnection = MinimumNumberOfConnection
	}

	return maxConnection // Return the maximum number of connections
}

// Returns the server's port as a formatted string, using the default port if the configured port is invalid.
func (c *Config) Port() string {
	port := c.Server.Port
	if port < 1 {
		port = DefaultPort //Default port
	}

	return fmt.Sprintf(":%d", port) // Format the port as a string (e.g., ":5001")
}

func (c *Config) NumberWorker() int {
	numberOfWorker := c.NumberOfWorker
	if numberOfWorker < 1 {
		numberOfWorker = DefaultNumberOfWorkers
	}

	return numberOfWorker	
}
