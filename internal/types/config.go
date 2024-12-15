package types

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
