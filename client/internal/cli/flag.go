package cli

import "flag"

func ParseFlag() string {
	// Define a flag for the configuration file path with a default value and description.
	// The flag is named "conf" and the default value is "./config.yaml".
	path := flag.String("conf", "./config.yaml", "is used to load the custom configuration file")

	// Parse the command-line flags.
	flag.Parse()

	// Return the value of the "conf" flag. Dereference the pointer to get the string value.
	return *path
}
