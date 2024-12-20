package types

import (
	"fmt"
	"log"
	"os"
	"root/cmd/internal/cli"
	"root/cmd/memory_allocator"

	"gopkg.in/yaml.v3"
)

const (
	MiB                       = 1024
	MinimumNumberOfConnection = 5
	IntDefaultValue           = 0
	DefaultPort               = 5001
	DefaultNumberOfWorkers    = 15
)

type Config struct {
	Server         ServerConfig `yaml:"server"`
	MemoryAllocate int          `yaml:"memory_for_allocate"` //By Default 5GiB
	NumberOfWorker int          `yaml:"number_of_worker"`
	DefaultSlab    []CustomSlab `yaml:"custom_slabs"`
}

func NewConfig() *Config {
	return &Config{}
}

type ServerConfig struct {
	Port          int `yaml:"port"`                  //By default 5001
	MaxConnection int `yaml:"max_number_connection"` //By default 100
}

type CustomSlab struct {
	Capacity          int `yaml:"chunk_capacity"`
	MaxMemoryAllocate int `yaml:"max_allocate_memory"`
}

func LoadConfiguration() *Config {
	path := cli.ParseFlag()

	configData, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}

	config := NewConfig()

	if err := yaml.Unmarshal(configData, config); err != nil {
		log.Fatal(err)
	}

	return config
}

func (c *Config) MemoryAllocator() *memory_allocator.Allocator {
	memorySize := c.MemoryAllocate

	if memorySize == IntDefaultValue {
		memorySize = 1 //can load minimum 1MiB
	}

	return memory_allocator.New(memorySize * MiB)
}

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

func (c *Config) Slabs(allocator *memory_allocator.Allocator) []memory_allocator.Slab {
	slabs := c.DefaultSlab

	if len(c.DefaultSlab) == IntDefaultValue {
		slabs = DefaultSlabs()
	}

	slabAllocator := make([]memory_allocator.Slab, len(slabs))
	for i := range slabAllocator {
		slab := slabs[i]
		slabAllocator[i] = memory_allocator.NewSlab(slab.Capacity, slab.MaxMemoryAllocate, allocator)
	}

	return slabAllocator
}

func (c *Config) MaxConnection() int {
	maxConnection := c.Server.MaxConnection

	if maxConnection < MinimumNumberOfConnection {
		maxConnection = MinimumNumberOfConnection
	}

	return maxConnection
}

func (c *Config) Port() string {
	port := c.Server.Port
	if port < 1 {
		port = DefaultPort //Default port
	}

	return fmt.Sprintf(":%d", port)
}
