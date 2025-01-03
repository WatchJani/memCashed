package types

import (
	"log"
	"os"

	"github.com/WatchJani/memCashed/client/internal/cli"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	// Client Client   `yaml:"client"`
	Server []Server `yaml:"server"`
}

type Server struct {
	IpAddr             string `yaml:"ip_address"`
	NumberOfConnection int    `yaml:"number_of_connection"`
}

func NewConfig() *Configuration {
	return &Configuration{}
}

func LoadConfiguration() *Configuration {
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
