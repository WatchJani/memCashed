package cli

import "flag"

func ParseFlag() string {
	path := flag.String("conf", "./config.yaml", "is used to load the custom configuration file")

	flag.Parse()

	return *path
}
